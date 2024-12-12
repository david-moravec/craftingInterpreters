package interpreter

import (
	"github.com/david-moravec/golox/internal/expr"
	"github.com/david-moravec/golox/internal/scanner"
)

type Interpreter struct {
}

func (i Interpreter) evaluate(e expr.Expr) any {
	return e.Accept(i)

}

func (i Interpreter) VisitGroupingExpr(e expr.GroupingExpr) any {
	return i.evaluate(e.Expression)
}

func (i Interpreter) VisitLiteralExpr(e expr.LiteralExpr) any {
	return ""
}

func (i Interpreter) VisitBinaryExpr(e expr.BinaryExpr) any {
	l := i.evaluate(e.Left)
	r := i.evaluate(e.Right)

	switch e.Operator.Kind {
	case scanner.Minus:
		return l.(float64) - r.(float64)
	case scanner.Plus:
		switch l.(type) {
		case float64:
			return l.(float64) + r.(float64)
		case string:
			return l.(string) + r.(string)

		}

	case scanner.Slash:
		return l.(float64) / r.(float64)
	case scanner.Star:
		return l.(float64) * r.(float64)
	case scanner.Greater:
		return l.(float64) > r.(float64)
	case scanner.GreaterEqual:
		return l.(float64) >= r.(float64)
	case scanner.Less:
		return l.(float64) < r.(float64)
	case scanner.LessEqual:
		return l.(float64) >= r.(float64)
	case scanner.BangEqual:
		return !isEqual(l, r)
	case scanner.EqualEqual:
		return isEqual(l, r)
	}

	return nil
}

func (i Interpreter) VisitUnaryExpr(e expr.UnaryExpr) any {
	a := i.evaluate(e.Right)

	switch e.Operator.Kind {
	case scanner.Minus:
		return -a.(float64)
	case scanner.Bang:
		return !isTruthy(a)
	}

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
