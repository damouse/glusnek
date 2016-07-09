package gosnake

// WARN: PYTHONPATH must be set to demo before tests can run!

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	b := NewBinding()
	e := b.Import("adder")
	assert.Nil(t, e)
}

//
// Go -> Py
//
func TestCallPy(t *testing.T) {
	b := NewBinding()
	b.Import("adder")

	r, e := b.Call("adder", "birthday", "bill", 15)

	assert.Nil(t, e)
	cast := r.([]interface{})
	assert.Equal(t, "bill", cast[0].(string))
	assert.Equal(t, 15, cast[1].(int))
}

//
// Py -> Go
//
// func TestCallGo(t *testing.T) {
// 	b := NewBinding()
// 	b.Import("adder")

// 	// callback := func()

// 	r, e := b.Call("adder", "callback", "bill", 15)

// 	assert.Nil(t, e)
// 	cast := r.([]interface{})
// 	assert.Equal(t, "bill", cast[0].(string))
// 	assert.Equal(t, 15, cast[1].(int))
// }
