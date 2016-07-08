package gosnake

// WARN: PYTHONPATH must be set to demo before tests can run!

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	b := NewBinding()
	e := b.Import("adder")
	assert.Nil(t, e)
}

func TestCall(t *testing.T) {
	b := NewBinding()

	b.Import("adder")
	r, e := b.Call("adder", "birthday", "bill", 15)

	fmt.Println(r)
	assert.Nil(t, e)
}
