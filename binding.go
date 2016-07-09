package gosnake

// Binding between go and python, with pthreading implmementation

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

var lock sync.Mutex

type Thread uintptr
type ThreadCallback func()

var create_callback chan ThreadCallback

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
	return &Binding{
		lock:           &sync.Mutex{},
		thread:         0,
		threadCallback: make(chan ThreadCallback, 1),
		imports:        map[string]*python.PyObject{},
		exports:        map[string]*ExportedFunction{},
	}
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

	if m, ok := b.imports[module]; !ok {
		return nil, fmt.Errorf("Cant call python function, module %s has not been imported. Did you call Import(moduleName)?", module)

	} else if fn := m.GetAttrString(function); m == fn {
		return nil, b.parseException()

	} else {

		b.lock.Lock()
		defer b.lock.Unlock()

		errChan := make(chan error)
		defer close(errChan)
		resultChan := make(chan interface{})
		defer close(resultChan)

		thread := PtCreate(func() {
			gil := C.PyGILState_Ensure()
			defer C.PyGILState_Release(gil)

			if tup, e := packTuple(args); e != nil {
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

//export pyInvocation
func pyInvocation(self *C.PyObject, args *C.PyObject) *C.PyObject {
	fmt.Println("GO: invocation", self, args)

	// a := python.PyObject_FromVoidPtr(unsafe.Pointer(args))
	// iter := python.PySeqIter_New(a)

	// converted := []interface{}{}

	// for i := 0; i < python.PyTuple_Size(iter); i++ {
	// 	p := python.PyTuple_GetItem(iter, i)
	// 	s := python.PyString_AsString(p)
	// 	converted = append(converted, s)
	// }

	// fmt.Printf("Go code called!: %s\n", converted)

	return C.PyNone()
}

//
// threading
//

//export createThreadCallback
func createThreadCallback() {
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

// signals the thread in question to terminate
func (t Thread) Kill() {
	C.pthread_kill(t.c(), C.SIGSEGV)
}

// helper function to convert the Thread object into a C.pthread_t object
func (t Thread) c() C.pthread_t {
	return *(*C.pthread_t)(unsafe.Pointer(t))
}
