package main

import "github.com/damouse/gosnake"

// Simple starter script to kick off the go core

func main() {
	// Init python environment
	gosnake.InitPyEnv()

	// Start the server
	end := make(chan bool)

	// os.Args

	for i := 0; i < 1; i++ {
		go gosnake.Create_thread(i)
	}

	<-end

	// Run forever
}
