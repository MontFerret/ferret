package operator

// ArrayComparator identifies the element predicate encoded by quantified array opcodes.
type ArrayComparator uint8

const (
	ArrayEqual ArrayComparator = iota
	ArrayNotEqual
	ArrayGreater
	ArrayGreaterOrEqual
	ArrayLess
	ArrayLessOrEqual
	ArrayIn
)

// UnknownArrayComparator represents an invalid quantified-comparison encoding.
const UnknownArrayComparator ArrayComparator = ^ArrayComparator(0)

// ArrayComparatorFromOffset decodes an offset within a quantified opcode group.
func ArrayComparatorFromOffset(offset int) (ArrayComparator, bool) {
	if offset < int(ArrayEqual) || offset > int(ArrayIn) {
		return UnknownArrayComparator, false
	}

	return ArrayComparator(offset), true
}

// ArrayComparatorFor converts a binary comparison operator to its array encoding.
func ArrayComparatorFor(op Binary) (ArrayComparator, bool) {
	switch op {
	case Equal:
		return ArrayEqual, true
	case NotEqual:
		return ArrayNotEqual, true
	case Greater:
		return ArrayGreater, true
	case GreaterOrEqual:
		return ArrayGreaterOrEqual, true
	case Less:
		return ArrayLess, true
	case LessOrEqual:
		return ArrayLessOrEqual, true
	case In:
		return ArrayIn, true
	default:
		return UnknownArrayComparator, false
	}
}

// Binary returns the comparison operator represented by the array encoding.
func (op ArrayComparator) Binary() (Binary, bool) {
	switch op {
	case ArrayEqual:
		return Equal, true
	case ArrayNotEqual:
		return NotEqual, true
	case ArrayGreater:
		return Greater, true
	case ArrayGreaterOrEqual:
		return GreaterOrEqual, true
	case ArrayLess:
		return Less, true
	case ArrayLessOrEqual:
		return LessOrEqual, true
	case ArrayIn:
		return In, true
	default:
		return Unknown, false
	}
}
