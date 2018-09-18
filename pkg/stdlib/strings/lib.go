package strings

import "github.com/MontFerret/ferret/pkg/runtime/core"

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"CONCAT":           Concat,
		"CONCAT_SEPARATOR": ConcatWithSeparator,
		"CONTAINS":         Contains,
		"FIND_FIRST":       FindFirst,
		"FIND_LAST":        FindLast,
		"JSON_PARSE":       JsonParse,
		"JSON_STRINGIFY":   JsonStringify,
		"LEFT":             Left,
		"LIKE":             Like,
		"LOWER":            Lower,
		"LTRIM":            LTrim,
		"RIGHT":            RTrim,
		"SPLIT":            Split,
		"SUBSTITUTE":       Substitute,
		"SUBSTRING":        Substring,
		"TRIM":             Trim,
		"UPPER":            Upper,
	}
}
