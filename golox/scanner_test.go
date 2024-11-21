package main

import (
	"testing"
)

func test_for_literal(t *testing.T, source string, literal string) {
	scanner := Scanner{source: source}

	tokens, err := scanner.scanTokens()

	if err != nil {
		t.Error("Gor error", err)
	}

	if tokens[0].literal != literal {
		t.Errorf("Incorect expected %s, got %s", literal, tokens[0].literal)
	}

}

func test_for_number(t *testing.T, source string, value float64) {
	scanner := Scanner{source: source}

	tokens, err := scanner.scanTokens()

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
