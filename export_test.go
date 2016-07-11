package gosnake

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func export() {}

func TestSuccessfulExport(t *testing.T) {
	err := Export(export)
	Nil(t, err)
}

func TestGoNoneNone(t *testing.T) {
	Export(export)

	module, _ := Import("testmodule")
	r, e := module.Call("reflect_call", "export")

	Nil(t, e)
	Nil(t, r)
}

func singleReturn(a int) int {
	return a
}

func TestGoSingleReturn(t *testing.T) {
	Export(singleReturn)

	module, _ := Import("testmodule")
	r, e := module.Call("reflect_call", "singleReturn", 1)

	// This cast is *not correct*. We should receive a single int instead of a slice
	// if returned an int
	badResults := r.([]interface{})

	Nil(t, e)
	Equal(t, 1, badResults[0].(int))
}
