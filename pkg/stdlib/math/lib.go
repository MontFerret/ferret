package math

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

const (
	RadToDeg  = 180 / math.Pi
	DegToRad  = math.Pi / 180
	RadToGrad = 200 / math.Pi
	GradToDeg = math.Pi / 200
)

func NewLib() map[string]core.Function {
	return map[string]core.Function{
		"ABS":        Abs,
		"ACOS":       Acos,
		"ASIN":       Asin,
		"ATAN":       Atan,
		"ATAN2":      Atan2,
		"AVERAGE":    Average,
		"CEIL":       Ceil,
		"COS":        Cos,
		"DEGREES":    Degrees,
		"EXP":        Exp,
		"EXP2":       Exp2,
		"FLOOR":      Floor,
		"LOG":        Log,
		"LOG2":       Log2,
		"LOG10":      Log10,
		"MAX":        Max,
		"MEDIAN":     Median,
		"MIN":        Min,
		"PERCENTILE": Percentile,
		"PI":         Pi,
		"POW":        Pow,
		"RADIANS":    Radians,
		"RAND":       Rand,
		"RANGE":      Range,
		"ROUND":      Round,
		"SIN":        Sin,
		"SQRT":       Sqrt,
	}
}

func toFloat(arg core.Value) float64 {
	if arg.Type() == core.IntType {
		return float64(arg.(values.Int))
	}

	return float64(arg.(values.Float))
}
