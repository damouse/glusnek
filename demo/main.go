package main

//#cgo pkg-config: python-2.7
//#include "Python.h"
import "C"

import (
	"fmt"
	"sync"

	"github.com/damouse/gosnake"
)

// Stress test the bindings
func stress() {
	module, err := gosnake.Import("adder")
	checkError(err)

	routines := 100
	iterations := 100

	wg := sync.WaitGroup{}
	wg.Add(routines)

	for i := 0; i < routines; i++ {
		go func(gid int) {
			for j := 0; j < iterations; j++ {

				ret, err := module.Call("run", "hello", 1)
				checkError(err)
				fmt.Printf("(%d  %d) \t%s\n", gid, j, ret)

			}

			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("\nInternal Done")
}

func demos() {
	module, _ := gosnake.Import("adder")
	result, err := module.Call("run", "hello", 1)

	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println("Success: ", result)
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// demos()
	stress()
}
