package gosnake

//#cgo pkg-config: python-2.7
//#include "binding.h"
import "C"

import (
	"encoding/json"

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
	args   string

	returnChan chan string
}

// Interestingly this doesnt work as a package init function. Not sure why
func initializeBinding() {
	if _INITIALIZED {
		return
	}

	C.pyinit(C.int(1))
	// C.register_sig_handler()

	opChan = make(chan *Operation)
	pool = make([]*pthread.Thread, POOL_SIZE)

	for i := 0; i < POOL_SIZE; i++ {
		t := pthread.Create(threadConsume)
		pool = append(pool, &t)
	}

	_INITIALIZED = true
}

// Consume an operation from the queue
func threadConsume() {
	for {
		// Grab a new operation
		op := <-opChan

		// Lock the gil
		gil := python.PyGILState_Ensure()

		// Import the module, target function
		m := python.PyImport_ImportModule(op.module.name)
		fn := m.GetAttrString(op.target)

		m.IncRef()
		fn.IncRef()

		// Pack the arguments
		args := python.PyTuple_New(2)
		a1 := python.PyString_FromString(op.args)
		a2 := python.PyInt_FromLong(1234)
		python.PyTuple_SET_ITEM(args, 0, a1)
		python.PyTuple_SET_ITEM(args, 1, a2)

		// Retain the arguments
		args.IncRef()
		a1.IncRef()
		a2.IncRef()

		// Call into Python
		ret := fn.CallObject(args)

		// Deserialize
		resultString := python.PyString_AsString(ret)

		// Cleanup
		ret.DecRef()
		args.DecRef()
		m.DecRef()
		fn.DecRef()
		args.DecRef()
		a1.DecRef()
		a2.DecRef()

		// Release the Gil
		python.PyGILState_Release(gil)

		op.returnChan <- resultString
	}
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

// We can likely pool the pthreads, but that sounds like a problem for future mickey
// Should we check if threads are "busy"?
func Call(args string) string {
	op := &Operation{
		args:       args,
		returnChan: make(chan string),
	}

	opChan <- op

	ret := <-op.returnChan
	return ret
}
