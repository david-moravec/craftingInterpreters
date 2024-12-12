package expr

import (
	"github.com/david-moravec/golox/internal/scanner"
	"testing"
)

func Test_printer(t *testing.T) {
	e := &BinaryExpr{
		Left: &UnaryExpr{
			Operator: Operator(*scanner.NewToken(scanner.Minus, "-", "", 0, 1)),
			Right: LiteralExpr{
				LitType: NumberType,
				Number:  123,
			},
		},
		Right: &GroupingExpr{
			Expression: LiteralExpr{
				LitType: NumberType,
				Number:  45,
			},
		},
		Operator: Operator(*scanner.NewToken(scanner.Star, "*", "", 0, 1)),
	}

	p := AstPrinter{}

	result := p.Print(e)

	ast := "(* (- 123) (group 45))"

	if result != ast {
		t.Errorf("Expected this:\n%s\nGot:\n%s\n", ast, result)
	}
}
