package main

import (
	"strings"
)

type Scanner struct {
	Source string
}

func (s Scanner) scanTokens() []string {
	return strings.Split(s.Source, " ")
}
