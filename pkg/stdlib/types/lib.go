package types

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("TO_BOOL", ToBool).
		Add("TO_INT", ToInt).
		Add("TO_FLOAT", ToFloat).
		Add("TO_STRING", ToString).
		Add("TO_DATETIME", ToDateTime).
		Add("TO_ARRAY", ToArray).
		Add("TO_BINARY", ToBinary).
		Add("TO_OBJECT", ToObject).
		Add("IS_NONE", IsNone).
		Add("IS_BOOL", IsBool).
		Add("IS_INT", IsInt).
		Add("IS_FLOAT", IsFloat).
		Add("IS_STRING", IsString).
		Add("IS_DATETIME", IsDateTime).
		Add("IS_LIST", IsList).
		Add("IS_ARRAY", IsArray).
		Add("IS_MAP", IsMap).
		Add("IS_OBJECT", IsObject).
		Add("IS_BINARY", IsBinary).
		Add("IS_NAN", IsNaN)
}

func isTypeof(value runtime.Value, ctype runtime.Type) runtime.Value {
	return runtime.NewBoolean(runtime.IsSameType(runtime.TypeOf(value), ctype))
}
