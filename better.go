package gosnake

/*
ToDo to make sure this actually works
    - Exporting functions into python
    - Getting results from python function


Advanced
    - "Permanent", goroutine-safe imports
    - Performance Cleanup
        - Dynamically create threads as needed to handle requests
        - Pre-create goroutines for outbound
    - Objects for modules
*/

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
inline PyObject *PYNONE() {Py_INCREF(Py_None); return Py_None;}

// Receiver function in Go. Cant be called directly from python since we cannot
// return the *PyObject pointer without violating the cgo rules
extern char *_gosnake_invoke(PyObject *self, PyObject *args);

static PyObject* _gosnake_receive(PyObject *self, PyObject *args) {
    PyObject *_module_name, *_module, *fn_name;
    PyGILState_STATE _gstate;

    // route the call into Go
    char *result;
    PyObject *pyStringResult;
    result = _gosnake_invoke(self, args);
    pyStringResult = PyString_FromString(result);

    // Initialize python GIL state. I think we don't have to do this?
    // _gstate = PyGILState_Ensure();

    // import json
    _module_name = PyString_FromString("json");
    _module = PyImport_Import(_module_name);

    // import json.dumps
    PyObject *_attr, *_result;
    fn_name = PyString_FromString("loads");
    _attr = PyObject_GetAttr(_module, fn_name);
    _result = PyObject_CallObject(_attr, pyStringResult);

    // Check for error message
    if (PyErr_Occurred()) {
        PyObject *ptype, *pvalue, *ptraceback;
        PyErr_Fetch(&ptype, &pvalue, &ptraceback);

        char *pStrErrorMessage = PyString_AsString(pvalue);
        fprintf(stdout, "gosnake: Could not parse return json: %s\n", pStrErrorMessage);

        Py_DECREF(ptype);
        Py_DECREF(pvalue);
        Py_DECREF(ptraceback);
    }

    // Clean up
    Py_DECREF(_module);
    Py_DECREF(_module_name);
    Py_DECREF(fn_name);
    Py_DECREF(_attr);
    Py_DECREF(pyStringResult);

    // PyGILState_Release(_gstate);

    // return _result;
    return PYNONE();
}

static PyMethodDef GosnakeMethods[] = {
    {"gocall", _gosnake_receive, METH_VARARGS, "doc string"},
    {NULL},
};

// Set up the python environment and inject gosnake as a virtual-python module
static void pyinit (int log) {
    if (Py_IsInitialized() == 0) {
        Py_Initialize();
    }

    if (PyEval_ThreadsInitialized() == 0) {
        PyEval_InitThreads();
    }

    Py_InitModule("gosnake", GosnakeMethods);
    PyEval_ReleaseThread(PyGILState_GetThisThreadState());

    if (log != 0) {
        fprintf(stdout, "gosnake: initialized python env\n");
    }
}
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/liamzdenek/go-pthreads"
	"github.com/sbinet/go-python"
)

// Thread Pool
var POOL_SIZE int = 10
var pool []*pthread.Thread

// Channel by which threads receive operations
var opChan chan *Operation

type Operation struct {
	args       string
	returnChan chan string
}

func Initialize() {
	C.pyinit(C.int(1))

	opChan = make(chan *Operation)
	pool = make([]*pthread.Thread, POOL_SIZE)

	for i := 0; i < POOL_SIZE; i++ {
		t := pthread.Create(threadConsume)
		pool = append(pool, &t)
	}

	// Export our new module methods
}

// Consume an operation from the queue
func threadConsume() {
	for {
		// Grab a new operation
		op := <-opChan

		// Lock the gil
		gil := python.PyGILState_Ensure()

		// Import the module, target function
		m := python.PyImport_ImportModule("adder")
		fn := m.GetAttrString("run")

		m.IncRef()
		fn.IncRef()

		// Pack the arguments
		args := python.PyTuple_New(2)
		a1 := python.PyString_FromString(op.args)
		a2 := python.PyInt_FromLong(1234)
		python.PyTuple_SET_ITEM(args, 0, a1)
		python.PyTuple_SET_ITEM(args, 1, a2)

		// Retain the arguments
		args.IncRef()
		a1.IncRef()
		a2.IncRef()

		// Call into Python
		ret := fn.CallObject(args)

		// Deserialize
		resultString := python.PyString_AsString(ret)

		// Cleanup
		ret.DecRef()
		args.DecRef()
		m.DecRef()
		fn.DecRef()
		args.DecRef()
		a1.DecRef()
		a2.DecRef()

		// Release the Gil
		python.PyGILState_Release(gil)

		op.returnChan <- resultString
	}
}

// An invocation FROM python to go
// TODO: throw errors back to python!
//export _gosnake_invoke
func _gosnake_invoke(self *C.PyObject, args *C.PyObject) *C.char {
	fmt.Println("GO: invocation", self, args)

	ret := "[Returning from go!]"

	cmode := C.CString(ret)
	defer C.free(unsafe.Pointer(cmode))
	return cmode
}

// We can likely pool the pthreads, but that sounds like a problem for future mickey
// Should we check if threads are "busy"?
func Call(args string) string {
	op := &Operation{
		args:       args,
		returnChan: make(chan string),
	}

	opChan <- op

	ret := <-op.returnChan
	return ret
}

func BetterTest(routines int, iterations int) {
	wg := sync.WaitGroup{}
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func(gid int) {
			for j := 0; j < iterations; j++ {

				ret := Call("Hey boyo")

				fmt.Printf("(%d  %d) \t%s\n", gid, j, ret)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("\nInternal Done")
}
