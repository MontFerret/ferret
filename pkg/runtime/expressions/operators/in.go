package operators

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type InOperator struct {
	*baseOperator
	all bool
	not bool
}

func NewInOperator(
	src core.SourceMap,
	left core.Expression,
	right core.Expression,
	all bool,
	not bool,
) (*InOperator, error) {
	if left == nil {
		return nil, core.Error(core.ErrMissedArgument, "left expression")
	}

	if right == nil {
		return nil, core.Error(core.ErrMissedArgument, "right expression")
	}

	base := &baseOperator{src, left, right}

	return &InOperator{
		base,
		all,
		not,
	}, nil
}

func (operator *InOperator) Exec(ctx context.Context, scope *core.Scope) (core.Value, error) {
	left, err := operator.left.Exec(ctx, scope)

	if err != nil {
		return values.False, core.SourceError(operator.src, err)
	}

	right, err := operator.right.Exec(ctx, scope)

	if err != nil {
		return values.False, core.SourceError(operator.src, err)
	}

	err = core.ValidateType(right, core.ArrayType)

	if err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	var found bool
	rightArr := right.(*values.Array)

	// it means that left value must be an array too
	if operator.all {
		err = core.ValidateType(left, core.ArrayType)

		if err != nil {
			// TODO: Return the error? AQL just returns false
			return values.False, nil
		}

		hashTable, err := collections.ToHashTable(collections.NewArrayIterator(rightArr))

		if err != nil {
			return values.False, err
		}

		leftArr := left.(*values.Array)

		leftArr.ForEach(func(value core.Value, _ int) bool {
			h := value.Hash()

			_, exists := hashTable[h]

			if !exists {
				found = false
				return false
			}

			found = true
			return true
		})
	} else {
		found = rightArr.IndexOf(left) > -1
	}

	if operator.not {
		return values.NewBoolean(found == false), nil
	}

	return values.NewBoolean(found == true), nil
}
