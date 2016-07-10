package gosnake

// Manages the interface between python and go code.
//
// Need at least one global to route incoming python calls into their respective handlers
// Arguably the rest of this can be trashed, but I'm going to leave it here for now

//#cgo pkg-config: python-2.7
//#include "binding.h"
import "C"

import (
	"github.com/liamzdenek/go-pthreads"
	"github.com/sbinet/go-python"
)

var (
	opChan chan *Operation // Channel by which threads receive operations

	POOL_SIZE int = 10 // Thread Pool
	pool      []*pthread.Thread

	_INITIALIZED bool = false // Silly initialization global. func init() doesn't work
)

// Initialize python environment and create the pthread pool
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

// Consumes an operation off the operation channel and sends it to the process function
func threadConsume() {
	// TODO: kill the thread if the channel is closed

	for {
		op := <-opChan

		if result, err := threadProcess(op); err != nil {
			op.errChan <- err
		} else {
			op.returnChan <- result
		}
	}
}

// Process an operation, aka make a call to python
func threadProcess(op *Operation) (interface{}, error) {
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
	C.check_pyerr()

	// Unpack result
	result, err := togo(ret)

	if err != nil {
		return "", err
	}

	// Cleanup
	ret.DecRef()
	args.DecRef()
	m.DecRef()
	fn.DecRef()

	python.PyGILState_Release(gil)
	return result, nil
}
