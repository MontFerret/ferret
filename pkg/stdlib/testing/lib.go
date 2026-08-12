package testing

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// @namespace T
func RegisterLib(ns runtime.Namespace) {
	t := ns.Namespace("T")

	registerNOT(t)

	registerPositive(t, "empty", Empty)
	registerPositive(t, "eq", Equal)
	registerPositive(t, "fail", Fail)
	registerPositive(t, "false", False)
	registerPositive(t, "gt", Gt)
	registerPositive(t, "gte", Gte)
	registerPositive(t, "include", Include)
	registerPositive(t, "len", Len)
	registerPositive(t, "match", Match)
	registerPositive(t, "lt", Lt)
	registerPositive(t, "lte", Lte)
	registerPositive(t, "none", None)
	registerPositive(t, "true", True)
	registerPositive(t, "string", String)
	registerPositive(t, "int", Int)
	registerPositive(t, "float", Float)
	registerPositive(t, "datetime", DateTime)
	registerPositive(t, "array", Array)
	registerPositive(t, "object", Object)
	registerPositive(t, "binary", Binary)
}

func registerNOT(ns runtime.Namespace) {
	t := ns.Namespace("NOT")

	registerNegative(t, "empty", Empty)
	registerNegative(t, "eq", Equal)
	registerNegative(t, "false", False)
	registerNegative(t, "gt", Gt)
	registerNegative(t, "gte", Gte)
	registerNegative(t, "include", Include)
	registerNegative(t, "len", Len)
	registerNegative(t, "match", Match)
	registerNegative(t, "lt", Lt)
	registerNegative(t, "lte", Lte)
	registerNegative(t, "none", None)
	registerNegative(t, "true", True)
	registerNegative(t, "string", String)
	registerNegative(t, "int", Int)
	registerNegative(t, "float", Float)
	registerNegative(t, "datetime", DateTime)
	registerNegative(t, "array", Array)
	registerNegative(t, "object", Object)
	registerNegative(t, "binary", Binary)
}
