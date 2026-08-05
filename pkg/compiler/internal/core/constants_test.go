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

	if integer.Hash() != floating.Hash() {
		t.Fatal("equal numeric values must have equal hashes")
	}

	// Int and Float compare equal numerically, but constants retain their runtime type.
	floatOperand := pool.Add(floating)
	if floatOperand == integerOperand {
		t.Fatal("equal numeric values with distinct runtime types were deduplicated")
	}
}

func TestConstantPoolUsesStrictDurationEquality(t *testing.T) {
	pool := NewConstantPool()
	stringOperand := pool.Add(runtime.NewString("1s"))
	intOperand := pool.Add(runtime.NewInt(1000))
	durationOperand := pool.Add(runtime.NewDuration(time.Second))
	equivalentDurationOperand := pool.Add(runtime.NewDuration(1000 * time.Millisecond))

	if stringOperand == intOperand || stringOperand == durationOperand || intOperand == durationOperand {
		t.Fatal("cross-type Duration representations were deduplicated")
	}
	if durationOperand != equivalentDurationOperand {
		t.Fatal("equivalent native Duration constants were not deduplicated")
	}
}
