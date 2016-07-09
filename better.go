package gosnake

/*
#cgo pkg-config: python-2.7
#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <pthread.h>
#include <unistd.h>
#include <stdio.h>

extern void createThreadCB();

static PyObject* thread_callback() {
    PyObject *_module_name, *_module;
    PyGILState_STATE _gstate;

    // Initialize python GIL state
    _gstate = PyGILState_Ensure();

    // Now execute some python code (call python functions)
    _module_name = PyString_FromString("json_dump");
    _module = PyImport_Import(_module_name);

    // Call a method of the class with no parameters
    PyObject *_attr, *_result;
    _attr = PyObject_GetAttr(_module, PyString_FromString("run"));
    _result = PyObject_CallObject(_attr, NULL);

    // Clean up
    Py_DECREF(_module);
    Py_DECREF(_module_name);
    Py_DECREF(_attr);

    PyGILState_Release(_gstate);
    return _result;
}

static void createThreade(pthread_t* pid) {
    pthread_create(pid, NULL, (void*)createThreadCB, pid);
}

static void initialize_pythone () {
    if (Py_IsInitialized() == 0) {
        Py_Initialize();
    }

    // make sure the GIL is correctly initialized
    if (PyEval_ThreadsInitialized() == 0) {
        PyEval_InitThreads();
    }

    PyEval_ReleaseThread(PyGILState_GetThisThreadState());
}

*/
import "C"

import (
	"encoding/json"
	"fmt"
	"sync"
)

var cbLock *sync.Mutex = &sync.Mutex{}
var callbacks map[*C.pthread_t]ThreadCB = map[*C.pthread_t]ThreadCB{}

type ThreadCB func(a *C.PyObject)

func init() {
	C.initialize_pythone()
}

//export createThreadCB
func createThreadCB(pid *C.pthread_t) {
	cbLock.Lock()

	if _cb, _ok := callbacks[pid]; !_ok {
		panic(fmt.Errorf("failed to found thread callback for `%v`", pid))
	} else {
		delete(callbacks, pid)
		cbLock.Unlock()

		_result := C.thread_callback()
		_cb(_result)
	}
}

func create_thread(cb ThreadCB) {
	_pid := new(C.pthread_t)

	cbLock.Lock()
	callbacks[_pid] = cb
	cbLock.Unlock()

	C.createThreade(_pid)
}

func BetterTest() {
	var _wg sync.WaitGroup
	// _wg.Add(n)

	for i := 0; i < 500; i++ {
		go create_thread(func(result *C.PyObject) {
			defer _wg.Done()

			_result_string := C.GoString(C.PyString_AsString(result))
			fmt.Printf("< got result string: %v (%T)\n", _result_string, _result_string)

			var _parsed []interface{}
			if _err := json.Unmarshal([]byte(_result_string), &_parsed); _err != nil {
				panic(fmt.Errorf("got invalid result from python function, `%v`\n", _result_string))
			}

		})
	}
	_wg.Wait()
}
