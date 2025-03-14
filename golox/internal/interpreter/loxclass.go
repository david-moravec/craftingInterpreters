package interpreter

import (
	"fmt"

	"github.com/david-moravec/golox/internal/scanner"
)

type LoxClass struct {
	Name       scanner.Token
	Methods    map[string]LoxFunction
	Superclass *LoxClass
}

func (c LoxClass) String() string {
	return c.Name.Lexeme
}

func (c LoxClass) Arity() int {
	meth := c.findMethod("init")
	if meth != nil {
		return meth.Arity()
	}
	return 0
}

func (c LoxClass) Call(i Interpreter, args []any) (any, error) {
	meth := c.findMethod("init")
	instance := &LoxInstance{class: c, fields: make(map[string]any)}
	if meth != nil {
		meth.bind(instance).Call(i, args)
	}
	return instance, nil
}

func (i LoxClass) findMethod(name string) *LoxFunction {
	meth, ok := i.Methods[name]
	if ok {
		return &meth
	} else {
		if i.Superclass != nil {
			return i.Superclass.findMethod(name)
		}

		return nil
	}
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
	meth := i.class.findMethod(name.Lexeme)
	if meth != nil {
		return meth.bind(i), nil
	}
	return nil, runtimeError{t: name, message: fmt.Sprintf("Undefined property %s.", name.Lexeme)}
}

func (i *LoxInstance) set(name scanner.Token, value any) (any, error) {
	i.fields[name.Lexeme] = value

	return value, nil
}

func (i LoxInstance) String() string {
	return fmt.Sprintf("%s instance", i.class.String())
}
