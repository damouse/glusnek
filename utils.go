package gosnake

// Entry point into core and rest of paradrop system

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

// Dummy entry point tester
func Create_thread(num int) {
	RunTest(num)
}
