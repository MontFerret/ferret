package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"regexp"
)

// FRAMES finds HTML frames by a given property selector.
// Returns an empty array if frames not found.
// @param {HTMLPage} page - HTML page.
// @param {String} property - Property selector.
// @param {String} exp - Regular expression to match property value.
// @return {HTMLDocument[]} - Returns an array of found HTML frames.
func Frames(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 3)

	if err != nil {
		return values.None, err
	}

	page, err := drivers.ToPage(args[0])

	if err != nil {
		return values.None, err
	}

	frames, err := page.GetFrames(ctx)

	if err != nil {
		return values.None, err
	}

	propName := values.ToString(args[1])
	matcher, err := regexp.Compile(values.ToString(args[2]).String())

	if err != nil {
		return values.None, err
	}

	result, _ := frames.Find(func(value core.Value, idx int) bool {
		doc, e := drivers.ToDocument(value)

		if e != nil {
			err = e
			return false
		}

		currentPropValue, e := doc.GetIn(ctx, []core.Value{propName})

		if e != nil {
			err = e

			return false
		}

		return matcher.MatchString(currentPropValue.String())
	})

	return result, err
}
