package runtime

import "context"

type (
	// Addable lets a host value implement binary addition in either operand position.
	//
	// For left + right, the runtime calls left.Add(ctx, right) first. If that
	// returns ErrUnsupportedOperands, and right also implements Addable, the
	// runtime calls right.RightAdd(ctx, left). Any other error stops negotiation.
	// The execution context is forwarded unchanged; implementations that may
	// block are responsible for observing cancellation while they retain control.
	Addable interface {
		// Add handles the receiver as the left operand of addition.
		Add(ctx context.Context, right Value) (Value, error)
		// RightAdd handles the receiver as the right operand of addition.
		RightAdd(ctx context.Context, left Value) (Value, error)
	}

	// Subtractable lets a host value implement binary subtraction in either operand position.
	// RightSubtract receives the original left operand; operands are never reversed implicitly.
	// Negotiation, error handling, and context ownership follow Addable.
	Subtractable interface {
		// Subtract handles the receiver as the left operand of subtraction.
		Subtract(ctx context.Context, right Value) (Value, error)
		// RightSubtract handles the receiver as the right operand of subtraction.
		RightSubtract(ctx context.Context, left Value) (Value, error)
	}

	// Multipliable lets a host value implement binary multiplication in either operand position.
	// Negotiation, error handling, and context ownership follow Addable.
	Multipliable interface {
		// Multiply handles the receiver as the left operand of multiplication.
		Multiply(ctx context.Context, right Value) (Value, error)
		// RightMultiply handles the receiver as the right operand of multiplication.
		RightMultiply(ctx context.Context, left Value) (Value, error)
	}

	// Dividable lets a host value implement binary division in either operand position.
	// RightDivide receives the original left operand; operands are never reversed implicitly.
	// Negotiation, error handling, and context ownership follow Addable.
	Dividable interface {
		// Divide handles the receiver as the left operand of division.
		Divide(ctx context.Context, right Value) (Value, error)
		// RightDivide handles the receiver as the right operand of division.
		RightDivide(ctx context.Context, left Value) (Value, error)
	}

	// Modulable lets a host value implement binary modulus in either operand position.
	// RightMod receives the original left operand; operands are never reversed implicitly.
	// Negotiation, error handling, and context ownership follow Addable.
	Modulable interface {
		// Mod handles the receiver as the left operand of modulus.
		Mod(ctx context.Context, right Value) (Value, error)
		// RightMod handles the receiver as the right operand of modulus.
		RightMod(ctx context.Context, left Value) (Value, error)
	}

	// Additive is a convenience contract for values supporting addition and subtraction.
	// Runtime operator dispatch checks Addable or Subtractable directly, not Additive.
	Additive interface {
		Addable
		Subtractable
	}

	// Multiplicative is a convenience contract for values supporting multiplication
	// and division. Runtime operator dispatch checks the primitive capability.
	Multiplicative interface {
		Multipliable
		Dividable
	}

	// Arithmetic is a convenience contract containing every binary arithmetic capability.
	// It exists for compile-time assertions and has no special runtime dispatch behavior.
	Arithmetic interface {
		Additive
		Multiplicative
		Modulable
	}
)
