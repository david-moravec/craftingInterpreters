package main

import (
	"bufio"
	"fmt"
	"os"
)

func run(source string) {
	scanner := Scanner{source: source}
	tokens := scanner.scanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

type GoLox struct {
	hadError bool
}

func (g *GoLox) runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		g.error(0, err.Error())
	}

	run(string(bytes))

	if g.hadError {
		os.Exit(64)
	}
}

func (g *GoLox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			g.error(0, err.Error())
		}
		if line == "" {
			break
		}
		run(line)
		g.hadError = false
	}
}

func (g *GoLox) report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s", line, where, message)
	g.hadError = true
}

func (g *GoLox) error(line int, message string) {
	g.report(line, "", message)
}
