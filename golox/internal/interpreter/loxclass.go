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
	return &LoxInstance{class: c, fields: make(map[string]any)}, nil
}

type LoxInstance struct {
	class  LoxClass
	fields map[string]any
}

func (i *LoxInstance) get(name scanner.Token) (any, error) {
	val, ok := i.fields[name.Lexeme]

	if !ok {
		return nil, runtimeError{t: name, message: fmt.Sprintf("Undefined property %s.", name.Lexeme)}
	}

	return val, nil
}

func (i *LoxInstance) set(name scanner.Token, value any) (any, error) {
	i.fields[name.Lexeme] = value

	return nil, nil
}

func (i LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.class.String())
}
