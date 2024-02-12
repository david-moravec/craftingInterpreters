package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type GoLox struct {
	hadError bool
}

func (g *GoLox) runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		error(0, err.Error())
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
			error(0, err.Error())
		}
		if line == "" {
			break
		}
		run(line)
		g.hadError = false
	}
}

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

type Scanner struct {
	Source string
}

func (s Scanner) scanTokens() []string {
	return strings.Split(s.Source, " ")
}

func run(source string) {
	scanner := Scanner{Source: source}
	tokens := scanner.scanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s", line, where, message)
	hadError := true
}
