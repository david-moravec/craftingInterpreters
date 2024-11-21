package main

import (
	"errors"
	"math"
	"strconv"
)

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

type Scanner struct {
	source       string
	tokens       []Token
	lexeme_start int
	current_char int
	line         int
}

func (s Scanner) scanTokens() ([]Token, error) {
	var errs []error

	for !s.isAtEnd() {
		s.lexeme_start = s.current_char

		tkind, literal, value, err := s.resolveTokenKind()

		if err != nil {
			errs = append(errs, err)
		} else {
			s.addToken(tkind, literal, value)
		}
	}

	s.tokens = append(s.tokens, Token{EOF, "", "", math.NaN(), s.line})

	return s.tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current_char >= len(s.source)
}

func (s *Scanner) advance() byte {
	current := s.source[s.current_char]
	s.current_char += 1
	return current
}

func (s *Scanner) peek() (byte, error) {
	if s.isAtEnd() {
		return 'x', errors.New("End of file reached")
	} else {
		return s.source[s.current_char], nil
	}
}

func (s *Scanner) peekNext() (byte, error) {
	if s.current_char+1 > len(s.source) {
		return 'x', errors.New("EOF reached")
	}

	return s.source[s.current_char+1], nil
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current_char] != expected {
		return false
	}

	s.current_char += 1
	return true
}

func (s *Scanner) addToken(token_kind TokenKind, literal string, value float64) {
	switch token_kind {
	case TokenKind(Meaningless):
		return
	default:
		s.tokens = append(s.tokens, Token{token_kind, s.source[s.lexeme_start:s.current_char], literal, value, s.line})
	}

}

func (s *Scanner) createString() (TokenKind, string, float64, error) {
	for {
		c, err := s.peek()

		if err != nil {
			break
		}

		if c == '"' {
			break
		}

		if c == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		return TokenKind(EOF), "", math.NaN(), errors.New("Unterminated string")
	}

	//closing "
	s.advance()

	literal := s.source[s.lexeme_start+1 : s.current_char-1]

	return TokenKind(String), literal, math.NaN(), nil
}

func (s *Scanner) createNumber() (TokenKind, string, float64, error) {
	for {
		c, err := s.peek()
		if err != nil || !isDigit(c) {
			break
		}
		s.advance()
	}
	c, err := s.peek()
	if err != nil {
		return TokenKind(Number), "", 0, err
	}
	next_c, err := s.peekNext()
	if err != nil {
		return TokenKind(Number), "", 0, err
	}
	if c == '.' && isDigit(next_c) {
		s.advance()
		for {
			c, err := s.peek()
			if err != nil || !isDigit(c) {
				break
			}
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.source[s.lexeme_start:s.current_char], 64)

	return TokenKind(Number), "", value, err
}

func (s *Scanner) resolveTokenKind() (TokenKind, string, float64, error) {
	c := s.advance()

	switch c {
	case '(':
		return TokenKind(LeftParenthesis), "", math.NaN(), nil
	case ')':
		return TokenKind(RightParenthesis), "", math.NaN(), nil
	case '{':
		return TokenKind(LeftBrace), "", math.NaN(), nil
	case '}':
		return TokenKind(RightBrace), "", math.NaN(), nil
	case ',':
		return TokenKind(Comma), "", math.NaN(), nil
	case '.':
		return TokenKind(Dot), "", math.NaN(), nil
	case '-':
		return TokenKind(Minus), "", math.NaN(), nil
	case '+':
		return TokenKind(Plus), "", math.NaN(), nil
	case ';':
		return TokenKind(Semicolon), "", math.NaN(), nil
	case '*':
		return TokenKind(Star), "", math.NaN(), nil
	case '!':
		if s.match('=') {
			return TokenKind(BangEqual), "", math.NaN(), nil
		} else {
			return TokenKind(Bang), "", math.NaN(), nil
		}
	case '=':
		if s.match('=') {
			return TokenKind(EqualEqual), "", math.NaN(), nil
		} else {
			return TokenKind(Equal), "", math.NaN(), nil
		}
	case '<':
		if s.match('=') {
			return TokenKind(LessEqual), "", math.NaN(), nil
		} else {
			return TokenKind(Less), "", math.NaN(), nil
		}
	case '>':
		if s.match('=') {
			return TokenKind(GreaterEqual), "", math.NaN(), nil
		} else {
			return TokenKind(Greater), "", math.NaN(), nil
		}
	case '\\':
		if s.match('\\') {
			for {
				c, err := s.peek()

				if err != nil {
					return TokenKind(EOF), "", math.NaN(), err
				}

				if c == '\n' {
					break
				}
				s.advance()
			}

			s.line += 1
			return TokenKind(Meaningless), "", math.NaN(), nil
		} else {
			return TokenKind(Slash), "", math.NaN(), nil
		}
	case ' ', '\r', '\t':
		return TokenKind(Meaningless), "", math.NaN(), nil
	case '\n':
		s.line += 1
		return TokenKind(Meaningless), "", math.NaN(), nil
	case '"':
		return s.createString()
	default:
		if isDigit(c) {
			return s.createNumber()
		}
		return TokenKind(EOF), "", math.NaN(), errors.New("Token not Recognized")
	}
}
