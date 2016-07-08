package gosnake

// Binding between go and python, with pthreading implmementation

/*
#cgo pkg-config: python-2.7

#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include <signal.h>
#include <unistd.h>
#include <stdio.h>

// Convenience method for returning None
inline PyObject *PyNone() {Py_INCREF(Py_None); return Py_None;}

// A call from python to go. This is implemented below in go
extern PyObject* PythonInvocation(PyObject *self, PyObject *args);

// Part of the threading implmentation
extern void createThreadCallback();
static void sig_func(int sig);

// Set up the python environment and inject gosnake as a virtual-python module
static void initialize_python () {
    if (Py_IsInitialized() == 0) {
        Py_Initialize();
    }

    if (PyEval_ThreadsInitialized() == 0) {
        PyEval_InitThreads();
    }

    static PyMethodDef ModuleMethods[] = {
        {"gocall", PythonInvocation, METH_VARARGS, "doc string"},
        {NULL},
    };

    Py_InitModule("gosnake", ModuleMethods);
    PyEval_ReleaseThread(PyGILState_GetThisThreadState());
    fprintf(stdout, "gosnake: initialized python env\n");
}

// Threading implementation
static void createThread(pthread_t* pid) {
    pthread_create(pid, NULL, (void*)createThreadCallback, NULL);
}

static void sig_func(int sig) {
    fprintf(stdout, "gosnake: handling exit signal\n");
    signal(SIGSEGV,sig_func);
    pthread_exit(NULL);
}

static void register_sig_handler() {
    signal(SIGSEGV,sig_func);
}
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/sbinet/go-python"
)

var lock sync.Mutex

type Thread uintptr
type ThreadCallback func()

var create_callback chan ThreadCallback

func InitPyEnv() {
	C.initialize_python()
	C.register_sig_handler()
	create_callback = make(chan ThreadCallback, 1)
}

//export PythonInvocation
func PythonInvocation(self *C.PyObject, args *C.PyObject) *C.PyObject {
	fmt.Println("GO: invocation", args)

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

// An example of publically exposing a method to python
//export port
func port(self *C.PyObject, args *C.PyObject) *C.PyObject {
	fmt.Println("GO: public exported called")
	return nil
}

func testFunctionTypes(name string, age int) {
	_module := python.PyImport_ImportModuleNoBlock("adder")
	attr := _module.GetAttrString("birthday")

	a := python.PyTuple_New(2)
	python.PyTuple_SET_ITEM(a, 0, python.PyString_FromString(name))
	python.PyTuple_SET_ITEM(a, 1, python.PyInt_FromLong(age))

	ret := attr.CallObject(a)

	// Python threw an exception! Return an error here to the go caller?
	if ret == nil {
		python.PyErr_PrintEx(false)
	}

	fmt.Println("GO: Done", ret)
}

func RunTest(num int) {
	lock.Lock()
	defer lock.Unlock()
	done := make(chan bool)

	thread := PtCreate(func() {
		gil := C.PyGILState_Ensure()
		defer C.PyGILState_Release(gil)

		testFunctionTypes("joe", num)

		done <- true
	})

	defer thread.Kill()
	<-done
	close(done)
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
