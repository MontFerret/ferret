package operators

import (
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func Like(left, right core.Value) (core.Value, error) {
	err := core.ValidateType(right, types.String)

	if err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	err = core.ValidateType(left, types.String)

	if err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	r, err := glob.Compile(right.String())

	if err != nil {
		return nil, errors.Wrap(err, "invalid glob pattern")
	}

	result := r.Match(left.String())

	return values.NewBoolean(result), nil
}
