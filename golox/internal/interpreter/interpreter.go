package interpreter

import (
	"errors"
	"fmt"

	"github.com/david-moravec/golox/internal/expr"
	"github.com/david-moravec/golox/internal/scanner"
	"github.com/david-moravec/golox/internal/stmt"
)

type GoloxError struct {
	t       scanner.Token
	message string
}

func (e GoloxError) Error() string {
	return fmt.Sprintf("[Line: %d]: %s", e.t.Line, e.message)
}

type runtimeError GoloxError

func (e runtimeError) Error() string {
	return fmt.Sprintf("[Line: %d]: %s", e.t.Line, e.message)
}

type unknownTypeError struct {
}

func (e unknownTypeError) Error() string {
	return "Unknown type error"
}

type Interpreter struct {
	env     Environment
	globals Environment
	locals  map[scanner.TokenKey]int
}

func NewInterpreter() Interpreter {
	env := NewGlobalEnvironment()
	env.define("clock", clock{})
	return Interpreter{env: env, globals: env, locals: make(map[scanner.TokenKey]int)}
}

func (i *Interpreter) Interpret(stmts []stmt.Stmt) error {
	var errs []error

	for _, s := range stmts {
		err := i.execute(s)

		if err != nil {
			fmt.Print(err)
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (i *Interpreter) execute(s stmt.Stmt) error {
	return s.Accept(i)
}

func (i *Interpreter) Resolve(name scanner.Token, depth int) {
	i.locals[name.Key()] = depth
}

func (i Interpreter) lookUpVar(e expr.Expr, name scanner.Token) (any, error) {
	dist, ok := i.locals[name.Key()]
	if ok {
		return i.env.getAt(dist, name)
	} else {
		return i.globals.get(name)
	}
}

func (i *Interpreter) executeStmts(b []stmt.Stmt, env Environment) error {
	orig_env := i.env
	i.env = env
	for _, s := range b {
		err := i.execute(s)
		if err != nil {
			return err
		}
	}
	i.env = orig_env
	return nil
}

func (i *Interpreter) evaluate(e expr.Expr) (any, error) {
	return e.Accept(i)
}

func (i Interpreter) VisitLogicalExpr(e expr.LogicalExpr) (any, error) {
	l, err := i.evaluate(e.Left)

	if e.Operator.Kind == scanner.Or {
		if isTruthy(l) {
			return l, err
		}
	} else {
		if !isTruthy(l) {
			return l, err
		}
	}
	return i.evaluate(e.Right)
}

func (i Interpreter) VisitGroupingExpr(e expr.GroupingExpr) (any, error) {
	return i.evaluate(e.Expression)
}

func (i Interpreter) VisitCallExpr(e expr.CallExpr) (any, error) {
	function, err := i.evaluate(e.Callee)
	if err != nil {
		return nil, err
	}
	var args []any
	for _, arg := range e.Arguments {
		a, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, a)
	}
	switch function.(type) {
	case LoxCallable:
		function := function.(LoxCallable)
		if len(args) != function.Arity() {
			return nil, runtimeError{t: e.Paren,
				message: fmt.Sprintf(
					"Expected %d arguments but got %d.", function.Arity(), len(args),
				),
			}

		}

		return function.Call(i, args)
	}
	return nil, runtimeError{t: e.Paren, message: "Can call only functions and classes."}
}

func (i Interpreter) VisitGetExpr(e expr.GetExpr) (any, error) {
	obj, err := i.evaluate(e.Obj)
	if err != nil {
		return nil, err
	}
	switch obj.(type) {
	case *LoxInstance:
		return obj.(*LoxInstance).get(e.Name)
	}

	return nil, runtimeError{t: e.Name, message: "Can get properties only on classes."}
}

func (i Interpreter) VisitSetExpr(e expr.SetExpr) (any, error) {
	obj, err := i.evaluate(e.Obj)
	if err != nil {
		return nil, err
	}
	switch obj.(type) {
	case *LoxInstance:
		val, err := i.evaluate(e.Val)
		if err != nil {
			return nil, err
		}
		return obj.(*LoxInstance).set(e.Name, val)
	}

	return nil, runtimeError{t: e.Name, message: "Can set properties only on classes."}

}

func (i Interpreter) VisitThisExpr(e expr.ThisExpr) (any, error) {
	return i.lookUpVar(e, e.Keyword)
}

func (i Interpreter) VisitSuperExpr(e expr.SuperExpr) (any, error) {
	dist, _ := i.locals[e.Keyword.Key()]
	super, err := i.env.getAt(dist, e.Keyword)
	if err != nil {
		return nil, err
	}
	var superclass *LoxClass = super.(*LoxClass)
	inst, err := i.env.getAt(dist-1, *scanner.DummyThisToken(e.Keyword.Line))
	if err != nil {
		return nil, err
	}
	var instance *LoxInstance
	instance = inst.(*LoxInstance)
	meth := superclass.findMethod(e.Method.Lexeme)
	if meth == nil {
		return nil, runtimeError{t: e.Method, message: fmt.Sprintf("Undefined property %s.", e.Method.Lexeme)}
	}

	return meth.bind(instance), nil
}

func (i Interpreter) VisitLiteralExpr(e expr.LiteralExpr) (any, error) {
	switch e.LitType {
	case expr.NumberType:
		return e.Number, nil
	case expr.StringType:
		return e.Str, nil
	case expr.BoolType:
		return e.Number != 0, nil
	case expr.NilType:
		return nil, nil
	}

	return nil, unknownTypeError{}
}

func (i Interpreter) VisitBinaryExpr(e expr.BinaryExpr) (any, error) {
	l, err := i.evaluate(e.Left)

	if err != nil {
		return nil, err
	}

	r, err := i.evaluate(e.Right)
	if err != nil {
		return nil, err
	}

	var result any

	switch e.Operator.Kind {
	case scanner.Minus:
		err := checkOperandsNumber(scanner.Token(e.Operator), l, r)
		if err != nil {
			return nil, err
		}
		result = l.(float64) - r.(float64)
	case scanner.Plus:
		switch l.(type) {
		case float64:
			switch r.(type) {
			case float64:
				result = l.(float64) + r.(float64)
			}
		case string:
			switch r.(type) {
			case string:
				result = l.(string) + r.(string)
			}
		default:
			return nil, runtimeError{scanner.Token(e.Operator), "Operands must be numbers or strings"}
		}

	case scanner.Slash:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) / r.(float64)
	case scanner.Star:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) * r.(float64)
	case scanner.Greater:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) > r.(float64)
	case scanner.GreaterEqual:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) >= r.(float64)
	case scanner.Less:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) < r.(float64)
	case scanner.LessEqual:
		if err = checkOperandsNumber(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = l.(float64) <= r.(float64)
	case scanner.BangEqual:
		if err = checkOperandsComparable(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = !isEqual(l, r)
	case scanner.EqualEqual:
		if err = checkOperandsComparable(scanner.Token(e.Operator), l, r); err != nil {
			return nil, err
		}
		result = isEqual(l, r)
	}

	return result, nil
}

func (i Interpreter) VisitUnaryExpr(e expr.UnaryExpr) (any, error) {
	a, err := i.evaluate(e.Right)

	if err != nil {
		return nil, err
	}

	switch e.Operator.Kind {
	case scanner.Minus:
		if err := checkOperandNumber(scanner.Token(e.Operator), a); err != nil {
			return nil, err
		}
		return -a.(float64), nil
	case scanner.Bang:
		return !isTruthy(a), nil
	}

	return nil, nil
}

func (i Interpreter) VisitVariableExpr(e expr.VariableExpr) (any, error) {
	return i.lookUpVar(e, e.Name)
}

func (i *Interpreter) VisitAssignExpr(e expr.AssignExpr) (any, error) {
	val, err := i.evaluate(e.Value)

	if err != nil {
		return val, err
	}
	dist, ok := i.locals[e.Name.Key()]
	if ok {
		err = i.env.assignAt(dist, e.Name, val)
	} else {
		err = i.globals.assign(e.Name, val)
	}
	return val, err
}

func (i Interpreter) VisitExpressionStmt(s stmt.ExpressionStmt) error {
	_, err := i.evaluate(s.Expression)

	return err
}

func (i Interpreter) VisitIfStmt(s stmt.IfStmt) error {
	val, err := i.evaluate(s.Condition)
	if err != nil {
		return err
	}
	if isTruthy(val) {
		err = i.execute(s.ThenBranch)
		if err != nil {
			return err
		}
	} else if s.ElseBranch != nil {
		err = i.execute(*s.ElseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i Interpreter) VisitPrintStmt(s stmt.PrintStmt) error {
	v, err := i.evaluate(s.Expression)

	if err != nil {
		return err
	}

	fmt.Println(stringify(v))

	return nil
}

func (i *Interpreter) VisitVarStmt(s stmt.VarStmt) error {
	var val any
	var err error

	if s.Initializer != nil {
		val, err = i.evaluate(*s.Initializer)
		if err != nil {
			return err
		}

	}

	i.env.define(s.Name.Lexeme, val)

	return nil
}

func (i *Interpreter) VisitFunctionStmt(s stmt.FunctionStmt) error {
	fun := LoxFunction{declaration: s, closure: i.env, isInit: false}
	i.env.define(s.Name.Lexeme, fun)
	return nil
}

func (i *Interpreter) VisitReturnStmt(s stmt.ReturnStmt) error {
	val, err := i.evaluate(s.Value)
	if err != nil {
		return err
	}
	return stmt.Return{Value: val}
}

func (i *Interpreter) VisitClassStmt(s stmt.ClassStmt) error {
	i.env.define(s.Name.Lexeme, nil)

	var superclass *LoxClass = nil

	if s.Superclass != nil {
		supercl, err := i.evaluate(s.Superclass)
		if err != nil {
			return err
		}
		switch supercl.(type) {
		case LoxClass:
			supercl := supercl.(LoxClass)
			superclass = &supercl
		default:
			return runtimeError{t: s.Superclass.Name, message: "Superclass must be class."}
		}
	}

	var enclosing Environment = i.env

	if superclass != nil {
		i.env = NewEnvironment(&enclosing)
		i.env.define("super", superclass)
	}

	var methods map[string]LoxFunction = make(map[string]LoxFunction)

	for _, meth := range s.Methods {
		methods[meth.Name.Lexeme] = LoxFunction{
			declaration: meth,
			closure:     i.env,
			isInit:      meth.Name.Lexeme == "init",
		}
	}

	if superclass != nil {
		i.env = enclosing
	}

	klass := LoxClass{Name: s.Name, Methods: methods, Superclass: superclass}

	return i.env.assign(s.Name, klass)
}

func (i *Interpreter) VisitWhileStmt(s stmt.WhileStmt) error {
	for {
		val, err := i.evaluate(s.Condition)
		if err != nil {
			return err
		}
		if !isTruthy(val) {
			break
		} else {
			err = i.execute(s.Body)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Interpreter) VisitBlockStmt(s stmt.BlockStmt) error {
	enclosing := i.env
	return i.executeStmts(s.Statements, NewEnvironment(&enclosing))
}

func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}

	return a == b
}

func isTruthy(a any) bool {
	switch a.(type) {
	case nil:
		return false
	case bool:
		return a.(bool)
	}

	return true
}

func checkOperandNumber(o scanner.Token, r any) error {
	switch r.(type) {
	case float64:
		return nil
	}

	return runtimeError{o, "Operand must be number"}
}

func checkOperandsNumber(o scanner.Token, l any, r any) error {
	switch r.(type) {
	case float64, bool:
		switch l.(type) {
		case float64, bool:
			return nil
		}
	}

	return runtimeError{o, "Operands must be numbers"}
}

func checkOperandsComparable(o scanner.Token, l any, r any) error {
	switch r.(type) {
	case float64, bool, string, nil:
		switch l.(type) {
		case float64, bool, string, nil:
			return nil
		}
	}

	return runtimeError{o, "Operands must be numbers, bools, strings to be comparable"}
}

func stringify(a any) string {
	if a == nil {
		return "nil"
	}
	switch a.(type) {
	case float64:
		return fmt.Sprintf("%f", a)
	case fmt.Stringer:
		return a.(fmt.Stringer).String()
	}

	return fmt.Sprint(a)
}
