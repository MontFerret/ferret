package types

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A1().
		Add("to_bool", ToBool).
		Add("to_int", ToInt).
		Add("to_float", ToFloat).
		Add("to_duration", ToDuration).
		Add("to_string", ToString).
		Add("to_array", ToArray).
		Add("to_binary", ToBinary).
		Add("to_number", ToNumber).
		Add("to_object", ToObject).
		Add("is_none", IsNone).
		Add("is_bool", IsBool).
		Add("is_int", IsInt).
		Add("is_float", IsFloat).
		Add("is_duration", IsDuration).
		Add("is_string", IsString).
		Add("is_datetime", IsDateTime).
		Add("is_list", IsList).
		Add("is_array", IsArray).
		Add("is_map", IsMap).
		Add("is_object", IsObject).
		Add("is_binary", IsBinary).
		Add("is_nan", IsNaN).
		Add("to_datetime", ToDateTime)

	ns.Function().A2().Add("to_datetime", toDateTime2)
}

func isTypeof(value runtime.Value, ctype runtime.Type) runtime.Value {
	return runtime.NewBoolean(runtime.IsSameType(runtime.TypeOf(value), ctype))
}
