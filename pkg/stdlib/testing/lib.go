package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

// @namespace T
func RegisterLib(ns runtime.Namespace) {
	t := ns.Namespace("T")

	registerNOT(t)

	t.Function().Var().
		Add("EQ", base.NewPositiveAssertion(Equal)).
		Add("FAIL", base.NewPositiveAssertion(Fail)).
		Add("FALSE", base.NewPositiveAssertion(False)).
		Add("GT", base.NewPositiveAssertion(Gt)).
		Add("GTE", base.NewPositiveAssertion(Gte)).
		Add("INCLUDE", base.NewPositiveAssertion(Include)).
		Add("LEN", base.NewPositiveAssertion(Len)).
		Add("MATCH", base.NewPositiveAssertion(Match)).
		Add("LT", base.NewPositiveAssertion(Lt)).
		Add("LTE", base.NewPositiveAssertion(Lte)).
		Add("NONE", base.NewPositiveAssertion(None)).
		Add("TRUE", base.NewPositiveAssertion(True)).
		Add("STRING", base.NewPositiveAssertion(String)).
		Add("INT", base.NewPositiveAssertion(Int)).
		Add("FLOAT", base.NewPositiveAssertion(Float)).
		Add("DATETIME", base.NewPositiveAssertion(DateTime)).
		Add("ARRAY", base.NewPositiveAssertion(Array)).
		Add("OBJECT", base.NewPositiveAssertion(Object)).
		Add("BINARY", base.NewPositiveAssertion(Binary))
}

func registerNOT(ns runtime.Namespace) {
	t := ns.Namespace("NOT")

	t.Function().Var().
		Add("EMPTY", base.NewNegativeAssertion(Empty)).
		Add("EQ", base.NewNegativeAssertion(Equal)).
		Add("FALSE", base.NewNegativeAssertion(False)).
		Add("GT", base.NewNegativeAssertion(Gt)).
		Add("GTE", base.NewNegativeAssertion(Gte)).
		Add("INCLUDE", base.NewNegativeAssertion(Include)).
		Add("LEN", base.NewNegativeAssertion(Len)).
		Add("MATCH", base.NewNegativeAssertion(Match)).
		Add("LT", base.NewNegativeAssertion(Lt)).
		Add("LTE", base.NewNegativeAssertion(Lte)).
		Add("NONE", base.NewNegativeAssertion(None)).
		Add("TRUE", base.NewNegativeAssertion(True)).
		Add("STRING", base.NewNegativeAssertion(String)).
		Add("INT", base.NewNegativeAssertion(Int)).
		Add("FLOAT", base.NewNegativeAssertion(Float)).
		Add("DATETIME", base.NewNegativeAssertion(DateTime)).
		Add("ARRAY", base.NewNegativeAssertion(Array)).
		Add("OBJECT", base.NewNegativeAssertion(Object)).
		Add("BINARY", base.NewNegativeAssertion(Binary))
}
