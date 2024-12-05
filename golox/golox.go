package main

import (
	"bufio"
	"fmt"
	"github.com/david-moravec/golox/internal/scanner"
	"os"
)

func run(source string) error {
	scanner := scanner.NewScanner(source)
	tokens, err := scanner.ScanTokens()

	if err != nil {
		return err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	return nil
}

type GoLox struct {
	hadError bool
}

func (g *GoLox) runFile(filename string) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		g.error(0, err)
	}

	err = run(string(bytes))

	if err != nil {
		g.error(0, err)
	}

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
			g.error(0, err)
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

func (g *GoLox) error(line int, err error) {
	g.report(line, "", err.Error())
}
