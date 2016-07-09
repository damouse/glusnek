package gosnake

/*
#cgo pkg-config: python-2.7
#include "Python.h"
#include <stdlib.h>
#include <string.h>
#include <pthread.h>
#include <signal.h>
#include <unistd.h>
#include <stdio.h>

static void caught_sig(int sig);

static void caught_sig(int sig) {
    PyErr_Print();
    fprintf(stdout, "gosnake: SIGSEGV on %d\n", sig);

    signal(SIGSEGV, caught_sig);
    pthread_exit(NULL);
}

static void register_sig_handler() {
    signal(SIGSEGV, caught_sig);
}
*/
import "C"

import (
	"fmt"
	"sync"

	"github.com/liamzdenek/go-pthreads"
	"github.com/sbinet/go-python"
)

var module *python.PyObject
var fn *python.PyObject

func Initialize() {
	gil := python.PyGILState_Ensure()
	module = python.PyImport_ImportModule("adder")
	fn = module.GetAttrString("run")
	python.PyGILState_Release(gil)
	checkPyErr()
}

// We can likely pool the pthreads, but that sounds like a problem for future mickey
func PyCall() string {
	done := make(chan string)

	t := pthread.Create(func() {
		gil := python.PyGILState_Ensure()

		// m := python.PyImport_ImportModule("adder")
		// m.IncRef()
		// fn := module.GetAttrString("run")
		// fn.IncRef()

		// Pack the arguments
		// args := python.PyTuple_New(2)
		// python.PyTuple_SET_ITEM(args, 0, python.PyString_FromString("Hello!"))
		// python.PyTuple_SET_ITEM(args, 1, python.PyInt_FromLong(1234))

		args := python.PyTuple_New(0)
		args.IncRef()

		ret := fn.CallObject(args)

		defer ret.DecRef()
		defer args.DecRef()

		checkPyErr()
		python.PyGILState_Release(gil)
		checkPyErr()

		done <- "Done"

		// done <- python.PyString_AsString(ret)
	})

	defer t.Kill()
	defer close(done)
	return <-done
}

func checkPyErr() {
	if python.PyErr_CheckSignals() {
		python.PyErr_PrintEx(false)
	}
}

func BetterTest() {
	n := 10
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(gid int) {
			for j := 0; j < 1000; j++ {
				PyCall()
				fmt.Printf("gorutine: %d iteration: %d\n", gid, j)

				// var retjson []interface{}
				// if _err := json.Unmarshal([]byte(rs), &retjson); _err != nil {
				// 	panic(fmt.Errorf("got invalid result from python function, `%v`\n", rs))
				// }

				// fmt.Printf("gorutine: %d iteration: %d pid: %f\n", gid, j, retjson[0])
			}

			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("\nInternal Done")
}
