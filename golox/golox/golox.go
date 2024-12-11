package golox

import (
	"bufio"
	"fmt"
	"github.com/david-moravec/golox/internal/expr"
	"github.com/david-moravec/golox/internal/parser"
	"github.com/david-moravec/golox/internal/scanner"
	"os"
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

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }

	fmt.Println(expr.NewPrinter().Print(e))

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
