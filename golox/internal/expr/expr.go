package expr

import (
	"github.com/david-moravec/golox/internal/scanner"
)

type Expr interface {
	accept(ExprVisitor) any
}

type Operator scanner.Token

func (o Operator) String() string {
	return scanner.Token(o).String()
}

type UnaryExpr struct {
	operator Operator
	right    Expr
}

func NewUnary(o Operator, r Expr) *UnaryExpr {
	return &UnaryExpr{operator: o, right: r}
}

func (e UnaryExpr) accept(visitor ExprVisitor) any {
	return visitor.visitUnaryExpr(e)
}

type BinaryExpr struct {
	left     Expr
	operator Operator
	right    Expr
}

func NewBinary(left Expr, right Expr, operator Operator) *BinaryExpr {
	return &BinaryExpr{left: left, operator: operator, right: right}
}

func (e BinaryExpr) accept(visitor ExprVisitor) any {
	return visitor.visitBinaryExpr(e)
}

type LiteralType int

const (
	NumberType LiteralType = iota
	StringType
	BoolType
	NilType
)

type LiteralExpr struct {
	litType LiteralType
	number  float64
	str     string
}

func NewLiteral(t LiteralType, n float64, str string) *LiteralExpr {
	return &LiteralExpr{litType: t, number: n, str: str}
}

func (e LiteralExpr) accept(visitor ExprVisitor) any {
	return visitor.visitLiteralExpr(e)
}

type GroupingExpr struct {
	expression Expr
}

func NewGroup(e Expr) *GroupingExpr {
	return &GroupingExpr{e}
}

func (e GroupingExpr) accept(visitor ExprVisitor) any {
	return visitor.visitGroupingExpr(e)
}

type ExprVisitor interface {
	visitUnaryExpr(expr UnaryExpr) any
	visitBinaryExpr(expr BinaryExpr) any
	visitLiteralExpr(expr LiteralExpr) any
	visitGroupingExpr(expr GroupingExpr) any
}
