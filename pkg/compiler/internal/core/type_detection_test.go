package core_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/runtime"
)

func TestRuntimeTypeToValueType(t *testing.T) {
	tests := []struct {
		name         string
		runtimeType  runtime.Type
		expectedType core.ValueType
	}{
		{"Int type", runtime.TypeInt, core.TypeInt},
		{"Float type", runtime.TypeFloat, core.TypeFloat},
		{"String type", runtime.TypeString, core.TypeString},
		{"Boolean type", runtime.TypeBoolean, core.TypeBool},
		{"Array type", runtime.TypeArray, core.TypeList},
		{"List type", runtime.TypeList, core.TypeList},
		{"Object type", runtime.TypeObject, core.TypeMap},
		{"Map type", runtime.TypeMap, core.TypeMap},
		{"Unknown type", runtime.TypeBinary, core.TypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := core.RuntimeTypeToValueType(tt.runtimeType)
			if result != tt.expectedType {
				t.Errorf("RuntimeTypeToValueType(%v) = %v, want %v", tt.runtimeType, result, tt.expectedType)
			}
		})
	}
}