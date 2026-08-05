package operator

// Binary identifies a binary operator used by runtime diagnostics.
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

// CannotApply formats the canonical diagnostic for incompatible operands.
func CannotApply(op Binary, leftType, rightType string) string {
	return "operator '" + op.String() + "' cannot be applied to " + leftType + " and " + rightType
}
