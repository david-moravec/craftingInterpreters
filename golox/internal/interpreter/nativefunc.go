package interpreter

import "time"

type clock struct{}

func (c clock) Arity() int {
	return 0
}

func (c clock) Call(_ Interpreter, _ []any) (any, error) {
	return time.Now().UnixMilli(), nil
}

func (c clock) String() string {
	return "<native fn>"
}
