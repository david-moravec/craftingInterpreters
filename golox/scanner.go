package main

type Scanner struct {
	source       string
	tokens       []Token
	lexeme_start int
	current_char int
	line         int
}

func (s Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.lexeme_start = s.current_char
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", "", s.line})

	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current_char >= len(s.source)
}

func (s *Scanner) scanToken() {}
