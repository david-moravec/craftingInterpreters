package expr

import (
	"bytes"
	"strconv"
)

type AstPrinter struct {
}

func (p AstPrinter) print(e Expr) string {
	return e.accept(p).(string)
}

func (p AstPrinter) visitUnaryExpr(expr UnaryExpr) any {
	return p.parenthesize(expr.operator.String(), []Expr{expr.right})

}
func (p AstPrinter) visitBinaryExpr(expr BinaryExpr) any {
	return p.parenthesize(expr.operator.String(), []Expr{expr.left, expr.right})

}
func (p AstPrinter) visitLiteralExpr(expr LiteralExpr) any {
	switch expr.litType {
	case NumberType:
		return strconv.Itoa(expr.number)
	case StringType:
		return expr.str
	case BoolType:
		return strconv.FormatBool(expr.number != 0)
	case NilType:
		return "nil"
	}

	return ""
}
func (p AstPrinter) visitGroupingExpr(expr GroupingExpr) any {
	return p.parenthesize("group", []Expr{expr.expression})
}

func (p *AstPrinter) parenthesize(name string, expressions []Expr) string {
	s := new(bytes.Buffer)

	s.WriteByte('(')
	s.WriteString(name)

	for _, e := range expressions {
		s.WriteString(" " + p.print(e))

	}
	s.WriteByte(')')

	return s.String()
}
