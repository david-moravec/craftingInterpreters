package interpreter

import (
	"fmt"

	"github.com/david-moravec/golox/internal/scanner"
)

type LoxClass struct {
	Name    scanner.Token
	Methods map[string]LoxFunction
}

func (c LoxClass) String() string {
	return c.Name.Lexeme
}

func (c LoxClass) Arity() int {
	meth, ok := c.Methods["init"]
	if ok {
		return meth.Arity()
	}
	return 0
}

func (c LoxClass) Call(i Interpreter, args []any) (any, error) {
	meth, ok := c.Methods["init"]
	instance := &LoxInstance{class: c, fields: make(map[string]any)}
	if ok {
		meth.bind(instance).Call(i, args)
	}
	return instance, nil
}

type LoxInstance struct {
	class  LoxClass
	fields map[string]any
}

func (i *LoxInstance) get(name scanner.Token) (any, error) {
	val, ok := i.fields[name.Lexeme]
	if ok {
		return val, nil
	}
	meth, ok := i.class.Methods[name.Lexeme]
	if ok {
		return meth.bind(i), nil
	}
	return nil, runtimeError{t: name, message: fmt.Sprintf("Undefined property %s.", name.Lexeme)}
}

func (i *LoxInstance) set(name scanner.Token, value any) (any, error) {
	i.fields[name.Lexeme] = value

	return nil, nil
}

func (i LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.class.String())
}
