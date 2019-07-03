package input

import (
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func randomDuration(delay int) time.Duration {
	max, min := core.NumberBoundaries(float64(delay))
	value := core.Random(max, min)

	return time.Duration(int64(value))
}
