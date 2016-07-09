package gosnake

/*
ToDo to make sure this actually works
    - Exporting functions into python
    - Getting results from python function


Advanced
    - "Permanent", goroutine-safe imports
    - Performance Cleanup
        - Dynamically create threads as needed to handle requests
        - Pre-create goroutines for outbound
    - Objects for modules
*/

//#cgo pkg-config: python-2.7
//#include "better.h"
import "C"

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/liamzdenek/go-pthreads"
	"github.com/sbinet/go-python"
)

// Thread Pool
var POOL_SIZE int = 10
var pool []*pthread.Thread

// Channel by which threads receive operations
var opChan chan *Operation

type Operation struct {
	args       string
	returnChan chan string
}

func Initialize() {
	C.pyinit(C.int(1))

	opChan = make(chan *Operation)
	pool = make([]*pthread.Thread, POOL_SIZE)

	for i := 0; i < POOL_SIZE; i++ {
		t := pthread.Create(threadConsume)
		pool = append(pool, &t)
	}
}

// Consume an operation from the queue
func threadConsume() {
	for {
		// Grab a new operation
		op := <-opChan

		// Lock the gil
		gil := python.PyGILState_Ensure()

		// Import the module, target function
		m := python.PyImport_ImportModule("adder")
		fn := m.GetAttrString("run")

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

func BetterTest(routines int, iterations int) {
	wg := sync.WaitGroup{}
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func(gid int) {
			for j := 0; j < iterations; j++ {

				ret := Call("Hey boyo")

				fmt.Printf("(%d  %d) \t%s\n", gid, j, ret)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("\nInternal Done")
}
