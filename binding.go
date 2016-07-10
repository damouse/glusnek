package gosnake

// Manages the interface between python and go code.
//
// Need at least one global to route incoming python calls into their respective handlers
// Arguably the rest of this can be trashed, but I'm going to leave it here for now

//#cgo pkg-config: python-2.7
//#include "binding.h"
import "C"

import (
	"encoding/json"
	"fmt"

	"github.com/liamzdenek/go-pthreads"
	"github.com/sbinet/go-python"
)

// Thread Pool
var POOL_SIZE int = 10
var pool []*pthread.Thread

// Channel by which threads receive operations
var opChan chan *Operation

// Silly initialization global. Please fix
var _INITIALIZED bool = false

type Operation struct {
	module *Module // previously imported module
	target string  // the name of the target function
	args   []interface{}

	returnChan chan string
	errChan    chan error
}

// Interestingly this doesnt work as a package init function. Not sure why
func initializeBinding() {
	if _INITIALIZED {
		return
	} else {
		_INITIALIZED = true
	}

	C.pyinit(C.int(0))
	C.register_sig_handler()

	opChan = make(chan *Operation)
	pool = make([]*pthread.Thread, POOL_SIZE)

	for i := 0; i < POOL_SIZE; i++ {
		t := pthread.Create(threadConsume)
		pool = append(pool, &t)
	}
}

// Make sure that a module can actually be imported
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

// Thread spinloop
// TODO: kill the thread if the channel is closed
func threadConsume() {
	for {
		op := <-opChan

		if result, err := threadProcess(op); err != nil {
			op.errChan <- err
		} else {
			op.returnChan <- result
		}
	}
}

// Process a call from go to python. Should only be called from threadConsume!
func threadProcess(op *Operation) (string, error) {
	gil := python.PyGILState_Ensure()

	// Import the module, target function
	m := python.PyImport_ImportModule(op.module.name)
	fn := m.GetAttrString(op.target)
	m.IncRef()
	fn.IncRef()

	// Pack the arguments
	args, err := packTuple(op.args)

	if err != nil {
		return "", err
	}

	// Retain the arguments
	args.IncRef()

	// Call into Python
	ret := fn.CallObject(args)

	// Deserialize
	resultString := python.PyString_AsString(ret)

	// Cleanup
	ret.DecRef()
	args.DecRef()
	m.DecRef()
	fn.DecRef()

	python.PyGILState_Release(gil)
	return resultString, nil
}

// An invocation FROM python to go
// TODO: raise errors from the go invocation back to python
//export _gosnake_invoke
func _gosnake_invoke(self *C.PyObject, args *C.PyObject) *C.char {
	// Building data structures ourselves and sending them to go can be... a little icky
	// For now everything goes back as json
	a := []interface{}{"hello"}
	b, _ := json.Marshal(a)
	ret := string(b)

	cmode := C.CString(ret)
	return cmode
}

// // We can likely pool the pthreads, but that sounds like a problem for future mickey
// // Should we check if threads are "busy"?
// func Call(args string) string {
// 	op := &Operation{
// 		args:       args,
// 		returnChan: make(chan string),
// 	}

// 	opChan <- op
// 	ret := <-op.returnChan

// 	return ret
// }
