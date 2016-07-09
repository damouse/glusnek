package gosnake

// Entry point into core and rest of paradrop system

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}
