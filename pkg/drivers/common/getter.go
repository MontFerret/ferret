package common

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func GetInDocument(ctx context.Context, doc drivers.HTMLDocument, path []core.Value) (core.Value, error) {
	if path == nil || len(path) == 0 {
		return values.None, nil
	}

	var result core.Value = doc
	var err error

	for i, segment := range path {
		if result == values.None || result == nil {
			break
		}

		if segment.Type() == types.String {
			segment := segment.(values.String)

			switch segment {
			case "URL":
				result = doc.GetURL()
			case "body":
				result = doc.QuerySelector("body")
			case "head":
				result = doc.QuerySelector("head")
			default:
				return GetInNode(ctx, doc.DocumentElement(), path[i:])
			}
		} else {
			return GetInNode(ctx, doc.DocumentElement(), path[i:])
		}
	}

	return result, err
}

func GetInElement(ctx context.Context, el drivers.HTMLElement, path []core.Value) (core.Value, error) {
	if path == nil || len(path) == 0 {
		return values.None, nil
	}

	var result core.Value = el
	var err error

	for i, segment := range path {
		if result == values.None || result == nil {
			break
		}

		if segment.Type() == types.String {
			segment := segment.(values.String)

			switch segment {
			case "innerText":
				result = el.InnerText()
			case "innerHTML":
				result = el.InnerHTML()
			case "value":
				result = el.GetValue()
			case "attributes":
				result = el.GetAttributes()
			default:
				return GetInNode(ctx, el, path[i:])
			}
		} else {
			return GetInNode(ctx, el, path[i:])
		}
	}

	return result, err
}

func GetInNode(ctx context.Context, node drivers.HTMLNode, path []core.Value) (core.Value, error) {
	if path == nil || len(path) == 0 {
		return values.None, nil
	}

	var result core.Value = node
	var err error

	for i, segment := range path {
		if result == values.None || result == nil {
			break
		}

		st := segment.Type()

		if st == types.Int {
			rt := result.Type()

			if rt == drivers.HTMLElementType {
				re := result.(drivers.HTMLElement)

				result = re.GetChildNode(segment.(values.Int))
			} else {
				result, err = values.GetIn(ctx, result, path[i:])

				if err != nil {
					return values.None, err
				}
			}
		} else if st == types.String {
			segment := segment.(values.String)

			switch segment {
			case "nodeType":
				result = node.NodeType()
			case "nodeName":
				result = node.NodeName()
			case "children":
				result = node.GetChildNodes()
			case "length":
				result = node.Length()
			default:
				result, err = values.GetIn(ctx, result, path[i:])

				if err != nil {
					return values.None, err
				}
			}
		} else {
			return values.None, core.TypeError(st, types.Int, types.String)
		}
	}

	return result, err
}
