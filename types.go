package gosnake

// Conversion between C.Python and python.PyObject types to and from G

/*
   cgo type mappings. Taken from http://blog.giorgis.io/cgo-examples

   char -->  C.char -->  byte
   signed char -->  C.schar -->  int8
   unsigned char -->  C.uchar -->  uint8
   short int -->  C.short -->  int16
   short unsigned int -->  C.ushort -->  uint16
   int -->  C.int -->  int
   unsigned int -->  C.uint -->  uint32
   long int -->  C.long -->  int32 or int64
   long unsigned int -->  C.ulong -->  uint32 or uint64
   long long int -->  C.longlong -->  int64
   long long unsigned int -->  C.ulonglong -->  uint64
   float -->  C.float -->  float32
   double -->  C.double -->  float64
   wchar_t -->  C.wchar_t  -->  [[https://github.com/orofarne/gowchar/blob/master/gowchar.go][wchar]]
   void * -> unsafe.Pointer
*/

// #cgo pkg-config: python-2.7
// #include "Python.h"
import "C"

import (
	"fmt"

	"github.com/sbinet/go-python"
)

// Conversion to python types (go -> python)
func topy(v interface{}) (ret *python.PyObject, err error) {
	// return (*C.PyObject)(unsafe.Pointer(self))

	switch v := v.(type) {
	case bool:
		val := 1
		if v {
			val = 0
		}

		ret = python.PyBool_FromLong(val)

	case float32:
		ret = python.PyFloat_FromDouble(v)

	case float64:
		ret = python.PyLong_FromDouble(v)

	case int:
		ret = python.PyInt_FromLong(v)

	case string:
		ret = python.PyString_FromString(v)

	case *python.PyObject:
		ret = v

	default:
		err = fmt.Errorf("python: unknown type (%v)", v)
	}

	return
}

// Conversion to go types (python -> go)
func togo(o *python.PyObject) (interface{}, error) {
	if python.PyString_Check(o) {
		return python.PyString_AsString(o), nil

	} else if python.PyInt_Check(o) {
		return python.PyInt_AsLong(o), nil

	} else {
		return nil, fmt.Errorf("Unknown type converting to go!")
	}
}
