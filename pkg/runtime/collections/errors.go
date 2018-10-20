package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

var (
	ErrExhausted         = core.Error(core.ErrInvalidOperation, "iterator has been exhausted")
	ErrResultSetMismatch = core.Error(core.ErrInvalidArgument, "count of values in result set is less that count of variables")
)

func ValidateDataSet(set DataSet, variables Variables) error {
	if len(variables) > len(set) {
		return ErrResultSetMismatch
	}

	return nil
}
