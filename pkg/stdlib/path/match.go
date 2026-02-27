package path

import (
	"context"
	"path"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// MATCH reports whether name matches the pattern.
// @param {String} pattern - The pattern.
// @param {String} name - The name.
// @return {Boolean} - True if the name matches the pattern.
func Match(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateType(arg1, runtime.TypeString); err != nil {
		return runtime.False, err
	}

	if err := runtime.ValidateType(arg2, runtime.TypeString); err != nil {
		return runtime.False, err
	}

	pattern := arg1.String()
	name := arg2.String()

	matched, err := path.Match(pattern, name)

	if err != nil {
		return runtime.False, runtime.Error(err, "match")
	}

	return runtime.NewBoolean(matched), nil
}
