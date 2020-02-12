package common

import (
	"context"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func GetInPage(ctx context.Context, page drivers.HTMLPage, path []core.Value) (core.Value, error) {
	if len(path) == 0 {
		return page, nil
	}

	segment := path[0]

	if segment.Type() == types.String {
		segment := segment.(values.String)

		switch segment {
		case "response":
			resp, err := page.GetResponse(ctx)

			if err != nil {
				return nil, errors.Wrap(err, "get response")
			}

			return resp.GetIn(ctx, path[1:])
		case "mainFrame", "document":
			return GetInDocument(ctx, page.GetMainFrame(), path[1:])
		case "frames":
			if len(path) == 1 {
				return page.GetFrames(ctx)
			}

			idx := path[1]

			if !values.IsNumber(idx) {
				return values.None, core.TypeError(idx.Type(), types.Int, types.Float)
			}

			value, err := page.GetFrame(ctx, values.ToInt(idx))

			if err != nil {
				return values.None, err
			}

			if len(path) == 2 {
				return value, nil
			}

			frame, err := drivers.ToDocument(value)

			if err != nil {
				return values.None, err
			}

			return GetInDocument(ctx, frame, path[2:])
		case "url", "URL":
			return page.GetMainFrame().GetURL(), nil
		case "cookies":
			cookies, err := page.GetCookies(ctx)

			if err != nil {
				return values.None, err
			}

			if len(path) == 1 {
				return cookies, nil
			}

			return cookies.GetIn(ctx, path[1:])
		case "isClosed":
			return page.IsClosed(), nil
		case "title":
			return page.GetMainFrame().GetTitle(), nil
		default:
			return GetInDocument(ctx, page.GetMainFrame(), path)
		}
	}

	return GetInDocument(ctx, page.GetMainFrame(), path)
}

func GetInDocument(ctx context.Context, doc drivers.HTMLDocument, path []core.Value) (core.Value, error) {
	if len(path) == 0 {
		return doc, nil
	}

	segment := path[0]

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
				return values.None, err
			}

			if parent == nil {
				return values.None, nil
			}

			if len(path) == 1 {
				return parent, nil
			}

			return GetInDocument(ctx, parent, path[1:])
		case "body", "head":
			out, err := doc.QuerySelector(ctx, segment)

			if err != nil {
				return values.None, err
			}

			if out == values.None {
				return out, nil
			}

			if len(path) == 1 {
				return out, nil
			}

			el, err := drivers.ToElement(out)

			if err != nil {
				return values.None, err
			}

			return GetInElement(ctx, el, path[1:])
		case "innerHTML":
			return doc.GetElement().GetInnerHTML(ctx)
		case "innerText":
			return doc.GetElement().GetInnerText(ctx)
		default:
			return GetInNode(ctx, doc.GetElement(), path)
		}
	}

	return GetInNode(ctx, doc.GetElement(), path)
}

func GetInElement(ctx context.Context, el drivers.HTMLElement, path []core.Value) (core.Value, error) {
	if len(path) == 0 {
		return el, nil
	}

	segment := path[0]

	if segment.Type() == types.String {
		segment := segment.(values.String)

		switch segment {
		case "innerText":
			return el.GetInnerText(ctx)
		case "innerHTML":
			return el.GetInnerHTML(ctx)
		case "value":
			return el.GetValue(ctx)
		case "attributes":
			attrs, err := el.GetAttributes(ctx)

			if err != nil {
				return values.None, err
			}

			if len(path) == 1 {
				return attrs, nil
			}

			return values.GetIn(ctx, attrs, path[1:])
		case "style":
			styles, err := el.GetStyles(ctx)

			if err != nil {
				return values.None, err
			}

			if len(path) == 1 {
				return styles, nil
			}

			return values.GetIn(ctx, styles, path[1:])
		default:
			return GetInNode(ctx, el, path)
		}
	}

	return GetInNode(ctx, el, path)
}

func GetInNode(ctx context.Context, node drivers.HTMLNode, path []core.Value) (core.Value, error) {
	if len(path) == 0 {
		return node, nil
	}

	nt := node.Type()
	segment := path[0]
	st := segment.Type()

	switch st {
	case types.Int:
		if nt == drivers.HTMLElementType || nt == drivers.HTMLDocumentType {
			re := node.(drivers.HTMLNode)

			return re.GetChildNode(ctx, values.ToInt(segment))
		}

		return values.GetIn(ctx, node, path[1:])
	case types.String:
		segment := segment.(values.String)

		switch segment {
		case "isDetached":
			return node.IsDetached(), nil
		case "nodeType":
			return node.GetNodeType(), nil
		case "nodeName":
			return node.GetNodeName(), nil
		case "children":
			children, err := node.GetChildNodes(ctx)

			if err != nil {
				return values.None, err
			}

			if len(path) == 1 {
				return children, nil
			}

			return values.GetIn(ctx, children, path[1:])
		case "length":
			return node.Length(), nil
		default:
			return values.None, nil
		}
	default:
		return values.None, core.TypeError(st, types.Int, types.String)
	}
}
