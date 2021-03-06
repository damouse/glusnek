package gosnake

//#cgo pkg-config: python-2.7
//#include "binding.h"
import "C"

import (
	"encoding/json"
	"fmt"
	"sync"
	"unsafe"

	"github.com/damouse/cumin"
	"github.com/sbinet/go-python"
)

// Global map of exported functions. Move this to the Module object
var exportLock = &sync.RWMutex{}
var exports map[string]*cumin.Curry = map[string]*cumin.Curry{}

// Called from python through the C code in binding.h. Returns the result across the binding
//export _gosnake_invoke
func _gosnake_invoke(self *C.PyObject, args *C.PyObject) *C.char {
	transform := python.PyObject_FromVoidPtr(unsafe.Pointer(args))

	if unpacked, e := togo(transform); e != nil {
		fmt.Println("GO: an error occured! ", e)

	} else if arr, ok := unpacked.([]interface{}); !ok {
		fmt.Println("GO: an error occured! ")

	} else if len(arr) < 1 {
		fmt.Println("GO: not enough arguments to gocall")

	} else if target, ok := arr[0].(string); !ok {
		fmt.Println("GO: method name wasnt a string")

	} else if curry, ok := getExport(target); !ok {
		fmt.Println("GO: no function exported as ", target)

	} else if results, err := curry.Invoke(arr[1:]); err != nil {
		fmt.Println("GO: exported function erred. ", target, err, arr[1:])

	} else if b, err := json.Marshal(results); err != nil {
		// Returns can be icky, so we're sticking to json for now
		fmt.Printf("Unable to marshal results: %v", results)

	} else {
		return C.CString(string(b))
	}

	panic("fatal error")
}

// Exports a go function to python. This must be an unbound, top-level function, not one with a
// receiver or an anonymous function. The name of the function in Go becomes the name of the
// function in python.
//
// Example (go):
//      func mygolangfunction() { }
//      gosnake.Export(mygolangfunction)
//
// python:
//      gosnake.gocall("mygolangfunction")
//
// Currently panics if the passed function is a method or anonymous
func Export(fn interface{}) error {
	if curry, err := cumin.NewCurry(fn); err != nil {
		return err
	} else {
		exports[curry.Name()] = curry
		return nil
	}
}

// Utility method. Grab an export using the lock
func getExport(target string) (*cumin.Curry, bool) {
	exportLock.RLock()
	defer exportLock.RUnlock()

	for name, f := range exports {
		if name == target {
			return f, true
		}
	}

	return nil, false
}
