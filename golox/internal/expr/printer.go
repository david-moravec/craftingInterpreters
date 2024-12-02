package expr

import (
	"bytes"
	"strconv"
)

type AstPrinter struct {
	expr_formated string
	ast_formated  string
}

func (p *AstPrinter) Print_ast(expr Expr) string {
	p.ast_formated = ""

	expr.accept(p)

	return p.ast_formated
}

func (p *AstPrinter) format_expr(expr Expr) string {
	p.expr_formated = ""

	expr.accept(p)

	return p.expr_formated
}

func (p *AstPrinter) visitUnaryExpr(expr *UnaryExpr) {
	p.parenthesize(string(expr.operator), []Expr{expr.right})

}
func (p *AstPrinter) visitBinaryExpr(expr *BinaryExpr) {
	p.parenthesize(string(expr.operator), []Expr{expr.left, expr.right})

}
func (p *AstPrinter) visitLiteralExpr(expr *LiteralExpr) {
	switch expr.lit_type {
	case NumberType:
		p.expr_formated = strconv.Itoa(expr.number)
	case StringType:
		p.expr_formated = expr.str
	case BoolType:
		p.expr_formated = strconv.FormatBool(expr.number != 0)
	case NilType:
		p.expr_formated = "nil"
	}

}
func (p *AstPrinter) visitGroupingExpr(expr *GroupingExpr) {
	p.parenthesize("group", []Expr{expr.expression})
}

func (p *AstPrinter) parenthesize(name string, expressions []Expr) {
	s := new(bytes.Buffer)

	s.WriteByte('(')
	s.WriteString(name)

	for _, e := range expressions {
		s.WriteString(" " + p.format_expr(e))

	}
	s.WriteByte(')')

	p.ast_formated += s.String()
}
