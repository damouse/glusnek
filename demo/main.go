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

	r, e := pymodule.Call("adder", "callback", "callme", 1, "2", 3.3)
	fmt.Println("Result: ", r, e)
}

func main() {
	demos()
}
