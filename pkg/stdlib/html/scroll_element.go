package html

import (
	"context"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// SCROLL_ELEMENT scrolls an element on.
// @param docOrEl (HTMLDocument|HTMLElement) - Target document or element.
// @param selector (String) - If document is passed, this param must represent an element selector.
// @param options (ScrollOptions) - Scroll options. Optional.
func ScrollInto(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 3)

	if err != nil {
		return values.None, err
	}

	var doc drivers.HTMLDocument
	var el drivers.HTMLElement
	var selector values.String
	var opts drivers.ScrollOptions

	if len(args) == 3 {
		if err = core.ValidateType(args[1], types.String); err != nil {
			return values.None, errors.Wrap(err, "selector")
		}

		if err = core.ValidateType(args[2], types.Object); err != nil {
			return values.None, errors.Wrap(err, "options")
		}

		doc, err = drivers.ToDocument(args[0])

		if err != nil {
			return values.None, errors.Wrap(err, "document")
		}

		selector = values.ToString(args[1])
		o, err := toScrollOptions(args[2])

		if err != nil {
			return values.None, errors.Wrap(err, "options")
		}

		opts = o
	} else if len(args) == 2 {
		if err = core.ValidateType(args[1], types.String, types.Object); err != nil {
			return values.None, err
		}

		if args[1].Type() == types.String {
			doc, err = drivers.ToDocument(args[0])

			if err != nil {
				return values.None, errors.Wrap(err, "document")
			}

			selector = values.ToString(args[1])
		} else {
			el, err = drivers.ToElement(args[0])
			o, err := toScrollOptions(args[1])

			if err != nil {
				return values.None, errors.Wrap(err, "options")
			}

			opts = o
		}
	} else {
		el, err = drivers.ToElement(args[0])

		if err != nil {
			return values.None, errors.Wrap(err, "element")
		}
	}

	if doc != nil {
		if selector != values.EmptyString {
			return values.None, doc.ScrollBySelector(ctx, selector, opts)
		}

		return values.None, doc.GetElement().ScrollIntoView(ctx, opts)
	}

	if el != nil {
		return values.None, el.ScrollIntoView(ctx, opts)
	}

	return values.None, core.TypeError(
		args[0].Type(),
		drivers.HTMLPageType,
		drivers.HTMLDocumentType,
		drivers.HTMLElementType,
	)
}
