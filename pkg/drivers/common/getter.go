package common

import (
	"context"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func GetInPage(ctx context.Context, path []core.Value, page drivers.HTMLPage) (core.Value, core.PathError) {
	if len(path) == 0 {
		return page, nil
	}

	segmentIdx := 0
	segment := path[segmentIdx]

	if segment.Type() == types.String {
		segment := segment.(values.String)

		switch segment {
		case "response":
			resp, err := page.GetResponse(ctx)

			if err != nil {
				return nil, core.NewPathError(
					errors.Wrap(err, "get response"),
					0,
				)
			}

			out, pathErr := resp.GetIn(ctx, path[segmentIdx+1:])

			if pathErr != nil {
				return values.None, core.NewPathErrorFrom(pathErr, segmentIdx)
			}

			return out, nil
		case "mainFrame", "document":
			out, pathErr := GetInDocument(ctx, path[segmentIdx+1:], page.GetMainFrame())

			if pathErr != nil {
				return values.None, core.NewPathErrorFrom(pathErr, segmentIdx)
			}

			return out, nil
		case "frames":
			if len(path) == 1 {
				out, err := page.GetFrames(ctx)

				if err != nil {
					return nil, core.NewPathError(
						errors.Wrap(err, "get response"),
						segmentIdx,
					)
				}

				return out, nil
			}

			segmentIdx = +1
			idx := path[segmentIdx]

			if !values.IsNumber(idx) {
				return values.None, core.NewPathError(
					core.TypeError(idx.Type(), types.Int, types.Float),
					segmentIdx,
				)
			}

			value, err := page.GetFrame(ctx, values.ToInt(idx))

			if err != nil {
				return values.None, core.NewPathError(err, segmentIdx)
			}

			if len(path) == 2 {
				return value, nil
			}

			frame, err := drivers.ToDocument(value)

			if err != nil {
				return values.None, core.NewPathError(err, segmentIdx)
			}

			out, pathErr := GetInDocument(ctx, path[segmentIdx+1:], frame)

			if err != nil {
				return values.None, core.NewPathErrorFrom(pathErr, segmentIdx)
			}

			return out, nil
		case "url", "URL":
			return page.GetURL(), nil
		case "cookies":
			cookies, err := page.GetCookies(ctx)

			if err != nil {
				return values.None, core.NewPathError(err, segmentIdx)
			}

			if len(path) == 1 {
				return cookies, nil
			}

			out, pathErr := cookies.GetIn(ctx, path[segmentIdx+1:])

			if err != nil {
				return values.None, core.NewPathErrorFrom(pathErr, segmentIdx)
			}

			return out, nil
		case "title":
			return page.GetMainFrame().GetTitle(), nil
		case "isClosed":
			return page.IsClosed(), nil
		default:
			return GetInDocument(ctx, path, page.GetMainFrame())
		}
	}

	return GetInDocument(ctx, path, page.GetMainFrame())
}

func GetInDocument(ctx context.Context, path []core.Value, doc drivers.HTMLDocument) (core.Value, core.PathError) {
	if len(path) == 0 {
		return doc, nil
	}

	var out core.Value
	var err error
	segmentIdx := 0
	segment := path[segmentIdx]

	if segment.Type() == types.String {
		segment := segment.(values.String)

		switch segment {
		case "url", "URL":
			return doc.GetURL(), nil
		case "name":
			return doc.GetName(), nil
		case "title":
			return doc.GetTitle(), nil
		case "parent":
			parent, err := doc.GetParentDocument(ctx)

			if err != nil {
				return values.None, core.NewPathError(err, segmentIdx)
			}

			if parent == nil {
				return values.None, nil
			}

			if len(path) == 1 {
				return parent, nil
			}

			out, pathErr := GetInDocument(ctx, path[segmentIdx+1:], parent)

			if pathErr != nil {
				return values.None, core.NewPathErrorFrom(pathErr, segmentIdx)
			}

			return out, nil
		case "body", "head":
			out, err := doc.QuerySelector(ctx, drivers.NewCSSSelector(segment))

			if err != nil {
				return values.None, core.NewPathError(err, segmentIdx)
			}

			if out == values.None {
				return out, nil
			}

			if len(path) == 1 {
				return out, nil
			}

			el, err := drivers.ToElement(out)

			if err != nil {
				return values.None, core.NewPathError(err, segmentIdx)
			}

			out, pathErr := GetInElement(ctx, path[segmentIdx+1:], el)

			if pathErr != nil {
				return values.None, core.NewPathErrorFrom(pathErr, segmentIdx)
			}

			return out, nil
		case "innerHTML":
			out, err = doc.GetElement().GetInnerHTML(ctx)
		case "innerText":
			out, err = doc.GetElement().GetInnerText(ctx)
		default:
			return GetInNode(ctx, path, doc.GetElement())
		}

		return values.ReturnOrNext(ctx, path, segmentIdx, out, err)
	}

	return GetInNode(ctx, path, doc.GetElement())
}

func GetInElement(ctx context.Context, path []core.Value, el drivers.HTMLElement) (core.Value, core.PathError) {
	if len(path) == 0 {
		return el, nil
	}

	segmentIdx := 0
	segment := path[segmentIdx]

	if segment.Type() == types.String {
		var out core.Value
		var err error

		segment := segment.(values.String)

		switch segment {
		case "innerText":
			out, err = el.GetInnerText(ctx)
		case "innerHTML":
			out, err = el.GetInnerHTML(ctx)
		case "value":
			out, err = el.GetValue(ctx)
		case "attributes":
			if len(path) == 1 {
				out, err = el.GetAttributes(ctx)
			} else {
				// e.g. attributes.href
				segmentIdx++
				attrName := path[segmentIdx]

				out, err = el.GetAttribute(ctx, values.ToString(attrName))
			}
		case "style":
			if len(path) == 1 {
				out, err = el.GetStyles(ctx)
			} else {
				// e.g. style.color
				segmentIdx++
				styleName := path[segmentIdx]

				out, err = el.GetStyle(ctx, values.ToString(styleName))
			}
		case "previousElementSibling":
			out, err = el.GetPreviousElementSibling(ctx)
		case "nextElementSibling":
			out, err = el.GetNextElementSibling(ctx)
		case "parentElement":
			out, err = el.GetParentElement(ctx)
		default:
			return GetInNode(ctx, path, el)
		}

		return values.ReturnOrNext(ctx, path, segmentIdx, out, err)
	}

	return GetInNode(ctx, path, el)
}

func GetInNode(ctx context.Context, path []core.Value, node drivers.HTMLNode) (core.Value, core.PathError) {
	if len(path) == 0 {
		return node, nil
	}

	segmentIdx := 0
	segment := path[segmentIdx]

	var out core.Value
	var err error

	switch segment.Type() {
	case types.Int:
		out, err = node.GetChildNode(ctx, values.ToInt(segment))
	case types.String:
		segment := segment.(values.String)

		switch segment {
		case "nodeType":
			out, err = node.GetNodeType(ctx)
		case "nodeName":
			out, err = node.GetNodeName(ctx)
		case "children":
			if len(path) == 1 {
				out, err = node.GetChildNodes(ctx)
			} else {
				segmentIdx++
				out, err = node.GetChildNode(ctx, values.ToInt(path[segmentIdx]))
			}
		case "length":
			return node.Length(), nil
		default:
			return values.None, nil
		}
	default:
		return values.None, core.NewPathError(
			core.TypeError(segment.Type(), types.Int, types.String),
			segmentIdx,
		)
	}

	return values.ReturnOrNext(ctx, path, segmentIdx, out, err)
}
