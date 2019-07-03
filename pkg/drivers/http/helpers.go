package http

import (
	"bytes"
	"golang.org/x/net/html"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
)

func parseXPathNode(nav *htmlquery.NodeNavigator) (core.Value, error) {
	node := nav.Current()

	if node == nil {
		return values.None, nil
	}

	switch nav.NodeType() {
	case xpath.ElementNode:
		return NewHTMLElement(&goquery.Selection{Nodes: []*html.Node{node}})
	case xpath.RootNode:
		url := htmlquery.SelectAttr(node, "url")
		return NewHTMLDocument(goquery.NewDocumentFromNode(node), url, nil)
	default:
		return values.Parse(node.Data), nil
	}
}

func outerHtml(s *goquery.Selection) (string, error) {
	var buf bytes.Buffer

	if len(s.Nodes) > 0 {
		c := s.Nodes[0]

		err := html.Render(&buf, c)

		if err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}
