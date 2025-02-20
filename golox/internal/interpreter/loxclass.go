package interpreter

import (
	"fmt"

	"github.com/david-moravec/golox/internal/scanner"
)

type LoxClass struct {
	Name scanner.Token
}

func (c LoxClass) String() string {
	return c.Name.Lexeme
}

func (c LoxClass) Arity() int {
	return 0
}

func (c LoxClass) Call(i Interpreter, args []any) (any, error) {
	return LoxInstance{class: c}, nil
}

type LoxInstance struct {
	class LoxClass
}

func (i LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.class.String())
}
