package interpreter

import (
	"fmt"

	"github.com/david-moravec/golox/internal/scanner"
	"github.com/david-moravec/golox/internal/stmt"
)

type LoxCallable interface {
	Call(i Interpreter, args []any) (any, error)
	Arity() int
	String() string
}

type LoxFunction struct {
	declaration stmt.FunctionStmt
	closure     Environment
	isInit      bool
}

func (l LoxFunction) Call(i Interpreter, args []any) (any, error) {
	env := NewEnvironment(&l.closure)

	for i := range l.declaration.Params {
		env.define(l.declaration.Params[i].Lexeme, args[i])
	}

	err := i.executeStmts(l.declaration.Body, env)
	switch err.(type) {
	case stmt.Return:
		if l.isInit {
			return l.closure.getAt(0, *scanner.DummyThisToken(l.declaration.Name.Line))
		}
		return err.(stmt.Return).Value, nil
	}
	if l.isInit {
		return l.closure.getAt(0, *scanner.DummyThisToken(l.declaration.Name.Line))
	}
	return nil, err
}

func (l LoxFunction) bind(i *LoxInstance) LoxFunction {
	l.closure.define("this", i)

	return LoxFunction{declaration: l.declaration, closure: l.closure, isInit: l.isInit}
}

func (l LoxFunction) Arity() int {
	return len(l.declaration.Params)
}

func (l LoxFunction) String() string {
	return fmt.Sprintf("fn< %s >", l.declaration.Name.Lexeme)
}
