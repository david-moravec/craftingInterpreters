package scanner

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

	"and":    And,
	"class":  Class,
	"else":   Else,
	"false":  False,
	"fun":    Fun,
	"for":    For,
	"if":     If,
	"nil":    Nil,
	"or":     Or,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"true":   True,
	"var":    Var,
	"while":  While,
}

func (t TokenKind) String() string {
	if LeftParenthesis <= t && t <= EOF {
		return kind_2_string[t]
	}

	return "Not a valid literal"
}

type Token struct {
	Kind    TokenKind
	Lexeme  string
	Literal string
	Value   float64
	Line    int
}

func NewToken(k TokenKind, l string, lit string, val float64, line int) *Token {
	return &Token{k, l, lit, val, line}
}

func (t Token) String() string {
	return t.Kind.String()
}
