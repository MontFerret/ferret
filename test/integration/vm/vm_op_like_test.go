package vm_test

import (
	"testing"
)

func TestLikeOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(`RETURN "foo" LIKE "f*"`, true),
		Case(`RETURN "foo" LIKE "b*"`, false),
		Case(`RETURN "foo" NOT LIKE "f*"`, false),
		Case(`RETURN "foo" NOT LIKE "b*"`, true),
		Case(`LET res = "foo" LIKE  "f*" RETURN res`, true),
		Case(`RETURN ("foo" LIKE  "b*") ? "foo" : "bar"`, "bar"),
		Case(`RETURN ("foo" NOT LIKE  "b*") ? "foo" : "bar"`, "foo"),
		Case(`RETURN true ? ("foo" NOT LIKE  "b*") : false`, true),
		Case(`RETURN true ? false : ("foo" NOT LIKE  "b*")`, false),
		Case(`RETURN false ? false : ("foo" NOT LIKE  "b*")`, true),
		CaseArray(`FOR str IN ["foo", "bar", "qaz"] FILTER str LIKE "*a*" RETURN str`, []any{"bar", "qaz"}),
		CaseArray(`FOR str IN ["foo", "bar", "qaz"] FILTER str NOT LIKE "*a*" RETURN str`, []any{"foo"}),
		CaseArray(`FOR str IN ["ar", "bar", "qaz"] FILTER str LIKE "a*" RETURN str`, []any{"ar"}),
		CaseArray(`FOR str IN ["ar", "bar", "qaz", "fa", "da"] FILTER str LIKE "*a" RETURN str`, []any{"fa", "da"}),
	})
}
