package strings

import "github.com/MontFerret/ferret/pkg/runtime/core"

func RegisterLib(ns core.Namespace) error {
	return ns.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
			"CONCAT":               Concat,
			"CONCAT_SEPARATOR":     ConcatWithSeparator,
			"CONTAINS":             Contains,
			"ESCAPE_HTML":          EscapeHTML,
			"DECODE_URI_COMPONENT": DecodeURIComponent,
			"ENCODE_URI_COMPONENT": EncodeURIComponent,
			"FIND_FIRST":           FindFirst,
			"FIND_LAST":            FindLast,
			"JSON_PARSE":           JSONParse,
			"JSON_STRINGIFY":       JSONStringify,
			"LEFT":                 Left,
			"LIKE":                 Like,
			"LOWER":                Lower,
			"LTRIM":                LTrim,
			"RANDOM_TOKEN":         RandomToken,
			"MD5":                  Md5,
			"REGEX_MATCH":          RegexMatch,
			"REGEX_SPLIT":          RegexSplit,
			"REGEX_TEST":           RegexTest,
			"REGEX_REPLACE":        RegexReplace,
			"RIGHT":                Right,
			"RTRIM":                RTrim,
			"SHA1":                 Sha1,
			"SHA512":               Sha512,
			"SPLIT":                Split,
			"SUBSTITUTE":           Substitute,
			"SUBSTRING":            Substring,
			"TO_BASE64":            ToBase64,
			"FROM_BASE64":          FromBase64,
			"TRIM":                 Trim,
			"UPPER":                Upper,
			"FMT":                  Fmt,
			"UNESCAPE_HTML":        UnescapeHTML,
		}))
}
