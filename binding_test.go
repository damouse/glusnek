package gosnake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessfulImport(t *testing.T) {
	module, err := Import("testmodule")

	assert.Nil(t, err)
	assert.Equal(t, "testmodule", module.name)
	assert.Equal(t, _INITIALIZED, true)
}

func TestBadImport(t *testing.T) {
	_, err := Import("idontexist")
	assert.NotNil(t, err)
}

//
// Test names are in the following format: Test{TargetLanguage}{Arguments}{Return}
//
func TestPyNoneNone(t *testing.T) {
	module, _ := Import("testmodule")
	_, e := module.Call("callee_none_none")

	assert.Nil(t, e)
	// TODO: return types
}

func TestPyThreeNone(t *testing.T) {
	module, _ := Import("testmodule")
	_, e := module.Call("callee_three_none", "joe", 12, false)

	assert.Nil(t, e)
	// TODO: return types
}

func TestPyNoneOne(t *testing.T) {
	module, _ := Import("testmodule")
	r, e := module.Call("callee_none_one")

	assert.Nil(t, e)
	assert.Equal(t, "higo", r.(string))
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
