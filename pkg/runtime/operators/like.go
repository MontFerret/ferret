package operators

import (
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func Like(left, right core.Value) (values.Boolean, error) {
	if err := values.AssertString(left); err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	if err := values.AssertString(right); err != nil {
		// TODO: Return the error? AQL just returns false
		return values.False, nil
	}

	r, err := glob.Compile(right.String())

	if err != nil {
		return values.False, errors.Wrap(err, "invalid glob pattern")
	}

	result := r.Match(left.String())

	return values.NewBoolean(result), nil
}
