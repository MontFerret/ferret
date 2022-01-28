package http

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func EvalXPathToNode(selection *goquery.Selection, expression string) (drivers.HTMLNode, error) {
	node := htmlquery.FindOne(fromSelectionToNode(selection), expression)

	if node == nil {
		return nil, nil
	}

	return parseXPathNode(node)
}

func EvalXPathToElement(selection *goquery.Selection, expression string) (drivers.HTMLElement, error) {
	node, err := EvalXPathToNode(selection, expression)

	if err != nil {
		return nil, err
	}

	if node == nil {
		return nil, nil
	}

	return drivers.ToElement(node)
}

func EvalXPathToNodes(selection *goquery.Selection, expression string) (*values.Array, error) {
	return EvalXPathToNodesWith(selection, expression, func(node *html.Node) (core.Value, error) {
		return parseXPathNode(node)
	})
}

func EvalXPathToNodesWith(selection *goquery.Selection, expression string, mapper func(node *html.Node) (core.Value, error)) (*values.Array, error) {
	out, err := evalXPathToInternal(selection, expression)

	if err != nil {
		return nil, err
	}

	switch res := out.(type) {
	case *xpath.NodeIterator:
		items := values.NewArray(10)

		for res.MoveNext() {
			item, err := mapper(res.Current().(*htmlquery.NodeNavigator).Current())

			if err != nil {
				return nil, err
			}

			if item != nil {
				items.Push(item)
			}
		}

		return items, nil
	default:
		return values.EmptyArray(), nil
	}
}

func EvalXPathTo(selection *goquery.Selection, expression string) (core.Value, error) {
	out, err := evalXPathToInternal(selection, expression)

	if err != nil {
		return nil, err
	}

	switch res := out.(type) {
	case *xpath.NodeIterator:
		items := values.NewArray(10)

		for res.MoveNext() {
			var item core.Value

			node := res.Current()

			switch node.NodeType() {
			case xpath.TextNode:
				item = values.NewString(node.Value())
			case xpath.AttributeNode:
				item = values.NewString(node.Value())
			default:
				i, err := parseXPathNode(node.(*htmlquery.NodeNavigator).Current())

				if err != nil {
					return nil, err
				}

				item = i
			}

			if item != nil {
				items.Push(item)
			}
		}

		return items, nil
	default:
		return values.Parse(res), nil
	}
}

func evalXPathToInternal(selection *goquery.Selection, expression string) (interface{}, error) {
	exp, err := xpath.Compile(expression)

	if err != nil {
		return nil, err
	}

	return exp.Evaluate(htmlquery.CreateXPathNavigator(fromSelectionToNode(selection))), nil
}

func parseXPathNode(node *html.Node) (drivers.HTMLNode, error) {
	if node == nil {
		return nil, nil
	}

	switch node.Type {
	case html.DocumentNode:
		url := htmlquery.SelectAttr(node, "url")
		return NewHTMLDocument(goquery.NewDocumentFromNode(node), url, nil)
	case html.ElementNode:
		return NewHTMLElement(&goquery.Selection{Nodes: []*html.Node{node}})
	default:
		return nil, nil
	}
}
