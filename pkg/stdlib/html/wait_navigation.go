package html

import (
	"context"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type WaitNavigationParams struct {
	TargetURL values.String
	Timeout   values.Int
	Frame     drivers.HTMLDocument
}

// WAIT_NAVIGATION waits for a given page to navigate to a new url.
// Stops the execution until the navigation ends or operation times out.
// @param {HTMLPage} page - Target page.
// @param {Int} [timeout=5000] - Navigation timeout.
// @param {Object} [params=None] - Navigation parameters.
// @param {Int} [params.timeout=5000] - Navigation timeout.
// @param {String} [params.target] - Navigation target url.
// @param {HTMLDocument} [params.frame] - Navigation frame.
func WaitNavigation(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	doc, err := drivers.ToPage(args[0])

	if err != nil {
		return values.None, err
	}

	var params WaitNavigationParams

	if len(args) > 1 {
		p, err := parseWaitNavigationParams(args[1])

		if err != nil {
			return values.None, err
		}

		params = p
	} else {
		params = defaultWaitNavigationParams()
	}

	ctx, fn := waitTimeout(ctx, params.Timeout)
	defer fn()

	if params.Frame == nil {
		return values.True, doc.WaitForNavigation(ctx, params.TargetURL)
	}

	return values.True, doc.WaitForFrameNavigation(ctx, params.Frame, params.TargetURL)
}

func parseWaitNavigationParams(arg core.Value) (WaitNavigationParams, error) {
	params := defaultWaitNavigationParams()
	err := core.ValidateType(arg, types.Int, types.Object)

	if err != nil {
		return params, err
	}

	if arg.Type() == types.Int {
		params.Timeout = arg.(values.Int)
	} else {
		obj := arg.(*values.Object)

		if v, exists := obj.Get("timeout"); exists {
			err := core.ValidateType(v, types.Int)

			if err != nil {
				return params, errors.Wrap(err, "navigation parameters: timeout")
			}

			params.Timeout = v.(values.Int)
		}

		if v, exists := obj.Get("target"); exists {
			err := core.ValidateType(v, types.String)

			if err != nil {
				return params, errors.Wrap(err, "navigation parameters: url")
			}

			params.TargetURL = v.(values.String)
		}

		if v, exists := obj.Get("frame"); exists {
			doc, err := drivers.ToDocument(v)

			if err != nil {
				return params, errors.Wrap(err, "navigation parameters: frame")
			}

			params.Frame = doc
		}
	}

	return params, nil
}

func defaultWaitNavigationParams() WaitNavigationParams {
	return WaitNavigationParams{
		TargetURL: "",
		Timeout:   values.NewInt(drivers.DefaultWaitTimeout),
	}
}
