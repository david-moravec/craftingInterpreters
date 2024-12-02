package expr

import (
	"testing"
)

func Test_printer(t *testing.T) {
	e := &BinaryExpr{
		left: &UnaryExpr{
			operator: '-',
			right: LiteralExpr{
				lit_type: LiteralType(NumberType),
				number:   123,
			},
		},
		right: &GroupingExpr{
			expression: LiteralExpr{
				lit_type: LiteralType(NumberType),
				number:   45,
			},
		},
		operator: '*',
	}

	p := AstPrinter{}

	result := p.print(e)

	ast := "(* (- 123) (group 45))"

	if result != ast {
		t.Errorf("Expected this:\n%s\nGot:\n%s\n", ast, result)
	}
}
