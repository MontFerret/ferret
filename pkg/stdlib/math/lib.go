package math

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"ABS":     Abs,
		"ACOS":    Acos,
		"ASIN":    Asin,
		"ATAN":    Atan,
		"ATAN2":   Atan2,
		"AVERAGE": Average,
		"CEIL":    Ceil,
		"COS":     Cos,
		"DEGREES": Degrees,
		"EXP":     Exp,
		"EXP2":    Exp2,
		"FLOOR":   Floor,
		"LOG":     Log,
		"LOG2":    Log2,
		"LOG10":   Log10,
		"MAX":     Max,
	}
}

func toFloat(arg core.Value) float64 {
	if arg.Type() == core.IntType {
		return float64(arg.(values.Int))
	}

	return float64(arg.(values.Float))
}
