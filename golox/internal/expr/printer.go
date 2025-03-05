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

func (p AstPrinter) VisitUnaryExpr(e UnaryExpr) (any, error) {
	return p.parenthesize(e.Operator.String(), []Expr{e.Right}), nil

}
func (p AstPrinter) VisitBinaryExpr(e BinaryExpr) (any, error) {
	return p.parenthesize(e.Operator.String(), []Expr{e.Left, e.Right}), nil

}
func (p AstPrinter) VisitLiteralExpr(e LiteralExpr) (any, error) {
	switch e.LitType {
	case NumberType:
		return strconv.FormatFloat(e.Number, 'f', -1, 64), nil
	case StringType:
		return e.Str, nil
	case BoolType:
		return strconv.FormatBool(e.Number != 0), nil
	case NilType:
		return "nil", nil
	}

	return "", nil
}
func (p AstPrinter) VisitGroupingExpr(e GroupingExpr) (any, error) {
	return p.parenthesize("group", []Expr{e.Expression}), nil
}

func (p AstPrinter) VisitLogicalExpr(e LogicalExpr) (any, error) {
	return e.Operator.String(), nil
}

func (p AstPrinter) VisitVariableExpr(e VariableExpr) (any, error) {
	return e.Name.String(), nil
}

func (p AstPrinter) VisitAssignExpr(e AssignExpr) (any, error) {
	return e.Name.String(), nil
}

func (p AstPrinter) VisitCallExpr(e CallExpr) (any, error) {
	return e.Paren.Lexeme, nil
}

func (p AstPrinter) VisitSetExpr(e SetExpr) (any, error) {
	return e.Name.Lexeme, nil
}

func (p AstPrinter) VisitGetExpr(e GetExpr) (any, error) {
	return e.Name.Lexeme, nil
}

func (p AstPrinter) VisitThisExpr(e ThisExpr) (any, error) {
	return e.Keyword.Lexeme, nil
}

func (p AstPrinter) VisitSuperExpr(e SuperExpr) (any, error) {
	return e.Keyword.Lexeme, nil
}

func (p *AstPrinter) parenthesize(name string, exprs []Expr) string {
	s := new(bytes.Buffer)

	s.WriteByte('(')
	s.WriteString(name)

	for _, e := range exprs {
		s.WriteString(" " + p.Print(e))

	}
	s.WriteByte(')')

	return s.String()
}
