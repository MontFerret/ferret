package source

import (
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

var (
	ErrResultSetMismatch = core.Error(core.ErrInvalidArgument, "count of values in result set is less that count of variables")
)

func ValidateResultSet(variables Variables, set collections.ResultSet) error {
	if len(variables) > len(set) {
		return ErrResultSetMismatch
	}

	return nil
}
