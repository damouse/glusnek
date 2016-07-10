package gosnake

import (
	"fmt"

	"github.com/sbinet/go-python"
)

// Importing functions from python go, calling them from go, and returning the results to go

// Try to import python module name
func tryImport(name string) error {
	gil := python.PyGILState_Ensure()
	defer python.PyGILState_Release(gil)

	if m := python.PyImport_ImportModule(name); m == nil {
		return fmt.Errorf("Could not find a python module named \"%s\"", name)
	} else {
		m.DecRef()
		return nil
	}
}
