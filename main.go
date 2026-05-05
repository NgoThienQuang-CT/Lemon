package main

import (
	"fmt"
	"os"

	"lemon/interpreter"
)

func main() {
	if len(os.Args) == 2 {
		interpreter.RunFile(os.Args[1])
	} else if len(os.Args) == 1 {
		interpreter.StartRepl()
	} else {
		fmt.Println("Usage: lemon <path>")
		os.Exit(64)
	}
}
