package interpreter

import (
	"fmt"

	"github.com/david-moravec/golox/internal/stmt"
)

type LoxCallable interface {
	Call(i Interpreter, args []any) (any, error)
	Arity() int
	String() string
}

type LoxFunction struct {
	declaration stmt.FunctionStmt
}

func (l LoxFunction) Call(i Interpreter, args []any) (any, error) {
	env := NewEnvironment(&i.globals)

	for i := range l.declaration.Params {
		env.define(l.declaration.Params[i].Lexeme, args[i])
	}

	err := i.executeBlock(l.declaration.Body, env)
	return nil, err
}

func (l LoxFunction) Arity() int {
	return len(l.declaration.Params)
}

func (l LoxFunction) String() string {
	return fmt.Sprintf("fn< %s >", l.declaration.Name.Lexeme)
}
