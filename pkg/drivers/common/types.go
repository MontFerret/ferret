package common

import (
	"golang.org/x/net/html"
)

func FromHTMLType(nt html.NodeType) int {
	switch nt {
	case html.DocumentNode:
		return 9
	case html.ElementNode:
		return 1
	case html.TextNode:
		return 3
	case html.CommentNode:
		return 8
	case html.DoctypeNode:
		return 10
	}

	return 0
}

func ToHTMLType(input int) html.NodeType {
	switch input {
	case 1:
		return html.ElementNode
	case 3:
		return html.TextNode
	case 8:
		return html.CommentNode
	case 9:
		return html.DocumentNode
	case 10:
		return html.DoctypeNode
	default:
		return html.ErrorNode
	}
}
