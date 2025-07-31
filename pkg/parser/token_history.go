package parser

import "github.com/antlr4-go/antlr/v4"

type (
	TokenNode struct {
		Token antlr.Token
		Prev  *TokenNode
		Next  *TokenNode
	}

	TokenHistory struct {
		head *TokenNode
		tail *TokenNode
		size int
		cap  int
	}
)

func NewTokenHistory(cap int) *TokenHistory {
	return &TokenHistory{cap: cap}
}

func (h *TokenHistory) Add(token antlr.Token) {
	if token == nil {
		return
	}

	node := &TokenNode{Token: token}

	if h.head != nil {
		node.Next = h.head
		h.head.Prev = node
	}

	h.head = node

	if h.tail == nil {
		h.tail = node
	}

	h.size++

	if h.size > h.cap {
		// Remove oldest
		h.tail = h.tail.Prev

		if h.tail != nil {
			h.tail.Next = nil
		}

		h.size--
	}
}

func (h *TokenHistory) LastN(n int) []antlr.Token {
	result := make([]antlr.Token, 0, n)
	curr := h.head

	for curr != nil && n > 0 {
		result = append(result, curr.Token)
		curr = curr.Next
		n--
	}

	return result
}

func (h *TokenHistory) Last() antlr.Token {
	if h.tail == nil {
		return nil
	}

	return h.tail.Token
}
