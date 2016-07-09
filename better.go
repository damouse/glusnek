package gosnake

/*
#cgo pkg-config: python-2.7
#include "Python.h"
*/
import "C"

import (
	"fmt"
	"sync"

	"github.com/liamzdenek/go-pthreads"
	"github.com/sbinet/go-python"
)

// The pthreads live in a pool. They should be sent messages through the opChan channel
var POOL_SIZE int = 5
var pool []*pthread.Thread
var opChan chan *Operation

type Operation struct {
	args       string
	returnChan chan string
}

func Initialize() {
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

		gil := python.PyGILState_Ensure()

		m := python.PyImport_ImportModule("adder")
		fn := m.GetAttrString("run")

		m.IncRef()
		fn.IncRef()

		// Pack the arguments
		// args := python.PyTuple_New(2)
		// python.PyTuple_SET_ITEM(args, 0, python.PyString_FromString("Hello!"))
		// python.PyTuple_SET_ITEM(args, 1, python.PyInt_FromLong(1234))

		args := python.PyTuple_New(0)
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

		// Release the Gil
		python.PyGILState_Release(gil)

		op.returnChan <- resultString
	}
}

// We can likely pool the pthreads, but that sounds like a problem for future mickey
func PyCall() string {
	op := &Operation{
		args:       "hello",
		returnChan: make(chan string),
	}

	opChan <- op

	ret := <-op.returnChan
	return ret
}

func BetterTest() {
	n := 1000
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(gid int) {
			for j := 0; j < 1000; j++ {

				ret := PyCall()

				fmt.Printf("(%d  %d) \t%s\n", gid, j, ret)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("\nInternal Done")
}
