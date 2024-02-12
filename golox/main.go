package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	golox := GoLox{hadError: false}

	if len(args) > 1 {
		fmt.Printf("Usage: golox [script]\n")
		os.Exit(64)
	} else if len(args) == 1 {
		golox.runFile(args[0])
	} else {
		golox.runPrompt()
	}

}
