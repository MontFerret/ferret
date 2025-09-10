package types

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) error {
	ns.Functions().
		Set1("TO_BOOL", ToBool).
		Set1("TO_INT", ToInt).
		Set1("TO_FLOAT", ToFloat).
		Set1("TO_STRING", ToString).
		Set1("TO_DATETIME", ToDateTime).
		Set1("TO_ARRAY", ToArray).
		Set1("TO_BINARY", ToBinary).
		Set1("TO_OBJECT", ToObject).
		Set1("IS_NONE", IsNone).
		Set1("IS_BOOL", IsBool).
		Set1("IS_INT", IsInt).
		Set1("IS_FLOAT", IsFloat).
		Set1("IS_STRING", IsString).
		Set1("IS_DATETIME", IsDateTime).
		Set1("IS_LIST", IsList).
		Set1("IS_ARRAY", IsArray).
		Set1("IS_MAP", IsMap).
		Set1("IS_OBJECT", IsObject).
		Set1("IS_HTML_ELEMENT", IsHTMLElement).
		Set1("IS_HTML_DOCUMENT", IsHTMLDocument).
		Set1("IS_BINARY", IsBinary).
		Set1("IS_NAN", IsNaN)

	return nil
}

func isTypeof(value runtime.Value, ctype runtime.Type) runtime.Value {
	return runtime.NewBoolean(runtime.Reflect(value) == ctype)
}
