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
	r, e := module.Call("callee_none_none")

	assert.Nil(t, e)
	assert.Equal(t, nil, r)
}

func TestPyThreeNone(t *testing.T) {
	module, _ := Import("testmodule")
	r, e := module.Call("callee_three_none", "joe", 12, false)

	assert.Nil(t, e)
	assert.Equal(t, nil, r)
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
func testexport() {}

func TestSuccessfulExport(t *testing.T) {
	err := Export(testexport)
	assert.Nil(t, err)
}

func TestGoNoneNone(t *testing.T) {
	Export(testexport)

	module, _ := Import("testmodule")
	r, e := module.Call("reflect_call", "testexport")

	assert.Nil(t, e)
	assert.Nil(t, r)
}
