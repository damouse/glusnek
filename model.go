package gosnake

// Models to expose interface

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
		returnChan: make(chan string),
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
