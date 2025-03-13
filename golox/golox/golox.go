package golox

import (
	"bufio"
	"fmt"
	"os"

	"github.com/david-moravec/golox/internal/interpreter"
	"github.com/david-moravec/golox/internal/parser"
	"github.com/david-moravec/golox/internal/resolver"
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
	stmts, errs := p.Parse()

	if len(errs) != 0 {
		fmt.Println(errs)

		return nil
	}

	res := resolver.NewResolver(interpreter)
	err := res.Resolve(stmts)
	if err != nil {
		fmt.Println(err)

		return nil
	}

	err = interpreter.Interpret(stmts)

	if err != nil {
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
