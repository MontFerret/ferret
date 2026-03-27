package persist

import (
	"errors"
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestValidateInstructionOperandForBitSize(t *testing.T) {
	tests := []struct {
		value   int64
		bitSize int
		name    string
		wantErr bool
	}{
		{
			name:    "32bit_max",
			value:   1<<31 - 1,
			bitSize: 32,
		},
		{
			name:    "32bit_min",
			value:   -1 << 31,
			bitSize: 32,
		},
		{
			name:    "32bit_above_max",
			value:   1 << 31,
			bitSize: 32,
			wantErr: true,
		},
		{
			name:    "32bit_below_min",
			value:   (-1 << 31) - 1,
			bitSize: 32,
			wantErr: true,
		},
		{
			name:    "64bit_max",
			value:   1<<63 - 1,
			bitSize: 64,
		},
		{
			name:    "64bit_min",
			value:   -1 << 63,
			bitSize: 64,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateInstructionOperandForBitSize(tc.value, tc.bitSize)
			if tc.wantErr {
				if !errors.Is(err, bytecode.ErrInvalidProgram) {
					t.Fatalf("expected ErrInvalidProgram, got %v", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestDecodeInstructionOperandAcceptsPlatformBounds(t *testing.T) {
	minValue, maxValue := operandRangeForBitSize(strconv.IntSize)

	for _, value := range []int64{minValue, maxValue, 0, -1, 1} {
		got, err := decodeInstructionOperand(value)
		if err != nil {
			t.Fatalf("unexpected error for %d: %v", value, err)
		}

		if int64(got) != value {
			t.Fatalf("unexpected operand for %d: got %d", value, got)
		}
	}
}
