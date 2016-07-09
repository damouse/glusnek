package gosnake

/*
#cgo pkg-config: python-2.7
#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include <signal.h>
#include <unistd.h>
#include <stdio.h>

static void c_call_python() {
    PyObject *_module_name, *_module;
    PyGILState_STATE _gstate;

    // Initialize python GIL state
    _gstate = PyGILState_Ensure();

    // Now execute some python code (call python functions)
    _module_name = PyString_FromString("adder");
    _module = PyImport_Import(_module_name);

    // Call a method of the class with no parameters
    PyObject *_attr, *_result;
    _attr = PyObject_GetAttr(_module, PyString_FromString("run"));
    _result = PyObject_CallObject(_attr, NULL);

    // Clean up
    Py_DECREF(_module);
    Py_DECREF(_module_name);
    Py_DECREF(_attr);
    Py_DECREF(_result);

    PyGILState_Release(_gstate);

    // return _result;
}
*/
import "C"

import (
	"fmt"
	"sync"

	"github.com/liamzdenek/go-pthreads"
	"github.com/sbinet/go-python"
)

// The pthreads live in a pool. They should be sent messages through the opChan channel
var POOL_SIZE int = 5
var pool []*pthread.Thread
var opChan chan func()

func Initialize() {
	opChan = make(chan func())
	pool = make([]*pthread.Thread, POOL_SIZE)

	for i := 0; i < POOL_SIZE; i++ {
		t := pthread.Create(threadConsume)
		pool = append(pool, &t)
	}
}

// Consume an operation from the queue
func threadConsume() {
	for {
		<-opChan
		// C.c_call_python()
		PyCall()
	}
}

// We can likely pool the pthreads, but that sounds like a problem for future mickey
func PyCall() string {
	gil := python.PyGILState_Ensure()

	m := python.PyImport_ImportModule("adder")
	fn := m.GetAttrString("run")

	m.IncRef()
	fn.IncRef()

	// Pack the arguments
	// args := python.PyTuple_New(2)
	// python.PyTuple_SET_ITEM(args, 0, python.PyString_FromString("Hello!"))
	// python.PyTuple_SET_ITEM(args, 1, python.PyInt_FromLong(1234))

	args := python.PyTuple_New(0)
	args.IncRef()

	ret := fn.CallObject(args)

	// Release everything
	ret.DecRef()
	args.DecRef()
	m.DecRef()
	fn.DecRef()

	python.PyGILState_Release(gil)

	return "Done"
}

func BetterTest() {
	n := 1000
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(gid int) {
			for j := 0; j < 1000; j++ {
				opChan <- func() {}
				fmt.Printf("gorutine: %d iteration: %d\n", gid, j)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("\nInternal Done")
}

// func checkPyErr() {
// 	if python.PyErr_CheckSignals() {
// 		python.PyErr_PrintEx(false)
// 	}
// }

// func BetterTest() {
// 	n := 10
// 	wg := sync.WaitGroup{}
// 	wg.Add(n)

// 	for i := 0; i < n; i++ {
// 		go func(gid int) {
// 			for j := 0; j < 1000; j++ {
// 				PyCall()
// 				fmt.Printf("gorutine: %d iteration: %d\n", gid, j)

// 				// var retjson []interface{}
// 				// if _err := json.Unmarshal([]byte(rs), &retjson); _err != nil {
// 				// 	panic(fmt.Errorf("got invalid result from python function, `%v`\n", rs))
// 				// }

// 				// fmt.Printf("gorutine: %d iteration: %d pid: %f\n", gid, j, retjson[0])
// 			}

// 			wg.Done()
// 		}(i)
// 	}

// 	wg.Wait()
// 	fmt.Println("\nInternal Done")
// }
