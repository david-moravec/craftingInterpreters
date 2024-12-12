package expr

import (
	"github.com/david-moravec/golox/internal/scanner"
)

type Expr interface {
	Accept(ExprVisitor) any
}

type Operator scanner.Token

func (o Operator) String() string {
	return scanner.Token(o).String()
}

type UnaryExpr struct {
	Operator Operator
	Right    Expr
}

func NewUnary(o Operator, r Expr) *UnaryExpr {
	return &UnaryExpr{Operator: o, Right: r}
}

func (e UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(e)
}

type BinaryExpr struct {
	Left     Expr
	Operator Operator
	Right    Expr
}

func NewBinary(left Expr, right Expr, operator Operator) *BinaryExpr {
	return &BinaryExpr{Left: left, Operator: operator, Right: right}
}

func (e BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(e)
}

type LiteralType int

const (
	NumberType LiteralType = iota
	StringType
	BoolType
	NilType
)

type LiteralExpr struct {
	LitType LiteralType
	Number  float64
	Str     string
}

func NewLiteral(t LiteralType, n float64, str string) *LiteralExpr {
	return &LiteralExpr{LitType: t, Number: n, Str: str}
}

func (e LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(e)
}

type GroupingExpr struct {
	Expression Expr
}

func NewGroup(e Expr) *GroupingExpr {
	return &GroupingExpr{e}
}

func (e GroupingExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(e)
}

type ExprVisitor interface {
	VisitUnaryExpr(expr UnaryExpr) any
	VisitBinaryExpr(expr BinaryExpr) any
	VisitLiteralExpr(expr LiteralExpr) any
	VisitGroupingExpr(expr GroupingExpr) any
}
