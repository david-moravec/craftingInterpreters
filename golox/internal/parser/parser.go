package parser

import (
	"errors"
	"fmt"

	"github.com/david-moravec/golox/internal/expr"
	"github.com/david-moravec/golox/internal/scanner"
	"github.com/david-moravec/golox/internal/stmt"
)

type parseError struct {
	line    int
	lexeme  string
	message string
}

func newParseError(token scanner.Token, message string) parseError {
	return parseError{token.Line, token.Lexeme, message}
}

func (e parseError) Error() string {
	return fmt.Sprintf("Error [line %d]: Parse error at '%s' %s", e.line, e.lexeme, e.message)

}

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() ([]stmt.Stmt, []error) {
	var stmts []stmt.Stmt
	var errs []error

	for {
		if p.isAtEnd() {
			break
		}

		stmt, err := p.declaration()

		if err != nil {
			errs = append(errs, err)
			p.synchronize()
		}

		stmts = append(stmts, stmt)
	}

	return stmts, errs
}

func (p *Parser) declaration() (stmt.Stmt, error) {
	if p.match(scanner.Class) {
		return p.classDeclaration()
	}
	if p.match(scanner.Fun) {
		return p.funDeclaration("function")
	}
	if p.match(scanner.Var) {
		return p.varDeclaration()
	} else {
		return p.statement()
	}
}

func (p *Parser) classDeclaration() (stmt.Stmt, error) {
	name, err := p.consume(scanner.Identifier, "Expected identifier after 'class'")
	if err != nil {
		return nil, err
	}
	var superclass *expr.VariableExpr = nil
	if p.match(scanner.Less) {
		super_name, err := p.consume(scanner.Identifier, "Expected identifier after '<'")
		if err != nil {
			return nil, err
		}
		superclass = &expr.VariableExpr{Name: *super_name}
	}
	_, err = p.consume(scanner.LeftBrace, "Expected '{' after identifier")
	if err != nil {
		return nil, err
	}
	var methods []stmt.FunctionStmt
	for {
		if p.isAtEnd() || p.checkCurrentKind(scanner.RightBrace) {
			break
		}
		meth, err := p.funDeclaration("method")
		if err != nil {
			return nil, err
		}
		methods = append(methods, meth.(stmt.FunctionStmt))
	}
	_, err = p.consume(scanner.RightBrace, "Expected '}' after method declarations")
	if err != nil {
		return nil, err
	}

	return stmt.ClassStmt{Name: *name, Methods: methods, Superclass: superclass}, nil
}

func (p *Parser) funDeclaration(k string) (stmt.Stmt, error) {
	name, err := p.consume(scanner.Identifier, fmt.Sprintf("Expect %s name.", k))
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.LeftParenthesis, fmt.Sprintf("Expect '(' after %s name.", k))
	var params []scanner.Token
	if !p.checkCurrentKind(scanner.RightParenthesis) {
		for {
			if len(params) >= 255 {
				peek := p.peek()
				return nil, newParseError(peek, "Can't have more than 255 arguments.")
			}
			param, err := p.consume(scanner.Identifier, "Expected parameter name.")
			if err != nil {
				return nil, err
			}
			params = append(params, *param)
			if !p.match(scanner.Comma) {
				break
			}
		}
	}

	p.consume(scanner.RightParenthesis, "Expect ')' after parameters.")
	p.consume(scanner.LeftBrace, fmt.Sprintf("Expect '{' after %s declaration.", k))
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	return stmt.FunctionStmt{Name: *name, Params: params, Body: body.Statements}, nil
}

func (p *Parser) statement() (stmt.Stmt, error) {
	if p.match(scanner.For) {
		return p.forStatement()
	}
	if p.match(scanner.While) {
		return p.whileStatement()
	}
	if p.match(scanner.If) {
		return p.ifStatement()
	}
	if p.match(scanner.Print) {
		return p.printStatement()
	}
	if p.match(scanner.LeftBrace) {
		return p.blockStatement()
	}
	if p.match(scanner.Return) {
		return p.returnStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) returnStatement() (stmt.Stmt, error) {
	keyword := p.previous()
	var val expr.Expr = nil
	var err error
	if !p.checkCurrentKind(scanner.Semicolon) {
		val, err = p.expression()
		if err != nil {
			return nil, err
		}

	}
	err = p.consume_semicolon()
	if err != nil {
		return nil, err
	}
	return stmt.ReturnStmt{Keyword: *keyword, Value: val}, nil
}

func (p *Parser) forStatement() (stmt.Stmt, error) {
	var errs []error
	_, err := p.consume(scanner.LeftParenthesis, "Expect '(' after 'if'.")
	errs = append(errs, err)

	var init stmt.Stmt = nil
	if p.match(scanner.Semicolon) {
		init = nil
	} else if p.match(scanner.Var) {
		i, err := p.varDeclaration()
		if err != nil {
			errs = append(errs, err)
		}
		init = &i
	} else {
		i, err := p.expressionStatement()
		if err != nil {
			errs = append(errs, err)
		}
		init = &i
	}

	var condition expr.Expr = nil
	if !p.checkCurrentKind(scanner.RightParenthesis) {
		if !p.checkCurrentKind(scanner.Semicolon) {
			condition, err = p.expression()
			errs = append(errs, err)
		}
		errs = append(errs, p.consume_semicolon())
	}

	var incr expr.Expr = nil
	if !p.checkCurrentKind(scanner.RightParenthesis) {
		incr, err = p.expression()
		errs = append(errs, err)
	}
	_, err = p.consume(scanner.RightParenthesis, "Expect ')' after for clauses.")
	errs = append(errs, err)
	body, err := p.statement()
	errs = append(errs, err)
	if incr != nil {
		var stmts = []stmt.Stmt{body, stmt.ExpressionStmt{Expression: incr}}
		body = stmt.BlockStmt{Statements: stmts}
	}
	if condition == nil {
		condition = expr.NewLiteral(expr.BoolType, 1, "")
	}
	body = stmt.WhileStmt{Condition: condition, Body: body}
	if init != nil {
		body = stmt.BlockStmt{Statements: []stmt.Stmt{init, body}}
	}
	return body, errors.Join(errs...)
}

func (p *Parser) whileStatement() (stmt.WhileStmt, error) {
	var errs []error
	_, err := p.consume(scanner.LeftParenthesis, "Expect '(' after 'if'.")
	errs = append(errs, err)
	condition, err := p.expression()
	_, err = p.consume(scanner.RightParenthesis, "Expect ')' after 'if'.")
	errs = append(errs, err)
	body, err := p.statement()
	errs = append(errs, err)
	return stmt.WhileStmt{Condition: condition, Body: body}, err
}

func (p *Parser) ifStatement() (stmt.IfStmt, error) {
	var errs []error
	_, err := p.consume(scanner.LeftParenthesis, "Expect '(' after 'if'.")
	errs = append(errs, err)
	condition, err := p.expression()
	_, err = p.consume(scanner.RightParenthesis, "Expect ')' after 'if'.")
	errs = append(errs, err)
	than, err := p.statement()
	errs = append(errs, err)
	var elsebr *stmt.Stmt = nil
	if p.match(scanner.Else) {
		elseb, err := p.statement()
		errs = append(errs, err)
		elsebr = &elseb
	}
	return stmt.IfStmt{Condition: condition, ThenBranch: than, ElseBranch: elsebr}, errors.Join(errs...)
}

func (p *Parser) expressionStatement() (stmt.ExpressionStmt, error) {
	expr, err := p.expression()
	if err != nil {
		return stmt.ExpressionStmt{Expression: expr}, err
	}
	return stmt.ExpressionStmt{Expression: expr}, p.consume_semicolon()

}

func (p *Parser) printStatement() (stmt.PrintStmt, error) {
	expr, err := p.expression()
	if err != nil {
		return stmt.PrintStmt{Expression: expr}, err
	}
	return stmt.PrintStmt{Expression: expr}, p.consume_semicolon()
}

func (p *Parser) blockStatement() (stmt.BlockStmt, error) {
	var stmts []stmt.Stmt
	var errs []error
	for {
		if p.checkCurrentKind(scanner.RightBrace) || p.isAtEnd() {
			break
		}
		s, err := p.declaration()
		errs = append(errs, err)
		stmts = append(stmts, s)
	}
	_, err := p.consume(scanner.RightBrace, "Expect '}' after block.")
	errs = append(errs, err)
	return stmt.BlockStmt{Statements: stmts}, errors.Join(errs...)
}

func (p *Parser) varDeclaration() (stmt.VarStmt, error) {
	name, err := p.consume(scanner.Identifier, "Expect variable name.")
	if err != nil {
		return stmt.VarStmt{Name: nil, Initializer: nil}, err
	}
	var init *expr.Expr = nil
	if p.match(scanner.Equal) {
		e, err := p.expression()
		if err != nil {
			return stmt.VarStmt{Name: name, Initializer: nil}, err
		}
		init = &e
	}
	return stmt.VarStmt{Name: name, Initializer: init}, p.consume_semicolon()

}

func (p Parser) checkCurrentKind(k scanner.TokenKind) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Kind == k
}

func (p Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p Parser) isAtEnd() bool {
	return p.peek().Kind == scanner.TokenKind(scanner.EOF)
}

func (p Parser) previous() *scanner.Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) expression() (expr.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (expr.Expr, error) {
	e, err := p.or()
	if err != nil {
		return e, err
	}
	if p.match(scanner.Equal) {
		eq := p.previous()
		val, err := p.assignment()
		if err != nil {
			return val, err
		}

		switch e.(type) {
		case expr.VariableExpr:
			name := e.(expr.VariableExpr).Name
			return expr.NewAssign(name, val), nil

		case expr.GetExpr:
			get := e.(expr.GetExpr)
			return expr.SetExpr{Obj: get.Obj, Name: get.Name, Val: val}, nil
		}

		return val, newParseError(*eq, "Invalid assignment target.")
	}

	return e, nil
}

func (p *Parser) or() (expr.Expr, error) {
	e, err := p.and()
	for {
		if !p.match(scanner.Or) {
			break
		}
		op := p.previous()
		r, err := p.and()
		if err != nil {
			return nil, err
		}
		e = expr.LogicalExpr{Right: r, Operator: *op, Left: e}
	}
	return e, err
}

func (p *Parser) and() (expr.Expr, error) {
	e, err := p.equality()
	for {
		if !p.match(scanner.And) {
			break
		}
		op := p.previous()
		r, err := p.equality()
		if err != nil {
			return nil, err
		}
		e = expr.LogicalExpr{Right: r, Operator: *op, Left: e}
	}
	return e, err
}

func (p *Parser) equality() (expr.Expr, error) {
	e, err := p.comparison()

	if err != nil {
		return e, err
	}

	for {
		if !p.match(
			scanner.BangEqual,
			scanner.EqualEqual,
		) {
			break
		}
		op := p.previous()
		r, err := p.comparison()

		if err != nil {
			return r, err
		}

		e = expr.NewBinary(e, r, expr.Operator(*op))

	}

	return e, err
}

func (p *Parser) comparison() (expr.Expr, error) {
	e, err := p.term()

	if err != nil {
		return e, err
	}

	for {
		if !p.match(
			scanner.Greater,
			scanner.GreaterEqual,
			scanner.Less,
			scanner.LessEqual,
		) {
			break
		}
		op := p.previous()
		r, err := p.term()

		if err != nil {
			return r, err
		}

		e = expr.NewBinary(e, r, expr.Operator(*op))
	}
	return e, err
}

func (p *Parser) term() (expr.Expr, error) {
	e, err := p.factor()

	if err != nil {
		return e, err
	}

	for {
		if !p.match(
			scanner.Minus,
			scanner.Plus,
		) {
			break
		}
		op := p.previous()
		r, err := p.factor()

		if err != nil {
			return r, err
		}

		e = expr.NewBinary(e, r, expr.Operator(*op))

	}

	return e, err
}

func (p *Parser) factor() (expr.Expr, error) {
	e, err := p.unary()
	if err != nil {
		return e, err
	}
	for {
		if !p.match(
			scanner.Slash,
			scanner.Star,
		) {
			break
		}
		op := p.previous()
		r, err := p.unary()
		if err != nil {
			return e, err
		}
		e = expr.NewBinary(e, r, expr.Operator(*op))
	}
	return e, err

}

func (p *Parser) unary() (expr.Expr, error) {
	if p.match(
		scanner.Bang,
		scanner.Minus,
	) {
		op := p.previous()
		e, err := p.unary()
		if err != nil {
			return e, err
		}
		return expr.NewUnary(expr.Operator(*op), e), nil
	}

	return p.call()
}

func (p *Parser) call() (expr.Expr, error) {
	e, err := p.primary()
	if err != nil {
		return e, err
	}
	for {
		if p.match(scanner.LeftParenthesis) {
			e, err = p.finishCall(e)
		} else if p.match(scanner.Dot) {
			name, err := p.consume(scanner.Identifier, "Expect property name after '.'.")
			if err != nil {
				return nil, err
			}
			e = expr.GetExpr{Obj: e, Name: *name}
		} else {
			break
		}
	}

	return e, err
}

func (p *Parser) finishCall(e expr.Expr) (expr.Expr, error) {
	var args []expr.Expr

	if !p.checkCurrentKind(scanner.RightParenthesis) {
		for {
			if len(args) >= 255 {
				peek := p.peek()
				return e, newParseError(peek, "Can't have more than 255 arguments.")
			}
			arg, err := p.expression()
			if err != nil {
				return arg, err
			}
			args = append(args, arg)
			if !p.match(scanner.Comma) {
				break
			}
		}
	}
	paren, err := p.consume(scanner.RightParenthesis, "Expect ')' after arguments.")
	return expr.CallExpr{Callee: e, Paren: *paren, Arguments: args}, err
}

func (p *Parser) primary() (expr.Expr, error) {
	if p.match(scanner.False) {
		return expr.NewLiteral(expr.BoolType, 0, ""), nil
	} else if p.match(scanner.True) {
		return expr.NewLiteral(expr.BoolType, 1, ""), nil
	} else if p.match(scanner.Nil) {
		return expr.NewLiteral(expr.NilType, 0, ""), nil
	} else if p.match(scanner.Number) {
		return expr.NewLiteral(expr.NumberType, p.previous().Value, ""), nil
	} else if p.match(scanner.String) {
		return expr.NewLiteral(expr.StringType, 0, p.previous().Literal), nil
	} else if p.match(scanner.Identifier) {
		return expr.NewVariable(*p.previous()), nil
	} else if p.match(scanner.This) {
		return expr.ThisExpr{Keyword: *p.previous()}, nil
	} else if p.match(scanner.Super) {
		kw := p.previous()
		p.consume(scanner.Dot, "Expected '.' after 'super'.")
		name, err := p.consume(scanner.Identifier, "Expected method name after '.'.")
		if err != nil {
			return nil, err
		}
		return expr.SuperExpr{Keyword: *kw, Method: *name}, nil
	} else if p.match(scanner.LeftParenthesis) {
		e, err := p.expression()
		if err != nil {
			return e, err
		}
		_, err = p.consume(scanner.RightParenthesis, "Expect ')' after expression.")

		if err != nil {
			return nil, err
		}

		return expr.NewGroup(e), nil
	} else {
		return nil, newParseError(*p.advance(), "Expect expression.")
	}
}

func (p *Parser) consume(k scanner.TokenKind, message string) (*scanner.Token, error) {
	if p.checkCurrentKind(k) {
		return p.advance(), nil

	}

	return nil, newParseError(p.peek(), message)
}

func (p *Parser) consume_semicolon() error {
	_, err := p.consume(scanner.Semicolon, "Expect ';' after expression.")

	return err
}

func (p *Parser) synchronize() {
	p.advance()

	for {
		if p.isAtEnd() {
			break
		}
		if p.previous().Kind == scanner.Semicolon {
			return
		}
		switch p.peek().Kind {
		case
			scanner.Class | scanner.Fun | scanner.Var | scanner.For | scanner.If | scanner.While | scanner.Print | scanner.Return:
			return
		}
		p.advance()
	}

}

func (p *Parser) match(kinds ...scanner.TokenKind) bool {
	for _, k := range kinds {
		if p.checkCurrentKind(k) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) advance() *scanner.Token {
	if p.isAtEnd() {
		return nil
	}

	p.current += 1
	return p.previous()
}
