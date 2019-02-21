package common

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func SetInDocument(ctx context.Context, doc drivers.HTMLDocument, path []core.Value, value core.Value) error {
	if path == nil || len(path) == 0 {
		return nil
	}

	segment := path[0]

	if segment.Type() == types.String {
		segment := segment.(values.String)

		switch segment {
		case "url", "URL":
			return doc.SetURL(ctx, values.NewString(value.String()))
		default:
			return SetInNode(ctx, doc, path, value)
		}
	}

	return SetInNode(ctx, doc, path, value)
}

func SetInElement(ctx context.Context, el drivers.HTMLElement, path []core.Value, value core.Value) error {
	if path == nil || len(path) == 0 {
		return nil
	}

	segment := path[0]

	if segment.Type() == types.String {
		segment := segment.(values.String)

		switch segment {
		case "attributes":
			if len(path) > 1 {
				attrName := path[1]

				return el.SetAttribute(ctx, values.NewString(attrName.String()), values.NewString(value.String()))
			}

			err := core.ValidateType(value, types.Object)

			if err != nil {
				return err
			}

			obj := value.(*values.Object)
			obj.ForEach(func(value core.Value, key string) bool {
				err = el.SetAttribute(ctx, values.NewString(key), values.NewString(value.String()))

				return err == nil
			})

			return err
		case "value":
			if len(path) > 1 {
				return core.Error(ErrInvalidPath, PathToString(path[1:]))
			}

			return el.SetValue(ctx, value)
		}
	}

	return SetInNode(ctx, el, path, value)
}

func SetInNode(_ context.Context, _ drivers.HTMLNode, path []core.Value, _ core.Value) error {
	if path == nil || len(path) == 0 {
		return nil
	}

	segment := path[0]
	st := segment.Type()

	if st == types.Int {
		return core.Error(core.ErrInvalidOperation, "children are read-only")
	}

	return core.Error(ErrReadOnly, PathToString(path))
}
