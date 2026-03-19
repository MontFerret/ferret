package core

// JoinValueTypes returns the most precise type that safely represents both
// operands for later compiler inference.
func JoinValueTypes(left, right ValueType) ValueType {
	if left == TypeAny || right == TypeAny {
		return TypeAny
	}

	if left == TypeUnknown {
		return right
	}

	if right == TypeUnknown {
		return left
	}

	if left == right {
		return left
	}

	return TypeAny
}
