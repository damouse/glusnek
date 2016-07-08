package gosnake

// WARN: PYTHONPATH must be set to demo before tests can run!

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWork(t *testing.T) {
	InitPyEnv()
	RunTest(2)
	assert.Nil(t, nil)
}
