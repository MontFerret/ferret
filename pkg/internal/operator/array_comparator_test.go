package operator_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
)

func TestArrayComparatorContract(t *testing.T) {
	tests := []struct {
		binary operator.Binary
		array  operator.ArrayComparator
		offset int
	}{
		{binary: operator.Equal, array: operator.ArrayEqual, offset: 0},
		{binary: operator.NotEqual, array: operator.ArrayNotEqual, offset: 1},
		{binary: operator.Greater, array: operator.ArrayGreater, offset: 2},
		{binary: operator.GreaterOrEqual, array: operator.ArrayGreaterOrEqual, offset: 3},
		{binary: operator.Less, array: operator.ArrayLess, offset: 4},
		{binary: operator.LessOrEqual, array: operator.ArrayLessOrEqual, offset: 5},
		{binary: operator.In, array: operator.ArrayIn, offset: 6},
	}

	for _, test := range tests {
		t.Run(test.binary.String(), func(t *testing.T) {
			if actual := int(test.array); actual != test.offset {
				t.Fatalf("numeric value = %d, want %d", actual, test.offset)
			}

			decoded, ok := operator.ArrayComparatorFromOffset(test.offset)
			if !ok || decoded != test.array {
				t.Fatalf("ArrayComparatorFromOffset(%d) = %v, %v; want %v, true", test.offset, decoded, ok, test.array)
			}

			encoded, ok := operator.ArrayComparatorFor(test.binary)
			if !ok || encoded != test.array {
				t.Fatalf("ArrayComparatorFor(%v) = %v, %v; want %v, true", test.binary, encoded, ok, test.array)
			}

			binary, ok := test.array.Binary()
			if !ok || binary != test.binary {
				t.Fatalf("Binary() = %v, %v; want %v, true", binary, ok, test.binary)
			}
		})
	}
}

func TestArrayComparatorMatchesQuantifiedOpcodeGroups(t *testing.T) {
	groups := []struct {
		name     string
		opcodes  []bytecode.Opcode
		baseCode bytecode.Opcode
	}{
		{name: "ANY", baseCode: bytecode.OpAnyEq, opcodes: []bytecode.Opcode{bytecode.OpAnyEq, bytecode.OpAnyNe, bytecode.OpAnyGt, bytecode.OpAnyGte, bytecode.OpAnyLt, bytecode.OpAnyLte, bytecode.OpAnyIn}},
		{name: "NONE", baseCode: bytecode.OpNoneEq, opcodes: []bytecode.Opcode{bytecode.OpNoneEq, bytecode.OpNoneNe, bytecode.OpNoneGt, bytecode.OpNoneGte, bytecode.OpNoneLt, bytecode.OpNoneLte, bytecode.OpNoneIn}},
		{name: "ALL", baseCode: bytecode.OpAllEq, opcodes: []bytecode.Opcode{bytecode.OpAllEq, bytecode.OpAllNe, bytecode.OpAllGt, bytecode.OpAllGte, bytecode.OpAllLt, bytecode.OpAllLte, bytecode.OpAllIn}},
	}

	for _, group := range groups {
		t.Run(group.name, func(t *testing.T) {
			for offset, expected := range group.opcodes {
				actual := bytecode.Opcode(int(group.baseCode) + offset)
				if actual != expected {
					t.Fatalf("opcode at offset %d = %s, want %s", offset, actual, expected)
				}
				if _, ok := operator.ArrayComparatorFromOffset(offset); !ok {
					t.Fatalf("offset %d is not a valid ArrayComparator", offset)
				}
			}
		})
	}
}

func TestArrayComparatorRejectsInvalidInput(t *testing.T) {
	for _, offset := range []int{-1, 7, 255, 256} {
		if actual, ok := operator.ArrayComparatorFromOffset(offset); ok || actual != operator.UnknownArrayComparator {
			t.Fatalf("ArrayComparatorFromOffset(%d) = %v, %v; want UnknownArrayComparator, false", offset, actual, ok)
		}
	}

	for _, binary := range []operator.Binary{operator.Unknown, operator.Add, operator.Divide} {
		if actual, ok := operator.ArrayComparatorFor(binary); ok || actual != operator.UnknownArrayComparator {
			t.Fatalf("ArrayComparatorFor(%v) = %v, %v; want UnknownArrayComparator, false", binary, actual, ok)
		}
	}

	if binary, ok := operator.UnknownArrayComparator.Binary(); ok || binary != operator.Unknown {
		t.Fatalf("UnknownArrayComparator.Binary() = %v, %v; want Unknown, false", binary, ok)
	}
}
