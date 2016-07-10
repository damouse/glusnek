package gosnake

//#cgo pkg-config: python-2.7
//#include "binding.h"
import "C"

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Global map of exported functions. Move this to the Module object
var exports map[string]interface{} = map[string]interface{}{}

// Called from python through the C code in binding.h. Returns the result across the binding
//export _gosnake_invoke
func _gosnake_invoke(self *C.PyObject, args *C.PyObject) *C.char {
	// Building data structures ourselves and sending them to go can be... a little icky
	// For now everything goes back as json
	a := []interface{}{"hello"}
	b, _ := json.Marshal(a)
	ret := string(b)

	cmode := C.CString(ret)
	return cmode
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
	// todo: deal with this err
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

	// Method above returns functions in the form :  main.foo
	parts := strings.Split(name, ".")

	if len(parts) < 2 {
		return fmt.Errorf("Cant resolve function name for exporting")
	}

	ending := parts[len(parts)-1]
	exports[ending] = fn

	return nil
}
