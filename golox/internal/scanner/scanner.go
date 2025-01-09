package scanner

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isAlpha(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isAlphanumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source,
		tokens:       []Token{},
		lexeme_start: 0,
		current_char: 0,
		line:         0}

}

type scannerError struct {
	line   int
	lexeme string
}

func (e scannerError) Error() string {
	return fmt.Sprintf("Error [line %d]: Token '%s' not recognized", e.line, e.lexeme)
}

type Scanner struct {
	source       string
	tokens       []Token
	lexeme_start int
	current_char int
	line         int
}

func (s Scanner) ScanTokens() ([]Token, []error) {
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

	return s.tokens, errs
}

func (s Scanner) isAtEnd() bool {
	return s.current_char >= len(s.source)
}

func (s Scanner) peek() (byte, error) {
	if s.isAtEnd() {
		return 'x', errors.New("End of file reached")
	} else {
		return s.source[s.current_char], nil
	}
}

func (s Scanner) peekNext() (byte, error) {
	if s.current_char+1 > len(s.source) {
		return 'x', errors.New("EOF reached")
	}

	return s.source[s.current_char+1], nil
}

func (s Scanner) peekCurrentLexeme() string {
	return s.source[s.lexeme_start:s.current_char]

}

func (s *Scanner) advance() byte {
	current := s.source[s.current_char]
	s.current_char += 1
	return current
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
	case Meaningless:
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
		return EOF, "", math.NaN(), errors.New("Unterminated string")
	}

	//closing "
	s.advance()

	literal := s.source[s.lexeme_start+1 : s.current_char-1]

	return String, literal, math.NaN(), nil
}

func (s *Scanner) createIdentifier() (TokenKind, string, float64, error) {
	for {
		c, err := s.peek()

		if err != nil {
			break
		}

		if isAlphanumeric(c) {
			s.advance()
		} else {
			break
		}
	}

	current := s.source[s.lexeme_start:s.current_char]

	t, ok := Keywords[current]

	if !ok {
		return Identifier, "", math.NaN(), nil
	}

	return t, "", math.NaN(), nil

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
		return Number, "", 0, err
	}
	next_c, err := s.peekNext()
	if err != nil {
		return Number, "", 0, err
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

	return Number, "", value, err
}

func (s *Scanner) resolveTokenKind() (TokenKind, string, float64, error) {
	c := s.advance()

	switch c {
	case '(':
		return LeftParenthesis, "", math.NaN(), nil
	case ')':
		return RightParenthesis, "", math.NaN(), nil
	case '{':
		return LeftBrace, "", math.NaN(), nil
	case '}':
		return RightBrace, "", math.NaN(), nil
	case ',':
		return Comma, "", math.NaN(), nil
	case '.':
		return Dot, "", math.NaN(), nil
	case '-':
		return Minus, "", math.NaN(), nil
	case '+':
		return Plus, "", math.NaN(), nil
	case ';':
		return Semicolon, "", math.NaN(), nil
	case '*':
		return Star, "", math.NaN(), nil
	case '!':
		if s.match('=') {
			return BangEqual, "", math.NaN(), nil
		} else {
			return Bang, "", math.NaN(), nil
		}
	case '=':
		if s.match('=') {
			return EqualEqual, "", math.NaN(), nil
		} else {
			return Equal, "", math.NaN(), nil
		}
	case '<':
		if s.match('=') {
			return LessEqual, "", math.NaN(), nil
		} else {
			return Less, "", math.NaN(), nil
		}
	case '>':
		if s.match('=') {
			return GreaterEqual, "", math.NaN(), nil
		} else {
			return Greater, "", math.NaN(), nil
		}
	case '\\':
		if s.match('\\') {
			for {
				c, err := s.peek()

				if err != nil {
					return EOF, "", math.NaN(), err
				}

				if c == '\n' {
					break
				}
				s.advance()
			}

			s.line += 1
			return Meaningless, "", math.NaN(), nil
		} else {
			return Slash, "", math.NaN(), nil
		}
	case ' ', '\r', '\t':
		return Meaningless, "", math.NaN(), nil
	case '\n':
		s.line += 1
		return Meaningless, "", math.NaN(), nil
	case '"':
		return s.createString()
	default:
		if isDigit(c) {
			return s.createNumber()
		}
		if isAlphanumeric(c) {
			return s.createIdentifier()
		}
		return EOF, "", math.NaN(), scannerError{s.line, s.peekCurrentLexeme()}
	}
}
