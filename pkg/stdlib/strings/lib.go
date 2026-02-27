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
		Add("UPPER", Upper).
		Add("RANDOM_TOKEN", RandomToken).
		Add("UNESCAPE_HTML", UnescapeHTML)

	ns.Function().A2().
		Add("LEFT", Left).
		Add("RIGHT", Right)

	ns.Function().Var().
		Add("CONCAT", Concat).
		Add("CONCAT_SEPARATOR", ConcatWithSeparator).
		Add("CONTAINS", Contains).
		Add("FIND_FIRST", FindFirst).
		Add("FIND_LAST", FindLast).
		Add("LIKE", Like).
		Add("LTRIM", LTrim).
		Add("REGEX_MATCH", RegexMatch).
		Add("REGEX_SPLIT", RegexSplit).
		Add("REGEX_TEST", RegexTest).
		Add("REGEX_REPLACE", RegexReplace).
		Add("RTRIM", RTrim).
		Add("SPLIT", Split).
		Add("SUBSTITUTE", Substitute).
		Add("SUBSTRING", Substring).
		Add("TRIM", Trim).
		Add("FMT", Fmt)
}
