package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// @namespace T
func RegisterLib(ns runtime.Namespace) {
	t := ns.Namespace("T")

	registerNOT(t)

	registerPositive(t, "EMPTY", Empty)
	registerPositive(t, "EQ", Equal)
	registerPositive(t, "FAIL", Fail)
	registerPositive(t, "FALSE", False)
	registerPositive(t, "GT", Gt)
	registerPositive(t, "GTE", Gte)
	registerPositive(t, "INCLUDE", Include)
	registerPositive(t, "LEN", Len)
	registerPositive(t, "MATCH", Match)
	registerPositive(t, "LT", Lt)
	registerPositive(t, "LTE", Lte)
	registerPositive(t, "NONE", None)
	registerPositive(t, "TRUE", True)
	registerPositive(t, "STRING", String)
	registerPositive(t, "INT", Int)
	registerPositive(t, "FLOAT", Float)
	registerPositive(t, "DATETIME", DateTime)
	registerPositive(t, "ARRAY", Array)
	registerPositive(t, "OBJECT", Object)
	registerPositive(t, "BINARY", Binary)
}

func registerNOT(ns runtime.Namespace) {
	t := ns.Namespace("NOT")

	registerNegative(t, "EMPTY", Empty)
	registerNegative(t, "EQ", Equal)
	registerNegative(t, "FALSE", False)
	registerNegative(t, "GT", Gt)
	registerNegative(t, "GTE", Gte)
	registerNegative(t, "INCLUDE", Include)
	registerNegative(t, "LEN", Len)
	registerNegative(t, "MATCH", Match)
	registerNegative(t, "LT", Lt)
	registerNegative(t, "LTE", Lte)
	registerNegative(t, "NONE", None)
	registerNegative(t, "TRUE", True)
	registerNegative(t, "STRING", String)
	registerNegative(t, "INT", Int)
	registerNegative(t, "FLOAT", Float)
	registerNegative(t, "DATETIME", DateTime)
	registerNegative(t, "ARRAY", Array)
	registerNegative(t, "OBJECT", Object)
	registerNegative(t, "BINARY", Binary)
}
