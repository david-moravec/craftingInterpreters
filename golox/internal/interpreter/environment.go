package interpreter

import (
	"errors"
	"fmt"
	"github.com/david-moravec/golox/internal/scanner"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewGlobalEnvironment() Environment {
	return Environment{values: map[string]any{}, enclosing: nil}
}

func NewEnvironment(enclosing *Environment) Environment {
	return Environment{values: map[string]any{}, enclosing: enclosing}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e Environment) get(name scanner.Token) (any, error) {
	val, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}
		return nil, errors.New(fmt.Sprintf("[Line %d] Undefined variable '%s'", name.Line, name.Lexeme))
	}
	return val, nil
}

func (e *Environment) assign(name scanner.Token, value any) error {
	_, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.assign(name, value)
		}
		return errors.New(fmt.Sprintf("[Line %d] Undefined variable '%s'", name.Line, name.Lexeme))
	}
	e.values[name.Lexeme] = value
	return nil
}
