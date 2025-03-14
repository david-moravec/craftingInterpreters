package resolver

import (
	"errors"
	"fmt"
	"github.com/david-moravec/golox/internal/expr"
	"github.com/david-moravec/golox/internal/interpreter"
	"github.com/david-moravec/golox/internal/scanner"
	"github.com/david-moravec/golox/internal/stmt"
)

type functionType int

const (
	NONE functionType = iota
	FUNCTION
	METHOD
	INIT
)

type classType int

const (
	NONe classType = iota
	CLASS
	SUBCLASS
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
	return fmt.Sprintf("ResolverError [Line: %d]: %s", e.t.Line, e.message)
}

type Resolver struct {
	interpreter    interpreter.Interpreter
	scopes         scopeStack
	currentFuncTy  functionType
	currentClassTy classType
}

func NewResolver(i interpreter.Interpreter) Resolver {
	return Resolver{interpreter: i, currentFuncTy: NONE, currentClassTy: NONe}
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
	err := r.Resolve(s.Statements)
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
	return r.resolveFunction(s, FUNCTION)
}

func (r *Resolver) VisitReturnStmt(s stmt.ReturnStmt) error {
	if r.currentFuncTy == NONE {
		return resolverError{t: s.Keyword, message: "Can't return from top-level code."}
	}
	var err error = nil
	if s.Value != nil {
		if r.currentFuncTy == INIT {
			return resolverError{t: s.Keyword, message: "Can't return value from inititalizer."}
		}
		_, err = r.resolveExpr(s.Value)
	}
	return err
}

func (r *Resolver) VisitClassStmt(s stmt.ClassStmt) error {
	orig := r.currentClassTy
	r.currentClassTy = CLASS
	err := r.declare(s.Name)
	if err != nil {
		return err
	}
	r.define(s.Name)
	if s.Superclass != nil {
		r.currentClassTy = SUBCLASS
		if s.Superclass.Name.Lexeme == s.Name.Lexeme {
			return resolverError{t: s.Superclass.Name, message: "Can't inherit from itself."}
		}
		r.resolveExpr(s.Superclass)
		r.beginScope()
		sc, _ := r.scopes.peek()
		sc["super"] = true
	}
	r.beginScope()
	sc, _ := r.scopes.peek()
	sc["this"] = true

	var errs []error
	for _, meth := range s.Methods {
		ty := METHOD
		if meth.Name.Lexeme == "init" {
			ty = INIT
		}
		errs = append(errs, r.resolveFunction(meth, ty))
	}

	r.endScope()
	if s.Superclass != nil {
		r.endScope()
	}
	r.currentClassTy = orig
	return errors.Join(errs...)
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

	return nil, r.resolveLocal(e.Name)
}

func (r *Resolver) VisitAssignExpr(e expr.AssignExpr) (any, error) {
	_, err := r.resolveExpr(e.Value)
	if err != nil {
		return nil, err
	}
	return nil, r.resolveLocal(e.Name)
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

func (r *Resolver) VisitGetExpr(e expr.GetExpr) (any, error) {
	return r.resolveExpr(e.Obj)
}

func (r *Resolver) VisitSetExpr(e expr.SetExpr) (any, error) {
	_, err := r.resolveExpr(e.Obj)
	if err != nil {
		return nil, err
	}

	return r.resolveExpr(e.Val)
}

func (r *Resolver) VisitThisExpr(e expr.ThisExpr) (any, error) {
	if r.currentClassTy == NONe {
		return nil, resolverError{t: e.Keyword, message: "Can't use 'this' outside of class."}
	}
	return nil, r.resolveLocal(e.Keyword)
}
func (r *Resolver) VisitSuperExpr(e expr.SuperExpr) (any, error) {
	if r.currentClassTy == NONe {
		return nil, resolverError{t: e.Keyword, message: "Can't use super outside of class."}
	}
	if r.currentClassTy != SUBCLASS {
		return nil, resolverError{t: e.Keyword, message: "Can't use super in class with no superclass."}
	}
	return nil, r.resolveLocal(e.Keyword)
}

func (r *Resolver) resolveFunction(f stmt.FunctionStmt, ty functionType) error {
	enclosingTy := r.currentFuncTy

	r.currentFuncTy = ty
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
	err := r.Resolve(f.Body)
	if err != nil {
		return err
	}
	r.endScope()
	r.currentFuncTy = enclosingTy
	return nil
}

func (r *Resolver) resolveLocal(name scanner.Token) error {
	for i := len(r.scopes) - 1; i >= 0; i = i - 1 {
		is_resolved, ok := r.scopes[i][name.Lexeme]
		if ok {
			if is_resolved {
				r.interpreter.Resolve(name, len(r.scopes)-1-i)
				break
			}

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

func (r *Resolver) Resolve(stmts []stmt.Stmt) error {
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
	_, ok := s[n.Lexeme]
	if ok {
		return resolverError{t: n, message: "Already a variable with this name in scope"}
	}
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
