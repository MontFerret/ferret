package parser

import "github.com/antlr4-go/antlr/v4"

type TokenHistory struct {
	head *TokenNode
	tail *TokenNode
	size int
	cap  int
}

func NewTokenHistory(cap int) *TokenHistory {
	return &TokenHistory{cap: cap}
}

func (h *TokenHistory) Size() int {
	return h.size
}

func (h *TokenHistory) Add(token antlr.Token) {
	if token == nil {
		return
	}

	// Avoid adding the same token twice in a row (by position, not just text)
	if h.head != nil {
		last := h.head.token
		if last.GetStart() == token.GetStart() &&
			last.GetStop() == token.GetStop() &&
			last.GetTokenType() == token.GetTokenType() {
			return
		}
	}

	node := &TokenNode{token: token}

	if h.head != nil {
		node.next = h.head
		h.head.prev = node
	}

	h.head = node

	if h.tail == nil {
		h.tail = node
	}

	h.size++

	if h.size > h.cap {
		// Remove oldest
		h.tail = h.tail.prev

		if h.tail != nil {
			h.tail.next = nil
		}

		h.size--
	}
}

func (h *TokenHistory) Last() *TokenNode {
	if h.head == nil {
		return nil
	}

	return h.head
}

func (h *TokenHistory) Iterate(yield func(token antlr.Token) bool) {
	curr := h.tail

	for curr != nil {
		if !yield(curr.token) {
			break
		}

		curr = curr.prev
	}
}
