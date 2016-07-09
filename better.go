package gosnake

/*
#cgo pkg-config: python-2.7
#include "Python.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/liamzdenek/go-pthreads"
	"github.com/sbinet/go-python"
)

// We can likely pool the pthreads, but that sounds like a problem for future mickey
func PyCall() string {
	done := make(chan string)

	t := pthread.Create(func() {

		// Lets try it with the gopython code...
		gil := C.PyGILState_Ensure()
		defer C.PyGILState_Release(gil)

		m := python.PyImport_ImportModule("adder")
		fn := m.GetAttrString("run")

		// Pack the arguments
		args := python.PyTuple_New(2)
		python.PyTuple_SET_ITEM(args, 0, python.PyString_FromString("Hello!"))
		python.PyTuple_SET_ITEM(args, 1, python.PyInt_FromLong(1234))

		ret := fn.CallObject(args)
		done <- python.PyString_AsString(ret)

		// cb(ret)
		// done <- nil
	})

	defer t.Kill()
	return <-done
}

func BetterTest() {
	n := 80
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(gid int) {
			for j := 0; j < 2; j++ {
				rs := PyCall()

				var retjson []interface{}
				if _err := json.Unmarshal([]byte(rs), &retjson); _err != nil {
					panic(fmt.Errorf("got invalid result from python function, `%v`\n", rs))
				}

				fmt.Printf("gorutine: %d iteration: %d pid: %f\n", gid, j, retjson[0])
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("\nInternal Done")
}
