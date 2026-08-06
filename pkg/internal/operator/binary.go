package operator

// Binary identifies a binary operator used by runtime dispatch and diagnostics.
type Binary uint8

const (
	Unknown Binary = iota
	Add
	Subtract
	Multiply
	Divide
	Modulus
	Less
	LessOrEqual
	Greater
	GreaterOrEqual
	Equal
	NotEqual
	In
)

// ParseBinary converts an exact FQL operator symbol into its diagnostic form.
func ParseBinary(input string) (Binary, bool) {
	switch input {
	case "+":
		return Add, true
	case "-":
		return Subtract, true
	case "*":
		return Multiply, true
	case "/":
		return Divide, true
	case "%":
		return Modulus, true
	case "<":
		return Less, true
	case "<=":
		return LessOrEqual, true
	case ">":
		return Greater, true
	case ">=":
		return GreaterOrEqual, true
	case "==":
		return Equal, true
	case "!=":
		return NotEqual, true
	case "IN":
		return In, true
	default:
		return Unknown, false
	}
}

// String returns the canonical FQL symbol for the operator.
func (op Binary) String() string {
	switch op {
	case Add:
		return "+"
	case Subtract:
		return "-"
	case Multiply:
		return "*"
	case Divide:
		return "/"
	case Modulus:
		return "%"
	case Less:
		return "<"
	case LessOrEqual:
		return "<="
	case Greater:
		return ">"
	case GreaterOrEqual:
		return ">="
	case Equal:
		return "=="
	case NotEqual:
		return "!="
	case In:
		return "IN"
	default:
		return "?"
	}
}

// IsRelational reports whether the operator requires relational comparison.
func (op Binary) IsRelational() bool {
	switch op {
	case Less, LessOrEqual, Greater, GreaterOrEqual:
		return true
	default:
		return false
	}
}

// IsEquality reports whether the operator performs equality comparison.
func (op Binary) IsEquality() bool {
	return op == Equal || op == NotEqual
}

// IsMembership reports whether the operator performs membership comparison.
func (op Binary) IsMembership() bool {
	return op == In
}

// IsComparison reports whether the operator produces a comparison predicate.
func (op Binary) IsComparison() bool {
	return op.IsEquality() || op.IsRelational() || op.IsMembership()
}

// CannotApply formats the canonical diagnostic for incompatible operands.
func CannotApply(op Binary, leftType, rightType string) string {
	return "operator '" + op.String() + "' cannot be applied to " + leftType + " and " + rightType
}
