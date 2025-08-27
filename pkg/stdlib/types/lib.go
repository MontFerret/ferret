package types

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

// TODO: Migrate to Func1 function types
func RegisterLib(ns runtime.Namespace) error {
	return ns.RegisterFunctions(
		runtime.NewFunctionsFromMap(map[string]runtime.Function{
			"TO_BOOL":          ToBool,
			"TO_INT":           ToInt,
			"TO_FLOAT":         ToFloat,
			"TO_STRING":        ToString,
			"TO_DATETIME":      ToDateTime,
			"TO_ARRAY":         ToArray,
			"TO_BINARY":        ToBinary,
			"IS_NONE":          IsNone,
			"IS_BOOL":          IsBool,
			"IS_INT":           IsInt,
			"IS_LIST":          IsList,
			"IS_FLOAT":         IsFloat,
			"IS_STRING":        IsString,
			"IS_DATETIME":      IsDateTime,
			"IS_ARRAY":         IsArray,
			"IS_OBJECT":        IsObject,
			"IS_HTML_ELEMENT":  IsHTMLElement,
			"IS_HTML_DOCUMENT": IsHTMLDocument,
			"IS_BINARY":        IsBinary,
			"IS_NAN":           IsNaN,
		}))
}

func isTypeof(value runtime.Value, ctype runtime.Type) runtime.Value {
	return runtime.NewBoolean(runtime.Reflect(value) == ctype)
}
