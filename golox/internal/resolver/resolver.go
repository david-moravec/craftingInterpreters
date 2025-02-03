package resolver

import (
	"errors"
	"github.com/david-moravec/golox/internal/expr"
	"github.com/david-moravec/golox/internal/interpreter"
	"github.com/david-moravec/golox/internal/scanner"
	"github.com/david-moravec/golox/internal/stmt"
)

type scope map[string]bool

type scopeStack []scope

func (s scopeStack) push(sc scope) scopeStack {
	return append(s, sc)
}

func (s scopeStack) pop() (scopeStack, scope, error) {
	l := len(s)
	if s.isEmpty() {
		return s, nil, errors.New("Empty scope stack")
	}

	return s[:l-1], s[l-1], nil
}

func (s scopeStack) peek() (scope, bool) {
	if s.isEmpty() {
		return nil, false
	}

	return s[len(s)-1], true
}

func (s scopeStack) isEmpty() bool {
	return len(s) == 0
}

type resolverError struct {
	t       scanner.Token
	message string
}

func (e resolverError) Error() string {
	return "Interpreter error"
}

type Resolver struct {
	interpreter interpreter.Interpreter
	scopes      scopeStack
}

func NewResolver(i interpreter.Interpreter) Resolver {
	return Resolver{interpreter: i}
}

func (r *Resolver) VisitPrintStmt(s stmt.PrintStmt) error {
	_, err := r.resolveExpr(s.Expression)
	return err
}

func (r *Resolver) VisitExpressionStmt(s stmt.ExpressionStmt) error {
	_, err := r.resolveExpr(s.Expression)
	return err
}

func (r *Resolver) VisitVarStmt(s stmt.VarStmt) error {
	err := r.declare(*s.Name)
	if err != nil {
		return err
	}
	if s.Initializer != nil {
		_, err := r.resolveExpr(*s.Initializer)
		if err != nil {
			return err
		}
	}
	return r.define(*s.Name)
}

func (r *Resolver) VisitBlockStmt(s stmt.BlockStmt) error {
	r.beginScope()
	err := r.resolveStmtMany(s.Statements)
	if err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *Resolver) VisitIfStmt(s stmt.IfStmt) error {
	_, err := r.resolveExpr(s.Condition)
	if err != nil {
		return err
	}
	err = r.resolveStmt(s.ThenBranch)
	if err != nil {
		return err
	}
	if s.ElseBranch != nil {
		err = r.resolveStmt(*s.ElseBranch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) VisitWhileStmt(s stmt.WhileStmt) error {
	var err error = nil
	if s.Condition != nil {
		_, err = r.resolveExpr(s.Condition)
		if err != nil {
			return err
		}
	}
	return r.resolveStmt(s.Body)
}

func (r *Resolver) VisitFunctionStmt(s stmt.FunctionStmt) error {
	err := r.declare(s.Name)
	if err != nil {
		return err
	}
	err = r.define(s.Name)
	if err != nil {
		return err
	}
	return r.resolveFunction(s)
}

func (r *Resolver) VisitReturnStmt(s stmt.ReturnStmt) error {
	var err error = nil
	if s.Value != nil {
		_, err = r.resolveExpr(s.Value)
	}
	return err
}

func (r *Resolver) VisitUnaryExpr(e expr.UnaryExpr) (any, error) {
	return r.resolveExpr(e.Right)
}

func (r *Resolver) VisitBinaryExpr(e expr.BinaryExpr) (any, error) {
	_, err := r.resolveExpr(e.Right)
	if err != nil {
		return nil, err
	}
	return r.resolveExpr(e.Left)
}

func (r *Resolver) VisitLiteralExpr(e expr.LiteralExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(e expr.GroupingExpr) (any, error) {
	return r.resolveExpr(e.Expression)
}

func (r *Resolver) VisitVariableExpr(e expr.VariableExpr) (any, error) {
	if !r.scopes.isEmpty() {
		sc, ok := r.scopes.peek()
		if ok {
			is_init, ok := sc[e.Name.Lexeme]
			if ok && !is_init {
				return nil, errors.New("Can't read variable in its initializer")

			}

		}
	}

	return nil, r.resolveLocal(e, e.Name)
}

func (r *Resolver) VisitAssignExpr(e expr.AssignExpr) (any, error) {
	_, err := r.resolveExpr(e.Value)
	if err != nil {
		return nil, err
	}
	return nil, r.resolveLocal(e, e.Name)
}

func (r *Resolver) VisitLogicalExpr(e expr.LogicalExpr) (any, error) {
	_, err := r.resolveExpr(e.Right)
	if err != nil {
		return nil, err
	}
	return r.resolveExpr(e.Left)
}

func (r *Resolver) VisitCallExpr(e expr.CallExpr) (any, error) {
	_, err := r.resolveExpr(e.Callee)
	if err != nil {
		return nil, err
	}
	for _, ex := range e.Arguments {
		_, err = r.resolveExpr(ex)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (r *Resolver) resolveFunction(f stmt.FunctionStmt) error {
	r.beginScope()
	for _, param := range f.Params {
		err := r.declare(param)
		if err != nil {
			return err
		}
		err = r.define(param)
		if err != nil {
			return err
		}
	}
	err := r.resolveStmt(f.Body)
	if err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *Resolver) resolveLocal(e expr.Expr, name scanner.Token) error {
	for i := len(r.scopes) - 1; i >= 0; i = i - 1 {
		_, ok := r.scopes[i][name.Lexeme]
		if ok {
			r.interpreter.resolve(name, len(r.scopes)-1-i)
		}
	}

	return nil
}

func (r *Resolver) resolveExpr(e expr.Expr) (any, error) {
	return e.Accept(r)
}

func (r *Resolver) resolveStmt(s stmt.Stmt) error {
	return s.Accept(r)
}

func (r *Resolver) resolveStmtMany(stmts []stmt.Stmt) error {
	for _, s := range stmts {
		err := r.resolveStmt(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) declare(n scanner.Token) error {
	// do not resolve global
	if r.scopes.isEmpty() {
		return nil
	}
	s, _ := r.scopes.peek()
	s[n.Lexeme] = false
	return nil
}

func (r *Resolver) define(n scanner.Token) error {
	// do not resolve global
	if r.scopes.isEmpty() {
		return nil
	}
	s, _ := r.scopes.peek()
	s[n.Lexeme] = true
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes = r.scopes.push(scope{})
}

func (r *Resolver) endScope() error {
	scopes, _, err := r.scopes.pop()
	if err != nil {
		return err
	}
	r.scopes = scopes
	return nil
}
