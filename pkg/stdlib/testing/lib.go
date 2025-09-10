package testing

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

// @namespace T
func RegisterLib(ns runtime.Namespace) error {
	t := ns.Namespace("T")

	if err := registerNOT(t); err != nil {
		return err
	}

	t.Functions().
		Set("EQ", base.NewPositiveAssertion(Equal)).
		Set("FAIL", base.NewPositiveAssertion(Fail)).
		Set("FALSE", base.NewPositiveAssertion(False)).
		Set("GT", base.NewPositiveAssertion(Gt)).
		Set("GTE", base.NewPositiveAssertion(Gte)).
		Set("INCLUDE", base.NewPositiveAssertion(Include)).
		Set("LEN", base.NewPositiveAssertion(Len)).
		Set("MATCH", base.NewPositiveAssertion(Match)).
		Set("LT", base.NewPositiveAssertion(Lt)).
		Set("LTE", base.NewPositiveAssertion(Lte)).
		Set("NONE", base.NewPositiveAssertion(None)).
		Set("TRUE", base.NewPositiveAssertion(True)).
		Set("STRING", base.NewPositiveAssertion(String)).
		Set("INT", base.NewPositiveAssertion(Int)).
		Set("FLOAT", base.NewPositiveAssertion(Float)).
		Set("DATETIME", base.NewPositiveAssertion(DateTime)).
		Set("ARRAY", base.NewPositiveAssertion(Array)).
		Set("OBJECT", base.NewPositiveAssertion(Object)).
		Set("BINARY", base.NewPositiveAssertion(Binary))

	return nil
}

func registerNOT(ns runtime.Namespace) error {
	t := ns.Namespace("NOT")

	t.Functions().
		Set("EMPTY", base.NewNegativeAssertion(Empty)).
		Set("EQ", base.NewNegativeAssertion(Equal)).
		Set("FALSE", base.NewNegativeAssertion(False)).
		Set("GT", base.NewNegativeAssertion(Gt)).
		Set("GTE", base.NewNegativeAssertion(Gte)).
		Set("INCLUDE", base.NewNegativeAssertion(Include)).
		Set("LEN", base.NewNegativeAssertion(Len)).
		Set("MATCH", base.NewNegativeAssertion(Match)).
		Set("LT", base.NewNegativeAssertion(Lt)).
		Set("LTE", base.NewNegativeAssertion(Lte)).
		Set("NONE", base.NewNegativeAssertion(None)).
		Set("TRUE", base.NewNegativeAssertion(True)).
		Set("STRING", base.NewNegativeAssertion(String)).
		Set("INT", base.NewNegativeAssertion(Int)).
		Set("FLOAT", base.NewNegativeAssertion(Float)).
		Set("DATETIME", base.NewNegativeAssertion(DateTime)).
		Set("ARRAY", base.NewNegativeAssertion(Array)).
		Set("OBJECT", base.NewNegativeAssertion(Object)).
		Set("BINARY", base.NewNegativeAssertion(Binary))

	return nil
}
