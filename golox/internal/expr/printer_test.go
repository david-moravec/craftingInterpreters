package expr

import (
	"github.com/david-moravec/golox/internal/scanner"
	"testing"
)

func Test_printer(t *testing.T) {
	e := &BinaryExpr{
		left: &UnaryExpr{
			operator: Operator(*scanner.NewToken(scanner.Minus, "-", "", 0, 1)),
			right: LiteralExpr{
				litType: NumberType,
				number:  123,
			},
		},
		right: &GroupingExpr{
			expression: LiteralExpr{
				litType: NumberType,
				number:  45,
			},
		},
		operator: Operator(*scanner.NewToken(scanner.Star, "*", "", 0, 1)),
	}

	p := AstPrinter{}

	result := p.Print(e)

	ast := "(* (- 123) (group 45))"

	if result != ast {
		t.Errorf("Expected this:\n%s\nGot:\n%s\n", ast, result)
	}
}
