package debug

import "fmt"

var DEBUG = true

func Printf(format string, i ...interface{}) {
	if !DEBUG {
		return
	}

	fmt.Printf(format, i...)
}

func Print(i ...interface{}) {
	if !DEBUG {
		return
	}

	fmt.Println(i...)
}
