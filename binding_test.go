package gosnake

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestSuccessfulImport(t *testing.T) {
	module, err := Import("testmodule")

	Nil(t, err)
	Equal(t, "testmodule", module.name)
	Equal(t, _INITIALIZED, true)
}

func TestBadImport(t *testing.T) {
	_, err := Import("idontexist")
	NotNil(t, err)
}

//
// Test names are in the following format: Test{TargetLanguage}{Arguments}{Return}
//
func TestPyNoneNone(t *testing.T) {
	module, _ := Import("testmodule")
	r, e := module.Call("callee_none_none")

	Nil(t, e)
	Equal(t, nil, r)
}

func TestPyThreeNone(t *testing.T) {
	module, _ := Import("testmodule")
	r, e := module.Call("callee_three_none", "joe", 12, false)

	Nil(t, e)
	Equal(t, nil, r)
}

func TestPyNoneOne(t *testing.T) {
	module, _ := Import("testmodule")
	r, e := module.Call("callee_none_one")

	Nil(t, e)
	Equal(t, "higo", r.(string))
}

//
// Py -> Go
//
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

// func singleReturn(a int) int {
// 	return a
// }

// func TestGoSingleReturn(t *testing.T) {
// 	Export(singleReturn)

// 	module, _ := Import("testmodule")
// 	r, e := module.Call("reflect_call", "singleReturn", 1)

// 	// This cast is *not correct*. We should receive a single int instead of a slice
// 	// if returned an int
// 	badResults := r.([]interface{})

// 	Nil(t, e)
// 	Equal(t, 1, badResults[0].(int))
// }
