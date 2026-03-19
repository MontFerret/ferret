package diagnostics

func hasPrevToken(node *TokenNode, tokenText string, steps int) bool {
	return findPrevToken(node, tokenText, steps) != nil
}

func findPrevToken(node *TokenNode, tokenText string, steps int) *TokenNode {
	current := node
	for i := 0; i < steps && current != nil; i++ {
		if is(current, tokenText) {
			return current
		}
		current = current.Prev()
	}

	return nil
}

func findNextToken(node *TokenNode, tokenText string, steps int) *TokenNode {
	current := node
	for i := 0; i < steps && current != nil; i++ {
		if is(current, tokenText) {
			return current
		}
		current = current.Next()
	}

	return nil
}
