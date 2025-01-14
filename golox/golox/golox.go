package golox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/david-moravec/golox/internal/interpreter"
	"github.com/david-moravec/golox/internal/parser"
	"github.com/david-moravec/golox/internal/scanner"
)

func run(source string, interpreter interpreter.Interpreter) error {
	scanner := scanner.NewScanner(source)
	tokens, errs := scanner.ScanTokens()

	if len(errs) != 0 {
		fmt.Println(errs)

		return nil
	}

	p := parser.NewParser(tokens)
	e, errs := p.Parse()

	if len(errs) != 0 {
		fmt.Println(errs)

		return nil
	}

	errs = interpreter.Interpret(e)

	if len(errs) != 0 {
		fmt.Println(errs)

		return nil
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
	interpreter := interpreter.NewInterpreter()

	err = run(string(bytes), interpreter)

	if err != nil {
		g.HandleError(err)
	}

	if g.hadError {
		os.Exit(64)
	}
}

func (g *GoLox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)
	interpreter := interpreter.NewInterpreter()

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			g.HandleError(err)
		}
		if line == "" {
			break
		}
		run(line, interpreter)
		g.hadError = false
	}
}

func (g *GoLox) HandleError(err error) {
	fmt.Printf(err.Error())
	g.hadError = true
}

var GoLoxGlobal = GoLox{false}
