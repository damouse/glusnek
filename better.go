package gosnake

/*
#cgo pkg-config: python-2.7
#include "Python.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"

	"github.com/liamzdenek/go-pthreads"
	"github.com/sbinet/go-python"
)

func create_thread(cb func(a *python.PyObject)) {
	done := make(chan error)

	var t pthread.Thread
	t = pthread.Create(func() {
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
		cb(ret)

		done <- nil
		// t.Kill()
	})

	<-done
	t.Kill()
}

func BetterTest() {
	for i := 0; i < 3; i++ {
		go func() {
			// for j := 0; j < 10; j++ {
			create_thread(func(result *python.PyObject) {
				rs := python.PyString_AsString(result)

				var retjson []interface{}
				if _err := json.Unmarshal([]byte(rs), &retjson); _err != nil {
					panic(fmt.Errorf("got invalid result from python function, `%v`\n", rs))
				}

				fmt.Printf("gorutine: %d iteration: %d pid: %f\n", i, i, retjson[0])
			})
			// }
		}()
	}
}
