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
	return e.Accept(p).(string)
}

func (p AstPrinter) VisitUnaryExpr(expr UnaryExpr) any {
	return p.parenthesize(expr.Operator.String(), []Expr{expr.Right})

}
func (p AstPrinter) VisitBinaryExpr(expr BinaryExpr) any {
	return p.parenthesize(expr.Operator.String(), []Expr{expr.Left, expr.Right})

}
func (p AstPrinter) VisitLiteralExpr(expr LiteralExpr) any {
	switch expr.LitType {
	case NumberType:
		return strconv.FormatFloat(expr.Number, 'f', -1, 64)
	case StringType:
		return expr.Str
	case BoolType:
		return strconv.FormatBool(expr.Number != 0)
	case NilType:
		return "nil"
	}

	return ""
}
func (p AstPrinter) VisitGroupingExpr(expr GroupingExpr) any {
	return p.parenthesize("group", []Expr{expr.Expression})
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
