package common

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func SetInPage(ctx context.Context, path []core.Value, page drivers.HTMLPage, value core.Value) core.PathError {
	if len(path) == 0 {
		return nil
	}

	return SetInDocument(ctx, path, page.GetMainFrame(), value)
}

func SetInDocument(ctx context.Context, path []core.Value, doc drivers.HTMLDocument, value core.Value) core.PathError {
	if len(path) == 0 {
		return nil
	}

	return SetInNode(ctx, path, doc, value)
}

func SetInElement(ctx context.Context, path []core.Value, el drivers.HTMLElement, value core.Value) core.PathError {
	if len(path) == 0 {
		return nil
	}

	segmentIdx := 0
	segment := path[segmentIdx]

	if segment.Type() == types.String {
		segment := segment.(values.String)

		switch segment {
		case "attributes":
			if len(path) > 1 {
				attrName := path[1]
				err := el.SetAttribute(ctx, values.NewString(attrName.String()), values.NewString(value.String()))

				if err != nil {
					return core.NewPathError(err, segmentIdx)
				}

				return nil
			}

			err := core.ValidateType(value, types.Object)

			if err != nil {
				return core.NewPathError(err, segmentIdx)
			}

			curr, err := el.GetAttributes(ctx)

			if err != nil {
				return core.NewPathError(err, segmentIdx)
			}

			// remove all previous attributes
			err = el.RemoveAttribute(ctx, curr.Keys()...)

			if err != nil {
				return core.NewPathError(err, segmentIdx)
			}

			obj := value.(*values.Object)
			obj.ForEach(func(value core.Value, key string) bool {
				err = el.SetAttribute(ctx, values.NewString(key), values.NewString(value.String()))

				return err == nil
			})

			if err != nil {
				return core.NewPathError(err, segmentIdx)
			}

			return nil
		case "style":
			if len(path) > 1 {
				attrName := path[1]

				err := el.SetStyle(ctx, values.NewString(attrName.String()), values.NewString(value.String()))

				if err != nil {
					return core.NewPathError(err, segmentIdx)
				}

				return nil
			}

			err := core.ValidateType(value, types.Object)

			if err != nil {
				return core.NewPathError(err, segmentIdx)
			}

			styles, err := el.GetStyles(ctx)

			if err != nil {
				return core.NewPathError(err, segmentIdx)
			}

			err = el.RemoveStyle(ctx, styles.Keys()...)

			obj := value.(*values.Object)
			obj.ForEach(func(value core.Value, key string) bool {
				err = el.SetStyle(ctx, values.NewString(key), values.NewString(value.String()))

				return err == nil
			})

			if err != nil {
				return core.NewPathError(err, segmentIdx)
			}

			return nil
		case "value":
			if len(path) > 1 {
				return core.NewPathError(ErrInvalidPath, segmentIdx+1)
			}

			err := el.SetValue(ctx, value)

			if err != nil {
				return core.NewPathError(err, segmentIdx)
			}

			return nil
		}
	}

	return SetInNode(ctx, path, el, value)
}

func SetInNode(_ context.Context, path []core.Value, _ drivers.HTMLNode, _ core.Value) core.PathError {
	if len(path) == 0 {
		return nil
	}

	return core.NewPathError(ErrReadOnly, 0)
}
