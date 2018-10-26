package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

var (
	ErrResultSetMismatch = core.Error(core.ErrInvalidArgument, "count of result in result set is less that count of variables")
)

func ValidateDataSet(set DataSet, variables Variables) error {
	if len(variables) > len(set) {
		return ErrResultSetMismatch
	}

	return nil
}
