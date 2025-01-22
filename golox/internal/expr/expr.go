package expr

import (
	"github.com/david-moravec/golox/internal/scanner"
)

type Expr interface {
	Accept(ExprVisitor) (any, error)
}

type Operator scanner.Token

func (o Operator) String() string {
	return scanner.Token(o).String()
}

type UnaryExpr struct {
	Operator Operator
	Right    Expr
}

func NewUnary(o Operator, r Expr) UnaryExpr {
	return UnaryExpr{Operator: o, Right: r}
}

func (e UnaryExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnaryExpr(e)
}

type BinaryExpr struct {
	Left     Expr
	Operator Operator
	Right    Expr
}

func NewBinary(left Expr, right Expr, operator Operator) BinaryExpr {
	return BinaryExpr{Left: left, Operator: operator, Right: right}
}

func (e BinaryExpr) Accept(visitor ExprVisitor) (any, error) {
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

func NewLiteral(t LiteralType, n float64, str string) LiteralExpr {
	return LiteralExpr{LitType: t, Number: n, Str: str}
}

func (e LiteralExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteralExpr(e)
}

type GroupingExpr struct {
	Expression Expr
}

func NewGroup(e Expr) GroupingExpr {
	return GroupingExpr{e}
}

func (e GroupingExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGroupingExpr(e)
}

type VariableExpr struct {
	Name scanner.Token
}

func NewVariable(name scanner.Token) VariableExpr {
	return VariableExpr{Name: name}
}

func (e VariableExpr) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariableExpr(e)
}

type AssignExpr struct {
	Name  scanner.Token
	Value Expr
}

func NewAssign(name scanner.Token, value Expr) AssignExpr {
	return AssignExpr{Name: name, Value: value}
}

func (e AssignExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitAssignExpr(e)
}

type LogicalExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (e LogicalExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitLogicalExpr(e)
}

type CallExpr struct {
	Callee    Expr
	Paren     scanner.Token
	Arguments []Expr
}

func (e CallExpr) Accept(v ExprVisitor) (any, error) {
	return v.VisitCallExpr(e)
}

type ExprVisitor interface {
	VisitUnaryExpr(e UnaryExpr) (any, error)
	VisitBinaryExpr(e BinaryExpr) (any, error)
	VisitLiteralExpr(e LiteralExpr) (any, error)
	VisitGroupingExpr(e GroupingExpr) (any, error)
	VisitVariableExpr(e VariableExpr) (any, error)
	VisitAssignExpr(e AssignExpr) (any, error)
	VisitLogicalExpr(e LogicalExpr) (any, error)
	VisitCallExpr(e CallExpr) (any, error)
}
