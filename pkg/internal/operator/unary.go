package operator

// Unary identifies a unary operator used by runtime diagnostics.
type Unary uint8

const (
	UnknownUnary Unary = iota
	Not
	Positive
	Negative
	Increment
	Decrement
)

// ParseUnary converts an exact FQL unary operator symbol into its diagnostic form.
func ParseUnary(input string) (Unary, bool) {
	switch input {
	case "!":
		return Not, true
	case "+":
		return Positive, true
	case "-":
		return Negative, true
	case "++":
		return Increment, true
	case "--":
		return Decrement, true
	default:
		return UnknownUnary, false
	}
}

// String returns the canonical FQL symbol for the operator.
func (op Unary) String() string {
	switch op {
	case Not:
		return "!"
	case Positive:
		return "+"
	case Negative:
		return "-"
	case Increment:
		return "++"
	case Decrement:
		return "--"
	default:
		return "?"
	}
}

// CannotApplyUnary formats the canonical diagnostic for an incompatible operand.
func CannotApplyUnary(op Unary, operandType string) string {
	return "operator '" + op.String() + "' cannot be applied to " + operandType
}
