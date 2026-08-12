package strings

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("escape_html", EscapeHTML).
		Add("decode_uri_component", DecodeURIComponent).
		Add("encode_uri_component", EncodeURIComponent).
		Add("md5", Md5).
		Add("sha1", Sha1).
		Add("sha512", Sha512).
		Add("to_base64", ToBase64).
		Add("from_base64", FromBase64).
		Add("json_parse", JSONParse).
		Add("json_stringify", JSONStringify).
		Add("lower", Lower).
		Add("ltrim", lTrim1).
		Add("upper", Upper).
		Add("random_token", RandomToken).
		Add("rtrim", rTrim1).
		Add("trim", trim1).
		Add("unescape_html", UnescapeHTML)

	ns.Function().A2().
		Add("contains", contains2).
		Add("find_first", findFirst2).
		Add("find_last", findLast2).
		Add("left", Left).
		Add("like", like2).
		Add("ltrim", lTrim2).
		Add("regex_match", regexMatch2).
		Add("regex_split", regexSplit2).
		Add("regex_test", regexTest2).
		Add("right", Right).
		Add("rtrim", rTrim2).
		Add("split", split2).
		Add("substitute", substitute2).
		Add("substring", substring2).
		Add("trim", trim2)

	ns.Function().A3().
		Add("contains", contains3).
		Add("find_first", findFirst3).
		Add("find_last", findLast3).
		Add("like", like3).
		Add("regex_match", regexMatch3).
		Add("regex_replace", regexReplace3).
		Add("regex_split", regexSplit3).
		Add("regex_test", regexTest3).
		Add("split", split3).
		Add("substitute", substitute3).
		Add("substring", substring3)

	ns.Function().A4().
		Add("find_first", findFirst4).
		Add("find_last", findLast4).
		Add("regex_replace", regexReplace4).
		Add("regex_split", regexSplit4).
		Add("substitute", substitute4)

	ns.Function().Var().
		Add("concat", Concat).
		Add("concat_separator", ConcatWithSeparator).
		Add("fmt", Fmt)
}
