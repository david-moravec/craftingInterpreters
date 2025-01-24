package stmt

import (
	"errors"

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

type BlockStmt struct {
	Statements []Stmt
}

func (s BlockStmt) Accept(v StmtVisitor) error {
	return v.VisitBlockStmt(s)
}

type IfStmt struct {
	Condition  expr.Expr
	ThenBranch Stmt
	ElseBranch *Stmt
}

func (s IfStmt) Accept(v StmtVisitor) error {
	return v.VisitIfStmt(s)
}

type WhileStmt struct {
	Condition expr.Expr
	Body      Stmt
}

func (s WhileStmt) Accept(v StmtVisitor) error {
	return v.VisitWhileStmt(s)
}

type FunctionStmt struct {
	Name   scanner.Token
	Params []scanner.Token
	Body   BlockStmt
}

func (s FunctionStmt) Accept(v StmtVisitor) error {
	return v.VisitFunctionStmt(s)
}

type ReturnStmt struct {
	Keyword scanner.Token
	Value   expr.Expr
}

func (s ReturnStmt) Accept(v StmtVisitor) error {
	return v.VisitReturnStmt(s)
}

type Return struct {
	Value any
}

func (r Return) Error() string {
	return "Return"
}

type StmtVisitor interface {
	VisitPrintStmt(PrintStmt) error
	VisitExpressionStmt(ExpressionStmt) error
	VisitVarStmt(VarStmt) error
	VisitBlockStmt(BlockStmt) error
	VisitIfStmt(IfStmt) error
	VisitWhileStmt(WhileStmt) error
	VisitFunctionStmt(FunctionStmt) error
	VisitReturnStmt(ReturnStmt) error
}

func DefaultVisitBlockStmt(s BlockStmt, v StmtVisitor) error {
	var errs []error
	for _, st := range s.Statements {
		err := st.Accept(v)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
