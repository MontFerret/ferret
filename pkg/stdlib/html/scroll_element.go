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
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selector - If document is passed, this param must represent an element selector.
// @param {Object} [params] - Scroll params.
// @param {String} [params.behavior="instant"] - Scroll behavior
// @param {String} [params.block="center"] - Scroll vertical alignment.
// @param {String} [params.inline="center"] - Scroll horizontal alignment.
func ScrollInto(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 3)

	if err != nil {
		return values.None, err
	}

	var doc drivers.HTMLDocument
	var el drivers.HTMLElement
	var selector drivers.QuerySelector
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

		selector, err = drivers.ToQuerySelector(args[1])

		if err != nil {
			return values.None, err
		}

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

			selector, err = drivers.ToQuerySelector(args[1])

			if err != nil {
				return values.None, err
			}

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
		if selector.String() != "" {
			return values.True, doc.ScrollBySelector(ctx, selector, opts)
		}

		return values.True, doc.GetElement().ScrollIntoView(ctx, opts)
	}

	if el != nil {
		return values.True, el.ScrollIntoView(ctx, opts)
	}

	return values.None, core.TypeError(
		args[0].Type(),
		drivers.HTMLPageType,
		drivers.HTMLDocumentType,
		drivers.HTMLElementType,
	)
}
