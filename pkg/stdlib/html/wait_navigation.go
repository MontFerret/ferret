package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

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

	page, err := drivers.ToPage(args[0])

	if err != nil {
		return values.None, err
	}

	var timeout values.Int

	sub := events.Subscription{
		EventName: drivers.NavigationPageEvent,
		Options:   values.NewObjectWith(),
	}

	if len(args) > 1 {
		if err := core.ValidateType(args[1], types.Object, types.Int); err != nil {
			return values.None, err
		}

		if types.Int == args[1].Type() {
			timeout = values.ToInt(args[1])
		} else {
			obj := values.ToObject(ctx, args[1])

			if obj.Has("timeout") {
				timeout = values.ToInt(obj.MustGet("timeout"))
				obj.Remove("timeout")
			}

			sub.Options = obj
		}
	}

	if timeout <= 0 {
		timeout = drivers.DefaultWaitTimeout
	}

	ctx, fn := waitTimeout(ctx, timeout)
	defer fn()

	ch, err := page.Subscribe(ctx, sub)

	if err != nil {
		return values.None, err
	}

	evt := <-ch

	return evt.Data, evt.Err
}
