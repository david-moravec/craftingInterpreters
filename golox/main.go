package main

import (
	"fmt"
	"github.com/david-moravec/golox/golox"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Printf("Usage: golox [script]\n")
		os.Exit(64)
	} else if len(args) == 1 {
		golox.GoLoxGlobal.RunFile(args[0])
	} else {
		golox.GoLoxGlobal.RunPrompt()
	}

}
