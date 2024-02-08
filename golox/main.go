package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func logErr(err error) {
	fmt.Print(err)
}

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Printf("Usage: golox [script]\n")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}

}

func runFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		logErr(err)
	}

	run(string(bytes))
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			logErr(err)
		}
		if line == "" {
			break
		}
		run(line)
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
