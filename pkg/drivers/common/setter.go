package common

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func SetIn(_ context.Context, to drivers.HTMLNode, path []core.Value, value core.Value) error {
	if path == nil || len(path) == 0 {
		return nil
	}

	segment := path[0]
	st := segment.Type()
	tt := to.Type()

	if st == types.Int {
		return core.Error(core.ErrInvalidOperation, "children is read-only")
	} else if st == types.String {
		segment := segment.(values.String)

		switch segment {
		case "nodeType":
			return core.Error(ErrReadOnly, "nodeType")
		case "nodeName":
			return core.Error(ErrReadOnly, "nodeName")
		case "innerText":
			return core.Error(ErrReadOnly, "innerText")
		case "innerHTML":
			return core.Error(ErrReadOnly, "innerHTML")
		case "value":
			if len(path) > 1 {
				return core.Error(ErrInvalidPath, "too long")
			}

			return to.SetValue(value)
		case "attributes":
			if len(path) > 1 {
				attrName := path[1]

				return to.SetAttribute(values.NewString(attrName.String()), values.NewString(value.String()))
			}

			err := core.ValidateType(value, types.Object)

			if err != nil {
				return err
			}

			obj := value.(*values.Object)
			obj.ForEach(func(value core.Value, key string) bool {
				err = to.SetAttribute(values.NewString(key), values.NewString(value.String()))

				return err == nil
			})

			return err
		case "children":
			return core.Error(ErrReadOnly, "children")
		case "length":
			return core.Error(ErrReadOnly, "length")
		case "url":
			if tt == drivers.HTMLDocumentType || tt == drivers.DHTMLDocumentType {
				doc, ok := to.(drivers.HTMLDocument)

				if ok {
					return doc.SetURL(values.NewString(value.String()))
				}
			}

			return core.TypeError(to.Type(), drivers.HTMLDocumentType, drivers.DHTMLDocumentType)
		default:
			return ErrInvalidPath
		}
	}

	return nil
}
