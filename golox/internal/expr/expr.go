package expr

type Expr interface {
	accept(ExprVisitor) any
}

type Operator byte

type UnaryExpr struct {
	operator Operator
	right    Expr
}

func (e UnaryExpr) accept(visitor ExprVisitor) any {
	return visitor.visitUnaryExpr(e)
}

type BinaryExpr struct {
	left     Expr
	right    Expr
	operator Operator
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
	lit_type LiteralType
	number   int
	str      string
}

func (e LiteralExpr) accept(visitor ExprVisitor) any {
	return visitor.visitLiteralExpr(e)
}

type GroupingExpr struct {
	expression Expr
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
