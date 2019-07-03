package eval

import (
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func Param(input core.Value) string {
	switch value := input.(type) {
	case values.String:
		return ParamString(string(value))
	case values.Float:
		return ParamFloat(float64(value))
	case values.Int:
		return ParamInt(int64(value))
	default:
		if input == values.None {
			return "null"
		}

		return value.String()
	}
}

func ParamString(param string) string {
	return "`" + param + "`"
}

func ParamFloat(param float64) string {
	return strconv.FormatFloat(param, 'f', 6, 64)
}

func ParamInt(param int64) string {
	return strconv.FormatInt(param, 64)
}
