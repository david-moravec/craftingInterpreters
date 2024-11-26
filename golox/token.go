package main

type TokenKind int

const (
	LeftParenthesis TokenKind = iota
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
	Meaningless
	Bang
	Equal
	Greater
	Less
	//double char
	BangEqual
	EqualEqual
	GreaterEqual
	LessEqual
	Identifier
	String
	Number
	// keyword
	And
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
	EOF
)

var kind_2_string = [...]string{
	"(",
	")",
	"{",
	"}",
	",",
	".",
	"-",
	"+",
	";",
	"\\",
	"*",
	"0",
	"!",
	"=",
	">",
	"<",
	"!=",
	"==",
	">=",
	"<=",
	"Identifier",
	"String",
	"Number",
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
	"EOF",
}

var Keywords = map[string]TokenKind{

	"And":    And,
	"Class":  Class,
	"Else":   Else,
	"False":  False,
	"Fun":    Fun,
	"For":    For,
	"If":     If,
	"Nil":    Nil,
	"Or":     Or,
	"Print":  Print,
	"Return": Return,
	"Super":  Super,
	"This":   This,
	"True":   True,
	"Var":    Var,
	"While":  While,
}

func (t TokenKind) String() string {
	if LeftParenthesis <= t && t <= EOF {
		return kind_2_string[t]
	}

	return "Not a valid literal"
}

type Token struct {
	kind    TokenKind
	lexeme  string
	literal string
	value   float64
	line    int
}

func (t Token) String() string {
	return t.kind.String()
}
