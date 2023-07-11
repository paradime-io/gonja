package parser

import (
	log "github.com/sirupsen/logrus"

	"github.com/paradime-io/gonja/nodes"
)

func (p *Parser) ParseTest(expr nodes.Expression) (nodes.Expression, error) {
	log.WithFields(log.Fields{
		"current": p.Current(),
	}).Trace("parseTest")

	expr, err := p.ParseFilterExpression(expr)
	if err != nil {
		return nil, err
	}

	if p.MatchName("is") != nil {
		not := p.MatchName("not")
		if p.End() {
			return nil, p.Error("is statement is incomplete", nil)
		}
		ident := p.Next()

		test := &nodes.TestCall{
			Token:  ident,
			Name:   ident.Val,
			Args:   []nodes.Expression{},
			Kwargs: map[string]nodes.Expression{},
		}

		if _, needsRightSide := p.Config.TestsNeedingRightSide[ident.Val]; needsRightSide {
			arg, argErr := p.ParseVariableOrLiteral()
			if argErr != nil {
				return nil, argErr
			}
			test.Args = append(test.Args, arg)
		}
		expr = &nodes.TestExpression{
			Expression: expr,
			Test:       test,
		}

		if not != nil {
			expr = &nodes.Negation{expr, not}
		}
	}

	log.WithFields(log.Fields{
		"expr": expr,
	}).Trace("parseTest return")
	return expr, nil
}
