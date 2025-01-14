package interpreter

import (
	"errors"
	"fmt"
	"github.com/david-moravec/golox/internal/scanner"
)

type Environment struct {
	values map[string]any
}

func NewEnvironment() Environment {
	return Environment{values: map[string]any{}}
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e Environment) get(name scanner.Token) (any, error) {
	val, ok := e.values[name.Lexeme]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Undefined variable '%s'", name.Lexeme))
	}
	return val, nil
}

func (e *Environment) assign(name scanner.Token, value any) error {
	_, ok := e.values[name.Lexeme]
	if !ok {
		return errors.New(fmt.Sprintf("Undefined variable '%s'", name.Lexeme))
	}
	e.values[name.Lexeme] = value
	return nil
}
