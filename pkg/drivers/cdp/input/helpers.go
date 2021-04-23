package input

import (
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func randomDuration(delay int) time.Duration {
	return time.Duration(core.Random2(float64(delay)))
}
