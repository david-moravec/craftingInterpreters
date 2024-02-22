package main

import (
	"fmt"
)

type TokenKind int

type SingleCharacter TokenKind

const (
	LeftParenthesis SingleCharacter = iota
	RightParenthesis
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star
)

type DoubleCharacter TokenKind

const (
	Bang DoubleCharacter = iota
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEquAL
	Less
	LessEqual
)

type Literal TokenKind

const (
	Identifier Literal = iota
	String
	Number
)

var literals = [...]string{
	"Identifier",
	"String",
	"Number",
}

func (l Literal) String() string {
	if Identifier <= l && l <= Number {
		return literals[l]
	}

	return "Not a valid literal"
}

type Keyword TokenKind

const (
	And Keyword = iota
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While
)

var keywords = [...]string{
	"And",
	"Class",
	"Else",
	"False",
	"Fun",
	"For",
	"If",
	"Nil",
	"Or",
	"Print",
	"Return",
	"Super",
	"This",
	"True",
	"Var",
	"While",
}

func (k Keyword) String() string {
	if And <= k && k <= While {
		return keywords[k]
	}

	return "Not a valid keyword"
}

const (
	EOF TokenKind = iota
)

type Token struct {
	kind    TokenKind
	lexeme  string
	literal string
	line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%i")
}
