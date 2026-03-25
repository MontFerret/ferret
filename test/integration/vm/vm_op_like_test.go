package vm_test

import (
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	"testing"
)

func TestLikeOperator(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`RETURN "foo" LIKE "f*"`, true),
		S(`RETURN "foo" LIKE "b*"`, false),
		S(`RETURN "foo" NOT LIKE "f*"`, false),
		S(`RETURN "foo" NOT LIKE "b*"`, true),
		S(`LET res = "foo" LIKE  "f*" RETURN res`, true),
		S(`RETURN ("foo" LIKE  "b*") ? "foo" : "bar"`, "bar"),
		S(`RETURN ("foo" NOT LIKE  "b*") ? "foo" : "bar"`, "foo"),
		S(`RETURN true ? ("foo" NOT LIKE  "b*") : false`, true),
		S(`RETURN true ? false : ("foo" NOT LIKE  "b*")`, false),
		S(`RETURN false ? false : ("foo" NOT LIKE  "b*")`, true),
		Array(`FOR str IN ["foo", "bar", "qaz"] FILTER str LIKE "*a*" RETURN str`, []any{"bar", "qaz"}),
		Array(`FOR str IN ["foo", "bar", "qaz"] FILTER str NOT LIKE "*a*" RETURN str`, []any{"foo"}),
		Array(`FOR str IN ["ar", "bar", "qaz"] FILTER str LIKE "a*" RETURN str`, []any{"ar"}),
		Array(`FOR str IN ["ar", "bar", "qaz", "fa", "da"] FILTER str LIKE "*a" RETURN str`, []any{"fa", "da"}),
	})
}
