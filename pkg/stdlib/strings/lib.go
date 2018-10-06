package strings

import "github.com/MontFerret/ferret/pkg/runtime/core"

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"CONCAT":               Concat,
		"CONCAT_SEPARATOR":     ConcatWithSeparator,
		"CONTAINS":             Contains,
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
		"REGEXP_MATCH":         RegexMatch,
		"REGEXP_SPLIT":         RegexSplit,
		"REGEXP_TEST":          RegexTest,
		"REGEXP_REPLACE":       RegexReplace,
		"REVERSE":              Reverse,
		"RIGHT":                Right,
		"RTRIM":                RTrim,
		"SHA1":                 Sha1,
		"SHA512":               Sha512,
		"SPLIT":                Split,
		"SUBSTITUTE":           Substitute,
		"SUBSTRING":            Substring,
		"TO_BASE64":            ToBase64,
		"TRIM":                 Trim,
		"UPPER":                Upper,
	}
}
