package vm

import (
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func relationalComparisonError(op operator.Binary, left, right runtime.Value, err error) error {
	if errors.Is(err, runtime.ErrInvalidOperation) {
		return binaryOperatorError(op, left, right)
	}

	return err
}

func binaryOperatorError(op operator.Binary, left, right runtime.Value) error {
	return runtime.Error(
		runtime.ErrInvalidOperation,
		operator.CannotApply(
			op,
			runtime.TypeName(runtime.TypeOf(left)),
			runtime.TypeName(runtime.TypeOf(right)),
		),
	)
}
