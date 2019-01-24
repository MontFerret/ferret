package common

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func SetIn(_ context.Context, to drivers.HTMLNode, path []core.Value, _ core.Value) error {
	if path == nil || len(path) == 0 {
		return nil
	}

	segment := path[0]
	st := segment.Type()

	if st == core.IntType {

	} else if st == core.StringType {
		segment := segment.(values.String)

		switch segment {
		case "nodeType":
			return core.Error(core.ErrInvalidOperation, "nodeType is read-only")
		case "nodeName":
			return core.Error(core.ErrInvalidOperation, "nodeName is read-only")
		case "innerText":
			return core.Error(core.ErrInvalidOperation, "innerText is read-only")
		case "innerHTML":
			return core.Error(core.ErrInvalidOperation, "innerHTML is read-only")
		case "value":
			return core.Error(core.ErrInvalidOperation, "value is read-only")
		case "attributes":
			return core.Error(core.ErrInvalidOperation, "attributes is read-only")
		case "children":
			return core.Error(core.ErrInvalidOperation, "children is read-only")
		case "length":
			return core.Error(core.ErrInvalidOperation, "length is read-only")
		case "url":
			//rt := to.Type()
			//
			//if rt == drivers.HTMLDocumentType || rt == drivers.DHTMLDocumentType {
			//	doc, ok := to.(drivers.HTMLDocument)
			//
			//	if ok {
			//		result = doc.URL()
			//	}
			//}
			return core.Error(core.ErrInvalidOperation, "url is read-only")
		default:
			return nil
		}
	}

	return nil
}
