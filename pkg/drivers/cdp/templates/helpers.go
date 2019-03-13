package templates

import (
	"github.com/MontFerret/ferret/pkg/drivers"
)

func WaitEventToEqOperator(when drivers.WaitEvent) string {
	if when == drivers.WaitEventPresence {
		return "=="
	}

	return "!="
}
