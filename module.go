package gosnake

import (
	"fmt"

	"github.com/sbinet/go-python"
)

// Models to expose interface

type Operation struct {
	module *Module // previously imported module
	target string  // the name of the target function

	args []interface{}

	returnChan chan interface{}
	errChan    chan error
}

// Representation of a python module
type Module struct {
	name string
}

// Python equivalent: import name
func Import(name string) (*Module, error) {
	// Ensure we're initialized
	initializeBinding()

	// Make sure the module actually exists
	if err := tryImport(name); err != nil {
		return nil, err
	}

	m := &Module{
		name: name,
	}

	return m, nil
}

// Python equivalent: module.function(args)
// If the call raises an exception in python its passed along as a go error
func (b *Module) Call(function string, args ...interface{}) (interface{}, error) {
	op := &Operation{
		module:     b,
		target:     function,
		args:       args,
		returnChan: make(chan interface{}),
		errChan:    make(chan error),
	}

	// Enqueue the operation
	opChan <- op

	select {
	case e := <-op.errChan:
		return nil, e
	case r := <-op.returnChan:
		return r, nil
	}
}

// Try to import python module name
func tryImport(name string) error {
	gil := python.PyGILState_Ensure()
	defer python.PyGILState_Release(gil)

	if m := python.PyImport_ImportModule(name); m == nil {
		return fmt.Errorf("Could not find a python module named \"%s\"", name)
	} else {
		m.DecRef()
		return nil
	}
}
