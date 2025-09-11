package operators

import (
	"fmt"

	"github.com/gobwas/glob"

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
		return runtime.False, fmt.Errorf("invalid glob pattern: %w", err)
	}

	result := r.Match(left.String())

	return runtime.NewBoolean(result), nil
}
