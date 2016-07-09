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

static PyObject* thread_callback() {
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

    PyGILState_Release(_gstate);
    return _result;
}
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/liamzdenek/go-pthreads"
	"github.com/sbinet/go-python"
)

func create_thread(cb func(a *C.PyObject)) {

	// C.createThreade(_pid)
	done := make(chan error)

	t := pthread.Create(func() {
		// Lets try it with the gopython code...
		gil := C.PyGILState_Ensure()
		defer C.PyGILState_Release(gil)

		m := python.PyImport_ImportModule("adder")
		fn := m.GetAttrString("run")

		// Pack the arguments
		args := python.PyTuple_New(2)
		python.PyTuple_SET_ITEM(args, 0, python.PyString_FromString("Hello!"))
		python.PyTuple_SET_ITEM(args, 1, python.PyInt_FromLong(1234))
		ret := fn.CallObject(args)

		resultAsPyobj := (*C.PyObject)(unsafe.Pointer(ret))
		cb(resultAsPyobj)

		// result := C.thread_callback()
		// cb(result)

		done <- nil
	})

	<-done
	t.Kill()
}

func BetterTest() {
	for i := 0; i < 500; i++ {
		go create_thread(func(result *C.PyObject) {
			for i := 0; i < 500; i++ {

				_result_string := C.GoString(C.PyString_AsString(result))
				fmt.Printf("< got result string: %v (%T)\n", _result_string, _result_string)

				var _parsed []interface{}
				if _err := json.Unmarshal([]byte(_result_string), &_parsed); _err != nil {
					panic(fmt.Errorf("got invalid result from python function, `%v`\n", _result_string))
				}
			}

		})
	}
}
