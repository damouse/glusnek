package gosnake

// WARN: PYTHONPATH must be set to demo before tests can run!

import (
	"testing"

	"github.com/sbinet/go-python"
	"github.com/stretchr/testify/assert"
)

//
// Python -> Go
func TestStringToGo(t *testing.T) {
	o := "there once was a man from peru"
	c, e := togo(python.PyString_FromString(o))

	assert.Nil(t, e)
	assert.Equal(t, o, c.(string))
}

func TestIntToGo(t *testing.T) {
	o := 123
	c, e := togo(python.PyInt_FromLong(o))

	assert.Nil(t, e)
	assert.Equal(t, o, c.(int))
}

func TestBoolToGo(t *testing.T) {
	o := 123
	c, e := togo(python.PyInt_FromLong(o))

	assert.Nil(t, e)
	assert.Equal(t, o, c.(int))
}

func TestFloatToGo(t *testing.T) {
	o := float64(1234)
	c, e := togo(python.PyLong_FromDouble(o))

	assert.Nil(t, e)
	assert.Equal(t, o, c.(float64))
}

func TestTupleToGo(t *testing.T) {
	tup := python.PyTuple_New(2)
	python.PyTuple_SET_ITEM(tup, 0, python.PyString_FromString("asdf"))
	python.PyTuple_SET_ITEM(tup, 1, python.PyInt_FromLong(12345))

	c, e := togo(tup)

	assert.Nil(t, e)
	cast := c.([]interface{})
	assert.Equal(t, "asdf", cast[0].(string))
	assert.Equal(t, 12345, cast[1].(int))
}

func TestListToGo(t *testing.T) {
	tup := python.PyList_New(2)

	python.PyList_SetItem(tup, 0, python.PyString_FromString("asdf"))
	python.PyList_SetItem(tup, 1, python.PyInt_FromLong(12345))

	c, e := togo(tup)

	assert.Nil(t, e)
	cast := c.([]interface{})
	assert.Equal(t, "asdf", cast[0].(string))
	assert.Equal(t, 12345, cast[1].(int))
}

func TestNoneToGo(t *testing.T) {
	o := python.Py_None
	c, e := togo(o)

	assert.Nil(t, e)
	assert.Equal(t, c, nil)
}

// Dictionaries

//
// Go -> Python
//
func TestIntToPy(t *testing.T) {
	o := 123
	c, e := topy(o)

	assert.Nil(t, e)
	assert.True(t, python.PyInt_Check(c))
}

func TestBoolToPy(t *testing.T) {
	o := true
	c, e := topy(o)

	assert.Nil(t, e)
	assert.True(t, python.PyBool_Check(c))
}

func TestStringToPy(t *testing.T) {
	c, e := topy("horses")
	assert.Nil(t, e)
	assert.True(t, python.PyString_Check(c))
}

func TestFloatToPy(t *testing.T) {
	c, e := topy(float32(1234))
	assert.Nil(t, e)
	assert.True(t, python.PyFloat_Check(c))
}

func TestLongToPy(t *testing.T) {
	c, e := topy(float64(1234))
	assert.Nil(t, e)
	assert.True(t, python.PyLong_Check(c))
}

func TestNoneToPy(t *testing.T) {
	c, e := topy(nil)

	assert.Nil(t, e)
	assert.True(t, isNone(c))
}

// func TestListToPy(t *testing.T) {
// 	c, e := topy([]string{"a", "b"})

// 	assert.Nil(t, e)
// 	assert.True(t, python.PyList_Check(c))
// }

// Arrays
// Dictionaries
