package core

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestConstantPoolVerifiesEqualityWithinHashBucket(t *testing.T) {
	t.Parallel()

	pool := NewConstantPool()
	first := runtime.NewDuration(time.Second)
	second := runtime.NewDuration(2 * time.Second)
	firstOperand := pool.Add(first)

	// Simulate a hash collision without coupling the test to a particular hash implementation.
	pool.index[second.Hash()] = constantBucket{first: firstOperand.Constant()}
	secondOperand := pool.Add(second)

	if firstOperand == secondOperand {
		t.Fatal("distinct constants sharing a hash bucket were deduplicated")
	}
	if got := pool.Get(secondOperand); got != second {
		t.Fatalf("second constant = %v, want %v", got, second)
	}
	if duplicate := pool.Add(second); duplicate != secondOperand {
		t.Fatalf("equal constant in collision bucket = %v, want %v", duplicate, secondOperand)
	}
}

func TestConstantPoolPreservesScalarTypeWithinHashBucket(t *testing.T) {
	t.Parallel()

	pool := NewConstantPool()
	integer := runtime.NewInt(1)
	floating := runtime.NewFloat(1)
	integerOperand := pool.Add(integer)

	// Int and Float compare equal numerically, but constants retain their runtime type.
	pool.index[floating.Hash()] = constantBucket{first: integerOperand.Constant()}
	floatOperand := pool.Add(floating)
	if floatOperand == integerOperand {
		t.Fatal("equal numeric values with distinct runtime types were deduplicated")
	}
}
