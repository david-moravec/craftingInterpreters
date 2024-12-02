package expr

type Expr interface {
	accept(ExprVisitor)
}

type Operator byte

type UnaryExpr struct {
	operator Operator
	right    Expr
}

func (e *UnaryExpr) accept(visitor ExprVisitor) {
	visitor.visitUnaryExpr(e)
}

type BinaryExpr struct {
	left     Expr
	right    Expr
	operator Operator
}

func (e *BinaryExpr) accept(visitor ExprVisitor) {
	visitor.visitBinaryExpr(e)
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

func (e LiteralExpr) accept(visitor ExprVisitor) {
	visitor.visitLiteralExpr(&e)
}

type GroupingExpr struct {
	expression Expr
}

func (e *GroupingExpr) accept(visitor ExprVisitor) {
	visitor.visitGroupingExpr(e)
}

type ExprVisitor interface {
	visitUnaryExpr(expr *UnaryExpr)
	visitBinaryExpr(expr *BinaryExpr)
	visitLiteralExpr(expr *LiteralExpr)
	visitGroupingExpr(expr *GroupingExpr)
}
