package scanner

import (
	"testing"
)

func test_for_token_kind(t *testing.T, source string, kinds []TokenKind) {
	scanner := Scanner{source: source}

	tokens, err := scanner.ScanTokens()

	if err != nil {
		t.Error("Gor error", err)
	}

	// get rid of the EOF token
	tokens = tokens[:len(tokens)-1]

	if len(kinds) != len(tokens) {
		t.Errorf("Inconsistent number of tokens parsed and token kinds (%d)"+
			"available for check (%d)", len(kinds), len(tokens))
		return
	}

	for i, token := range tokens {
		if token.Kind != kinds[i] {
			t.Errorf("Incorect %d token expected %s, got %s", i, kinds[i], token.Kind)
		}
	}
}

func test_for_literal(t *testing.T, source string, literal string) {
	scanner := Scanner{source: source}

	tokens, err := scanner.ScanTokens()

	if err != nil {
		t.Error("Gor error", err)
	}

	if tokens[0].literal != literal {
		t.Errorf("Incorect expected %s, got %s", literal, tokens[0].literal)
	}

}

func test_for_number(t *testing.T, source string, value float64) {
	scanner := Scanner{source: source}

	tokens, err := scanner.ScanTokens()

	if err != nil {
		t.Error("Gor error", err)
	}

	if tokens[0].value != value {
		t.Errorf("Incorect expected %f, got %f", value, tokens[0].value)
	}

}

func Test_createString(t *testing.T) {
	test_for_literal(t, "\"Ahoj\"", "Ahoj")
}

func Test_createNumber(t *testing.T) {
	test_for_number(t, "88.8", 88.8)
}

func Test_createIdent(t *testing.T) {
	test_for_token_kind(t,
		"ahoj And and While lol",
		[]TokenKind{
			TokenKind(Identifier),
			TokenKind(And),
			TokenKind(Identifier),
			TokenKind(While),
			TokenKind(Identifier),
		})
}
