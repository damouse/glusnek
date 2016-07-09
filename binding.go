package gosnake

// Binding between go and python, with pthreading implmementation
// Python and C are leaking badly. Comment out big chunks of functionality and parsing until its found,
// then work backwards.
// Also try working off just one pthread.

//#cgo pkg-config: python-2.7
//#include "Python.h"
//#include "binding.h"
import "C"

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/sbinet/go-python"
)

type ExportedFunction func([]interface{}, map[string]interface{}) ([]interface{}, error)

type Thread uintptr
type ThreadCallback func()

var create_callback chan ThreadCallback

// Because the invocation method is static we have to hold a global reference to these guys. Not ecstatic about that
var bindingsLock *sync.Mutex = &sync.Mutex{}
var allBindings []*Binding = []*Binding{}

func init() {
	C.initialize_python(C.int(0))
	C.register_sig_handler()
	create_callback = make(chan ThreadCallback, 1)
}

// Wraps all state and management information for a binding
type Binding struct {
	lock *sync.Mutex

	thread uintptr

	threadCallback chan ThreadCallback

	imports map[string]*python.PyObject
	exports map[string]*ExportedFunction
}

func NewBinding() *Binding {
	b := &Binding{
		lock:           &sync.Mutex{},
		thread:         0,
		threadCallback: make(chan ThreadCallback, 1),
		imports:        map[string]*python.PyObject{},
		exports:        map[string]*ExportedFunction{},
	}

	bindingsLock.Lock()
	defer bindingsLock.Unlock()
	allBindings = append(allBindings, b)
	return b
}

// Equivalent to "import name" in python
func (b *Binding) Import(name string) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	done := make(chan error)
	defer close(done)

	thread := PtCreate(func() {
		gil := C.PyGILState_Ensure()
		defer C.PyGILState_Release(gil)

		if m := python.PyImport_ImportModuleNoBlock(name); m == nil {
			done <- b.parseException()
		} else {
			b.imports[name] = m
			done <- nil
		}
	})

	defer thread.Kill()
	return <-done
}

// Call a python function on passed module. Fails if Binding.Import(moduleName) has not already succeeded
// If the call raises an exception in python its passed along as a go error
func (b *Binding) Call(module string, function string, args ...interface{}) (interface{}, error) {

	// See if the module has been imported
	m, ok := b.imports[module]

	if !ok {
		return nil, fmt.Errorf("Cant call python function, module %s has not been imported. Did you call Import(moduleName)?", module)
	}

	b.lock.Lock()
	defer b.lock.Unlock()

	errChan := make(chan error)
	defer close(errChan)
	resultChan := make(chan interface{})
	defer close(resultChan)

	thread := PtCreate(func() {
		gil := C.PyGILState_Ensure()
		defer C.PyGILState_Release(gil)

		if fn := m.GetAttrString(function); fn == nil {
			errChan <- b.parseException()

		} else if tup, e := packTuple(args); e != nil {
			errChan <- e

		} else if ret := fn.CallObject(tup); ret == nil {
			python.PyErr_PrintEx(false)
			errChan <- b.parseException()

		} else if cast, e := togo(ret); e != nil {
			errChan <- e

		} else {
			resultChan <- cast
		}
	})

	defer thread.Kill()

	select {
	case e := <-errChan:
		return nil, e
	case r := <-resultChan:
		return r, nil
	}

}

// Make a go method that takes a slice and a map callable from python
// TODO
//      Raise exceptions in python
//      Check if name is already taken
func (b *Binding) Export(meth string, fn ExportedFunction) {
	b.exports[meth] = &fn
}

// Process a python exception: return the reason, the stack trace, and clear the exception flag
func (b *Binding) parseException() error {
	python.PyErr_PrintEx(false)

	// TODO: extract the exception information and format it nicely instead of just printing out
	return fmt.Errorf("Python operation failed")
}

// An invocation FROM python to go
// TODO: throw errors back to python!
//export pyInvocation
func pyInvocation(self *C.PyObject, args *C.PyObject) *C.PyObject {
	fmt.Println("GO: invocation", self, args)
	transform := python.PyObject_FromVoidPtr(unsafe.Pointer(args))

	var target string
	var goArgs []interface{}

	if unpacked, e := togo(transform); e != nil {
		fmt.Println("GO: an error occured! ", e)
	} else if arr, ok := unpacked.([]interface{}); !ok {
		fmt.Println("GO: an error occured! ")
	} else {
		goArgs = arr
	}

	if len(goArgs) < 1 {
		fmt.Println("GO: not enough arguments to gocall")
	} else if t, ok := goArgs[0].(string); !ok {
		fmt.Println("GO: method name wasnt a string")
	} else {
		target = t
	}

	if len(goArgs) != 2 {
		goArgs = nil
	} else if f, ok := goArgs[1].([]interface{}); !ok {
		fmt.Println("GO: not passed a list as method arguments")
	} else {
		goArgs = f
	}

	// Search for exported method
	bindingsLock.Lock()
	var fn ExportedFunction
	for _, binding := range allBindings {
		for name, f := range binding.exports {
			if name == target {
				fn = *f
			}
		}
	}
	bindingsLock.Unlock()

	if fn == nil {
		fmt.Println("GO: no function exported as ", target)
	}

	results, err := fn(goArgs, nil)

	if err != nil {
		fmt.Println("GO: exported function erred. Name:", target, err)
	}

	// Works, is a c type
	return C.PyLong_FromLongLong(1)

	// cstr := python.PyString_FromString("asdfasdfasdf")
	// nptr := (*C.PyObject)(unsafe.Pointer(cstr))
	// defer C.free(unsafe.Pointer(cstr))
	// return (*C.PyObject)(unsafe.Pointer(cstr))

	// cstr := C.CString("Pee")
	// uptr := unsafe.Pointer(cstr), "s"

	// return (*C.PyObject)(unsafe.Pointer(python.Py_BuildValue("i", 1)))

	return C.PyNone()

	// Attempts to return things to C code
	// This may be a solution: https://docs.python.org/2/c-api/arg.html#c.Py_BuildValue
	s := python.PyString_FromString("asdf")
	unptr := (*C.PyObject)(unsafe.Pointer(&s))
	defer C.free(unsafe.Pointer(unptr))
	return unptr

	if ret, err := topy(results); err != nil {
		fmt.Println("GO: could not convert types back to python. Name:", target, err)
	} else {
		// This pointer doesnt work... but is it coming from the slice itself, or internal pointers?
		unptr := (*C.PyObject)(unsafe.Pointer(&ret))
		defer C.free(unsafe.Pointer(unptr))
		return unptr
	}

	return C.PyNone()
}

//
// threading
//

//export createThreadCallbacke
func createThreadCallbacke() {
	C.register_sig_handler()
	C.pthread_setcanceltype(C.PTHREAD_CANCEL_ASYNCHRONOUS, nil)
	(<-create_callback)()
}

// calls C's sleep function
func PtSleep(seconds uint) {
	C.sleep(C.uint(seconds))
}

// initializes a thread using pthread_create
func PtCreate(cb ThreadCallback) Thread {
	var pid C.pthread_t
	pidptr := &pid
	create_callback <- cb

	C.createThread(pidptr)

	return Thread(uintptr(unsafe.Pointer(&pid)))
}

// determines if the thread is running
func (t Thread) Running() bool {
	// magic number "3". oops
	return int(C.pthread_kill(t.c(), 0)) != 3
}

func (t Thread) Kill() {
	C.pthread_kill(t.c(), C.SIGSEGV)
}

// helper function to convert the Thread object into a C.pthread_t object
func (t Thread) c() C.pthread_t {
	return *(*C.pthread_t)(unsafe.Pointer(t))
}
