package gosnake

//#cgo pkg-config: python-2.7
//#include "binding.h"
import "C"

import "encoding/json"

// Exported go functions to python
// Go functions exposed to python
type ExportedFunction func([]interface{}, map[string]interface{}) ([]interface{}, error)

var exports map[string]*ExportedFunction

// An invocation FROM python to go
// TODO: raise errors from the go invocation back to python
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
