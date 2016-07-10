package gosnake

//#cgo pkg-config: python-2.7
//#include "binding.h"
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/damouse/cumin"
	"github.com/sbinet/go-python"
)

// TODO: locked map!
// Global map of exported functions. Move this to the Module object
var exports map[string]*cumin.Curry = map[string]*cumin.Curry{}

// Called from python through the C code in binding.h. Returns the result across the binding
//export _gosnake_invoke
func _gosnake_invoke(self *C.PyObject, args *C.PyObject) *C.char {
	// Unmarshal arguments
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
	var curry *cumin.Curry
	for name, f := range exports {
		if name == target {
			fn = f
		}
	}

	if curry == nil {
		fmt.Println("GO: no function exported as ", target)
	}

	results, err := curry.Invoke(goArgs)

	if err != nil {
		panic(err)
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
