package interpreter

import (
	"fmt"
	"github.com/david-moravec/golox/internal/scanner"
)

type environmentError struct {
	name    scanner.Token
	message string
}

func (e environmentError) Error() string {
	return fmt.Sprintf("EnvError [Line %d] %s '%s'", e.name.Line, e.message, e.name.Lexeme)
}

func undefinedError(name scanner.Token) environmentError {
	return environmentError{message: "Undefined variable", name: name}
}

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
		return nil, undefinedError(name)
	}
	return val, nil
}

func (e Environment) getAt(dist int, name scanner.Token) (any, error) {
	return e.ancestor(dist).get(name)
}

func (e Environment) ancestor(dist int) *Environment {
	env := &e
	for range dist {
		env = env.enclosing
	}
	return env
}

func (e *Environment) assignAt(dist int, name scanner.Token, value any) error {
	return e.ancestor(dist).assign(name, value)
}

func (e *Environment) assign(name scanner.Token, value any) error {
	_, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.assign(name, value)
		}
		return undefinedError(name)
	}
	e.values[name.Lexeme] = value
	return nil
}
