package stmt

import (
	"github.com/david-moravec/golox/internal/expr"
	"github.com/david-moravec/golox/internal/scanner"
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

type VarStmt struct {
	Name        *scanner.Token
	Initializer *expr.Expr
}

func (s VarStmt) Accept(v StmtVisitor) error {
	return v.VisitVarStmt(s)
}

type StmtVisitor interface {
	VisitPrintStmt(PrintStmt) error
	VisitExpressionStmt(ExpressionStmt) error
	VisitVarStmt(VarStmt) error
}
