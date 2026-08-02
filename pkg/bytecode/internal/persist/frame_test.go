package persist

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestDurationConstantRoundTrip(t *testing.T) {
	t.Parallel()

	expected := runtime.NewDuration(1500 * time.Millisecond)
	frame, err := encodeConstant(expected)
	if err != nil {
		t.Fatal(err)
	}
	if frame.Type != constantTypeDuration || frame.Duration == nil || *frame.Duration != "1.5s" {
		t.Fatalf("unexpected duration frame: %#v", frame)
	}

	actual, err := decodeConstant(frame)
	if err != nil {
		t.Fatal(err)
	}
	if actual != expected {
		t.Fatalf("decoded duration = %v, want %v", actual, expected)
	}

	invalid := "not-a-duration"
	if _, err := decodeConstant(ConstantFrame{Type: constantTypeDuration, Duration: &invalid}); !errors.Is(err, bytecode.ErrInvalidConstant) {
		t.Fatalf("malformed duration error = %v", err)
	}
}

func TestValidateInstructionOperandForBitSize(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		bitSize int
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
