package golox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/david-moravec/golox/internal/interpreter"
	"github.com/david-moravec/golox/internal/parser"
	"github.com/david-moravec/golox/internal/scanner"
)

func run(source string) error {
	scanner := scanner.NewScanner(source)
	tokens, err := scanner.ScanTokens()

	if err != nil {
		return err
	}

	p := parser.NewParser(tokens)
	e, err := p.Parse()

	if err != nil {
		return err
	}

	err = interpreter.NewInterpreter().Interpret(e)

	if err != nil {
		return err
	}

	return nil
}

type GoLox struct {
	hadError bool
}

func (g *GoLox) RunFile(filename string) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		g.HandleError(err)
	}

	err = run(string(bytes))

	if err != nil {
		g.HandleError(err)
	}

	if g.hadError {
		os.Exit(64)
	}
}

func (g *GoLox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			g.HandleError(err)
		}
		if line == "" {
			break
		}
		run(line)
		g.hadError = false
	}
}

func (g *GoLox) HandleError(err error) {
	fmt.Printf(err.Error())
	g.hadError = true
}

var GoLoxGlobal = GoLox{false}
