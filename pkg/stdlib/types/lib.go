package types

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func RegisterLib(ns core.Namespace) error {
	return ns.RegisterFunctions(
		core.NewFunctionsFromMap(map[string]core.Function{
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
			"IS_FLOAT":         IsFloat,
			"IS_STRING":        IsString,
			"IS_DATETIME":      IsDateTime,
			"IS_ARRAY":         IsArray,
			"IS_OBJECT":        IsObject,
			"IS_HTML_ELEMENT":  IsHTMLElement,
			"IS_HTML_DOCUMENT": IsHTMLDocument,
			"IS_BINARY":        IsBinary,
			"IS_NAN":           IsNaN,
			"TYPENAME":         TypeName,
		}))
}

func isTypeof(value core.Value, ctype core.Type) core.Value {
	return values.NewBoolean(core.IsTypeOf(value, ctype))
}
