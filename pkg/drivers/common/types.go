package common

import (
	"golang.org/x/net/html"
)

func ToHTMLType(nt html.NodeType) int {
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
