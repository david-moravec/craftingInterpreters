package parser

import (
	"github.com/david-moravec/golox/internal/expr"
	"github.com/david-moravec/golox/internal/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) expression() expr.Expr {
	return p.equality()
}

func (p *Parser) equality() expr.Expr {
	e := p.comparison()
	for {
		if !p.match(
			scanner.BangEqual,
			scanner.EqualEqual,
		) {
			break
		}
		op := p.previous()
		r := p.comparison()
		e = expr.NewBinary(e, r, expr.Operator(*op))

	}
	return e
}

func (p *Parser) comparison() expr.Expr {
	e := p.term()
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
		r := p.term()
		e = expr.NewBinary(e, r, expr.Operator(*op))
	}
	return e
}

func (p *Parser) term() expr.Expr {
	e := p.factor()
	for {
		if !p.match(
			scanner.Minus,
			scanner.Plus,
		) {
			break
		}
		op := p.previous()
		r := p.factor()
		e = expr.NewBinary(e, r, expr.Operator(*op))
	}
	return e
}

func (p *Parser) factor() expr.Expr {
	e := p.unary()
	for {
		if !p.match(
			scanner.Slash,
			scanner.Star,
		) {
			break
		}
		op := p.previous()
		r := p.unary()
		e = expr.NewBinary(e, r, expr.Operator(*op))
	}
	return e

}

func (p *Parser) unary() expr.Expr {
	if p.match(
		scanner.Bang,
		scanner.Minus,
	) {
		return expr.NewUnary(expr.Operator(*p.previous()), p.unary())
	}

	return p.primary()
}

func (p *Parser) primary() expr.Expr {
	if p.match(scanner.False) {
		return expr.NewLiteral(expr.BoolType, 0, "")
	} else if p.match(scanner.True) {
		return expr.NewLiteral(expr.BoolType, 1, "")
	} else if p.match(scanner.Nil) {
		return expr.NewLiteral(expr.NilType, 0, "")
	} else if p.match(scanner.Number) {
		return expr.NewLiteral(expr.NumberType, p.previous().Value, "")
	} else if p.match(scanner.String) {
		return expr.NewLiteral(expr.NumberType, 0, p.previous().Literal)
	} else if p.match(scanner.LeftParenthesis) {
		e := p.expression()
		p.consume(scanner.RightParenthesis, "Expect ')' after expression.")
		return expr.NewGroup(e)
	}

	// TODO: return error
	panic("")
}

func (p *Parser) consume(k scanner.TokenKind, message string) *scanner.Token {
	if p.checkCurrentKind(k) {
		return p.advance()
	}

	panic(message)
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

func (p *Parser) checkCurrentKind(k scanner.TokenKind) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Kind == k
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Kind == scanner.TokenKind(scanner.EOF)
}

func (p *Parser) previous() *scanner.Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) advance() *scanner.Token {
	if p.isAtEnd() {
		return nil
	}

	p.current += 1
	return p.previous()
}
