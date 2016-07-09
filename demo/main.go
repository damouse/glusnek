package main

//#cgo pkg-config: python-2.7
//#include "Python.h"
import "C"

import (
	"fmt"

	"github.com/damouse/gosnake"
)

func demos() {
	pymodule := gosnake.NewBinding()
	pymodule.Import("adder")

	pymodule.Export("callme", func(args []interface{}, kwargs map[string]interface{}) ([]interface{}, error) {
		fmt.Println("Go function called!", args, kwargs)
		// return args, nil
		return nil, nil
	})

	c := make(chan bool)

	for i := 0; i < 400; i++ {

		go func() {
			pymodule.Call("adder", "callback", "callme", 1, "2", 3.3)
			// fmt.Println("Result: ", r, e)
		}()
	}

	<-c
}

func main() {
	demos()
}

// fatal error: unexpected signal during runtime execution
// [signal 0x7 code=0x80 addr=0x0 pc=0x7fcdf668fcda]
