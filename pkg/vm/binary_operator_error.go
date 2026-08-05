package vm

import (
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func relationalComparisonError(operator string, left, right runtime.Value, err error) error {
	if errors.Is(err, runtime.ErrInvalidOperation) {
		return binaryOperatorError(operator, left, right)
	}

	return err
}

func binaryOperatorError(operator string, left, right runtime.Value) error {
	return runtime.Errorf(
		runtime.ErrInvalidOperation,
		"operator '%s' cannot be applied to %s and %s",
		operator,
		runtime.TypeName(runtime.TypeOf(left)),
		runtime.TypeName(runtime.TypeOf(right)),
	)
}
