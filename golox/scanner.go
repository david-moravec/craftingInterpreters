package main

import (
	"errors"
)

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

		tkind, err := s.resolveTokenKind()

		if err != nil {
			errs = append(errs, err)
		} else {
			s.addToken(tkind)
		}
	}

	s.tokens = append(s.tokens, Token{EOF, "", "", s.line})

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

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '0'
	} else {
		return s.source[s.current_char]
	}
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

func (s *Scanner) addToken(token_kind TokenKind) {
	switch token_kind {
	case TokenKind(Meaningless):
		return
	default:
		s.tokens = append(s.tokens, Token{token_kind, s.source[s.lexeme_start:s.current_char], "", s.line})
	}

}

func (s *Scanner) resolveTokenKind() (TokenKind, error) {
	switch s.advance() {
	case '(':
		return TokenKind(LeftParenthesis), nil
	case ')':
		return TokenKind(RightParenthesis), nil
	case '{':
		return TokenKind(LeftBrace), nil
	case '}':
		return TokenKind(RightBrace), nil
	case ',':
		return TokenKind(Comma), nil
	case '.':
		return TokenKind(Dot), nil
	case '-':
		return TokenKind(Minus), nil
	case '+':
		return TokenKind(Plus), nil
	case ';':
		return TokenKind(Semicolon), nil
	case '*':
		return TokenKind(Star), nil
	case '!':
		if s.match('=') {
			return TokenKind(BangEqual), nil
		} else {
			return TokenKind(Bang), nil
		}
	case '=':
		if s.match('=') {
			return TokenKind(EqualEqual), nil
		} else {
			return TokenKind(Equal), nil
		}
	case '<':
		if s.match('=') {
			return TokenKind(LessEqual), nil
		} else {
			return TokenKind(Less), nil
		}
	case '>':
		if s.match('=') {
			return TokenKind(GreaterEqual), nil
		} else {
			return TokenKind(Greater), nil
		}
	case '\\':
		if s.match('\\') {
			for s.peek() != '\n' && s.isAtEnd() {
				s.advance()
			}
			s.line += 1
			return TokenKind(Meaningless), nil
		} else {
			return TokenKind(Slash), nil
		}
	case ' ', '\r', '\t':
		return TokenKind(Meaningless), nil
	case '\n':
		s.line += 1
		return TokenKind(Meaningless), nil
	default:
		return TokenKind(EOF), errors.New("Token not Recognized")
	}
}
