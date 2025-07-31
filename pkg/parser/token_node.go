package parser

import (
	"github.com/antlr4-go/antlr/v4"
)

type TokenNode struct {
	token antlr.Token
	prev  *TokenNode
	next  *TokenNode
}

func (t *TokenNode) Token() antlr.Token {
	return t.token
}

func (t *TokenNode) Prev() *TokenNode {
	return t.prev
}

func (t *TokenNode) Next() *TokenNode {
	return t.next
}

func (t *TokenNode) PrevAt(n int) *TokenNode {
	if n <= 0 {
		return t
	}

	node := t

	for i := 0; i < n && node != nil; i++ {
		node = node.prev
	}

	return node
}

func (t *TokenNode) NextAt(n int) *TokenNode {
	if n <= 0 {
		return t
	}

	node := t

	for i := 0; i < n && node != nil; i++ {
		node = node.next
	}

	return node
}

func (t *TokenNode) String() string {
	return t.token.GetText()
}
