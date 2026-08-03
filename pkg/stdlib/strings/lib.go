package strings

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("ESCAPE_HTML", EscapeHTML).
		Add("DECODE_URI_COMPONENT", DecodeURIComponent).
		Add("ENCODE_URI_COMPONENT", EncodeURIComponent).
		Add("MD5", Md5).
		Add("SHA1", Sha1).
		Add("SHA512", Sha512).
		Add("TO_BASE64", ToBase64).
		Add("FROM_BASE64", FromBase64).
		Add("JSON_PARSE", JSONParse).
		Add("JSON_STRINGIFY", JSONStringify).
		Add("LOWER", Lower).
		Add("LTRIM", lTrim1).
		Add("UPPER", Upper).
		Add("RANDOM_TOKEN", RandomToken).
		Add("RTRIM", rTrim1).
		Add("TRIM", trim1).
		Add("UNESCAPE_HTML", UnescapeHTML)

	ns.Function().A2().
		Add("CONTAINS", contains2).
		Add("FIND_FIRST", findFirst2).
		Add("FIND_LAST", findLast2).
		Add("LEFT", Left).
		Add("LIKE", like2).
		Add("LTRIM", lTrim2).
		Add("REGEX_MATCH", regexMatch2).
		Add("REGEX_SPLIT", regexSplit2).
		Add("REGEX_TEST", regexTest2).
		Add("RIGHT", Right).
		Add("RTRIM", rTrim2).
		Add("SPLIT", split2).
		Add("SUBSTITUTE", substitute2).
		Add("SUBSTRING", substring2).
		Add("TRIM", trim2)

	ns.Function().A3().
		Add("CONTAINS", contains3).
		Add("FIND_FIRST", findFirst3).
		Add("FIND_LAST", findLast3).
		Add("LIKE", like3).
		Add("REGEX_MATCH", regexMatch3).
		Add("REGEX_REPLACE", regexReplace3).
		Add("REGEX_SPLIT", regexSplit3).
		Add("REGEX_TEST", regexTest3).
		Add("SPLIT", split3).
		Add("SUBSTITUTE", substitute3).
		Add("SUBSTRING", substring3)

	ns.Function().A4().
		Add("FIND_FIRST", findFirst4).
		Add("FIND_LAST", findLast4).
		Add("REGEX_REPLACE", regexReplace4).
		Add("REGEX_SPLIT", regexSplit4).
		Add("SUBSTITUTE", substitute4)

	ns.Function().Var().
		Add("CONCAT", Concat).
		Add("CONCAT_SEPARATOR", ConcatWithSeparator).
		Add("FMT", Fmt)
}
