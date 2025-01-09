package stmt

import (
	"github.com/david-moravec/golox/internal/expr"
)

type Stmt interface {
	Accept(StmtVisitor) error
}

type ExpressionStmt struct {
	Expression expr.Expr
}

func (s ExpressionStmt) Accept(v StmtVisitor) error {
	return v.VisitExpressionStmt(s)
}

type PrintStmt struct {
	Expression expr.Expr
}

func (s PrintStmt) Accept(v StmtVisitor) error {
	return v.VisitPrintStmt(s)
}

type StmtVisitor interface {
	VisitPrintStmt(PrintStmt) error
	VisitExpressionStmt(ExpressionStmt) error
}
