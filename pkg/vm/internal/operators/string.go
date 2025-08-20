package operators

import (
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func Like(left, right runtime.Value) (runtime.Boolean, error) {
	if err := runtime.AssertString(left); err != nil {
		// TODO: Return the error? AQL just returns false
		return runtime.False, nil
	}

	if err := runtime.AssertString(right); err != nil {
		// TODO: Return the error? AQL just returns false
		return runtime.False, nil
	}

	r, err := glob.Compile(right.String())

	if err != nil {
		return runtime.False, errors.Wrap(err, "invalid glob pattern")
	}

	result := r.Match(left.String())

	return runtime.NewBoolean(result), nil
}
