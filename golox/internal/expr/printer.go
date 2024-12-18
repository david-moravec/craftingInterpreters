package expr

import (
	"bytes"
	"strconv"
)

type AstPrinter struct {
}

func NewPrinter() AstPrinter {
	return AstPrinter{}
}

func (p AstPrinter) Print(e Expr) string {
	s, _ := e.Accept(p)
	return s.(string)
}

func (p AstPrinter) VisitUnaryExpr(expr UnaryExpr) (any, error) {
	return p.parenthesize(expr.Operator.String(), []Expr{expr.Right}), nil

}
func (p AstPrinter) VisitBinaryExpr(expr BinaryExpr) (any, error) {
	return p.parenthesize(expr.Operator.String(), []Expr{expr.Left, expr.Right}), nil

}
func (p AstPrinter) VisitLiteralExpr(expr LiteralExpr) (any, error) {
	switch expr.LitType {
	case NumberType:
		return strconv.FormatFloat(expr.Number, 'f', -1, 64), nil
	case StringType:
		return expr.Str, nil
	case BoolType:
		return strconv.FormatBool(expr.Number != 0), nil
	case NilType:
		return "nil", nil
	}

	return "", nil
}
func (p AstPrinter) VisitGroupingExpr(expr GroupingExpr) (any, error) {
	return p.parenthesize("group", []Expr{expr.Expression}), nil
}

func (p *AstPrinter) parenthesize(name string, expressions []Expr) string {
	s := new(bytes.Buffer)

	s.WriteByte('(')
	s.WriteString(name)

	for _, e := range expressions {
		s.WriteString(" " + p.Print(e))

	}
	s.WriteByte(')')

	return s.String()
}
