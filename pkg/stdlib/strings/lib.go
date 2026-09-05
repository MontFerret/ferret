package strings

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// RegisterLib registers global string operations.
func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().Add("lower", Lower).Add("upper", Upper).
		Add("trim", trim1).Add("ltrim", lTrim1).Add("rtrim", rTrim1)
	ns.Function().A2().Add("contains", Contains).Add("starts_with", StartsWith).Add("ends_with", EndsWith).
		Add("find_first", findFirst2).Add("find_last", findLast2).
		Add("left", Left).Add("right", Right).Add("substring", substring2).
		Add("trim", trim2).Add("ltrim", lTrim2).Add("rtrim", rTrim2).
		Add("split", split2).Add("join", Join).Add("repeat", Repeat).Add("like", like2).
		Add("regex_test", RegexTest).Add("regex_find", RegexFind).Add("regex_find_all", RegexFindAll).
		Add("regex_split", regexSplit2)
	ns.Function().A3().Add("find_first", findFirst3).Add("find_last", findLast3).
		Add("substring", substring3).Add("split", split3).Add("replace", replace3).
		Add("like", like3).Add("regex_replace", RegexReplace).Add("regex_split", regexSplit3)
	ns.Function().A4().Add("find_first", findFirst4).Add("find_last", findLast4).Add("replace", replace4)
	ns.Function().Var().Add("concat", Concat).Add("fmt", Fmt)
}
