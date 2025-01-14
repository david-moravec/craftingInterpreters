package interpreter

import (
	"fmt"

	"github.com/david-moravec/golox/internal/expr"
	"github.com/david-moravec/golox/internal/scanner"
	"github.com/david-moravec/golox/internal/stmt"
)

type interpreterError struct {
	t       scanner.Token
	message string
}

func (e interpreterError) Error() string {
	return "Interpreter error"
}

type unknownTypeError struct {
}

func (e unknownTypeError) Error() string {
	return "Unknown type error"
}

type Interpreter struct {
	env Environment
}

func NewInterpreter() Interpreter {
	return Interpreter{env: NewEnvironment()}
}

func (i *Interpreter) Interpret(stmts []stmt.Stmt) []error {
	var errs []error

	for _, s := range stmts {
		err := i.execute(s)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (i *Interpreter) execute(s stmt.Stmt) error {
	return s.Accept(i)
}

func (i Interpreter) evaluate(e expr.Expr) (any, error) {
	return e.Accept(i)

}

func (i Interpreter) VisitGroupingExpr(e expr.GroupingExpr) (any, error) {
	return i.evaluate(e.Expression)
}

func (i Interpreter) VisitLiteralExpr(e expr.LiteralExpr) (any, error) {
	switch e.LitType {
	case expr.NumberType:
		return e.Number, nil
	case expr.StringType:
		return e.Str, nil
	case expr.BoolType:
		return e.Number != 0, nil
	case expr.NilType:
		return nil, nil
	}

	return nil, unknownTypeError{}
}

func (i Interpreter) VisitBinaryExpr(e expr.BinaryExpr) (any, error) {
	l, err := i.evaluate(e.Left)

	if err != nil {
		return nil, err
	}

	r, err := i.evaluate(e.Right)
	if err != nil {
		return nil, err
	}

	var result any

	switch e.Operator.Kind {
	case scanner.Minus:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) - r.(float64)
	case scanner.Plus:
		switch l.(type) {
		case float64:
			switch r.(type) {
			case float64:
				result = l.(float64) + r.(float64)
			}
		case string:
			switch r.(type) {
			case string:
				result = l.(string) + r.(string)
			}
		default:
			return nil, interpreterError{scanner.Token(e.Operator), "Operands must be numbers or strings"}
		}

	case scanner.Slash:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) / r.(float64)
	case scanner.Star:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) * r.(float64)
	case scanner.Greater:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) > r.(float64)
	case scanner.GreaterEqual:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) >= r.(float64)
	case scanner.Less:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) < r.(float64)
	case scanner.LessEqual:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) >= r.(float64)
	case scanner.BangEqual:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = !isEqual(l, r)
	case scanner.EqualEqual:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = isEqual(l, r)
	}

	return result, err
}

func (i Interpreter) VisitUnaryExpr(e expr.UnaryExpr) (any, error) {
	a, err := i.evaluate(e.Right)

	if err != nil {
		return nil, err
	}

	switch e.Operator.Kind {
	case scanner.Minus:
		if err := checkOperandNumber(scanner.Token(e.Operator), a); err != nil {
			return nil, err
		}
		return -a.(float64), nil
	case scanner.Bang:
		return !isTruthy(a), nil
	}

	return nil, nil
}

func (i Interpreter) VisitVariableExpr(e expr.VariableExpr) (any, error) {
	return i.env.get(e.Name)
}

func (i Interpreter) VisitExpressionStmt(s stmt.ExpressionStmt) error {
	_, err := i.evaluate(s.Expression)

	return err
}

func (i Interpreter) VisitPrintStmt(s stmt.PrintStmt) error {
	v, err := i.evaluate(s.Expression)

	if err != nil {
		return err
	}

	fmt.Println(stringify(v))

	return nil
}

func (i *Interpreter) VisitVarStmt(s stmt.VarStmt) error {
	val, err := i.evaluate(*s.Initializer)

	if err != nil {
		return err
	}

	i.env.define(s.Name.Lexeme, val)

	return nil
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}

	return a == b
}

func isTruthy(a any) bool {
	switch a.(type) {
	case nil:
		return false
	case bool:
		return a.(bool)
	}

	return true
}

func checkOperandNumber(o scanner.Token, r any) error {
	switch r.(type) {
	case float64:
		return nil
	}

	return interpreterError{o, "Operand must be number"}
}

func checkOperandsNumber(o scanner.Token, l any, r any) error {
	switch r.(type) {
	case float64:
		switch l.(type) {
		case float64:
			return nil
		}
	}

	return interpreterError{o, "Operands must be numbers"}
}

func stringify(a any) string {
	if a == nil {
		return "Nil"
	}
	switch a.(type) {
	case float64:
		return fmt.Sprintf("%.2f", a)
	case fmt.Stringer:
		return a.(fmt.Stringer).String()
	}

	return fmt.Sprint(a)
}
