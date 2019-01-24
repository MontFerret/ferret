package common

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func GetIn(ctx context.Context, el drivers.HTMLNode, path []core.Value) (core.Value, error) {
	if path == nil || len(path) == 0 {
		return values.None, nil
	}

	var result core.Value = el
	var err error

	for i, segment := range path {
		if result == values.None || result == nil {
			break
		}

		st := segment.Type()

		if st == core.IntType {
			rt := result.Type()

			if rt == drivers.HTMLElementType || rt == drivers.DHTMLElementType {
				re := result.(drivers.HTMLNode)

				result = re.GetChildNode(segment.(values.Int))
			} else {
				result, err = values.GetIn(ctx, result, path[i:])

				if err != nil {
					return values.None, err
				}
			}
		} else if st == core.StringType {
			segment := segment.(values.String)

			switch segment {
			case "nodeType":
				result = el.NodeType()
			case "nodeName":
				result = el.NodeName()
			case "innerText":
				result = el.InnerText()
			case "innerHTML":
				result = el.InnerHTML()
			case "value":
				result = el.Value()
			case "attributes":
				result = el.GetAttributes()
			case "children":
				result = el.GetChildNodes()
			case "length":
				result = el.Length()
			case "url":
				rt := result.Type()

				if rt == drivers.HTMLDocumentType || rt == drivers.DHTMLDocumentType {
					doc, ok := result.(drivers.HTMLDocument)

					if ok {
						result = doc.URL()
					}
				}
			default:
				result, err = values.GetIn(ctx, result, path[i:])

				if err != nil {
					return values.None, err
				}
			}
		} else {
			return values.None, core.TypeError(st, core.IntType, core.StringType)
		}
	}

	return result, err
}
