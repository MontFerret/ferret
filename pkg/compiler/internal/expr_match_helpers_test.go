package internal

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/parser"
)

func TestSelectMatchConstantFoldExpressionDeclinesOperationalEqualityError(t *testing.T) {
	p := parser.New(`MATCH 1s (1s => 10, _ => 20)`)
	expression := p.Expression()
	if !p.AtEOF() {
		t.Fatal("expected parser to consume the MATCH expression")
	}

	predicate := expression.Predicate()
	if predicate == nil || predicate.ExpressionAtom() == nil {
		t.Fatal("expected MATCH expression atom")
	}

	match := predicate.ExpressionAtom().MatchExpression()
	if match == nil || match.MatchPatternArms() == nil {
		t.Fatal("expected MATCH pattern arms")
	}

	armsContext := match.MatchPatternArms()
	armList := armsContext.MatchPatternArmList()
	if armList == nil {
		t.Fatal("expected MATCH pattern arm list")
	}

	sentinel := errors.New("host equality failed")
	called := 0
	selected, ok := selectMatchConstantFoldExpression(
		matchFoldErrorValue{err: sentinel, called: &called},
		armList.AllMatchPatternArm(),
		armsContext.MatchDefaultArm(),
	)
	if ok || selected != nil {
		t.Fatalf("expected operational equality error to decline folding, got selected=%v ok=%t", selected, ok)
	}
	if called != 1 {
		t.Fatalf("expected exactly one host equality call, got %d", called)
	}
}
