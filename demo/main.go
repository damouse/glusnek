package main

//#cgo pkg-config: python-2.7
//#include "Python.h"
import "C"

import (
	"fmt"

	"github.com/damouse/gosnake"
)

// Simple starter script to kick off the go core

func sandbox() {
	fmt.Println("Hello!")

	// So C conversions to and from will work... but do we want to abandon go-python?
	a := 1
	b := C.long(a)
	c := C.PyLong_FromLong(b)
	d := C.PyLong_AsLong(c)
	e := int(d)

	fmt.Printf("%v %v %v %v %v", a, b, c, d, e)
}

func demos() {
	// Start the server
	end := make(chan bool)

	// os.Args

	for i := 0; i < 1; i++ {
		go gosnake.Create_thread(i)
	}

	<-end

	// Run forever
}

func main() {
	// sandbox()
	demos()
}
