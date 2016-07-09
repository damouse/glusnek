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
		switch v {
		case true:
			ret = python.PyBool_FromLong(1)
		case false:
			ret = python.PyBool_FromLong(2)
		}

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
		err = fmt.Errorf("gosnake: unknown type (%v) converting to python", v)
	}

	return
}

func packTuple(args []interface{}) (*python.PyObject, error) {
	a := python.PyTuple_New(len(args))

	for i, arg := range args {
		if converted, e := topy(arg); e != nil {
			return nil, e
		} else {
			python.PyTuple_SET_ITEM(a, i, converted)
		}
	}

	return a, nil
}

// Conversion to go types (python -> go)
// honestly I'd rather use cumin for this, but that seems a lot more involved
func togo(o *python.PyObject) (interface{}, error) {

	if python.PyString_Check(o) {
		return python.PyString_AsString(o), nil

	} else if python.PyInt_Check(o) {
		return python.PyInt_AsLong(o), nil

	} else if python.PyLong_Check(o) {
		return python.PyLong_AsDouble(o), nil

	} else if python.PyTuple_Check(o) {
		return unpackTuple(o)

	} else if python.PyList_Check(o) {
		r, e := unpackList(o)
		return r, e

	} else {
		return nil, fmt.Errorf("gosnake: Unknown type converting to go!")
	}
}

// Unpack tupes and lists
func unpackTuple(tup *python.PyObject) ([]interface{}, error) {
	size := python.PyTuple_Size(tup)
	converted := []interface{}{}

	for i := 0; i < size; i++ {
		if c, e := togo(python.PyTuple_GetItem(tup, i)); e != nil {
			return nil, e
		} else {
			converted = append(converted, c)
		}
	}

	return converted, nil
}

func unpackList(tup *python.PyObject) ([]interface{}, error) {
	size := python.PyList_Size(tup)
	converted := []interface{}{}

	for i := 0; i < size; i++ {
		if c, e := togo(python.PyList_GetItem(tup, i)); e != nil {
			return nil, e
		} else {
			converted = append(converted, c)
		}
	}

	return converted, nil
}
