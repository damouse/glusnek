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
