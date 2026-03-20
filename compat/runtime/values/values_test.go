package values_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/compat/runtime/core"
	"github.com/MontFerret/ferret/v2/compat/runtime/values"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// --- Type alias identity tests ---

func TestTypeAlias_Int(t *testing.T) {
	var v values.Int = values.NewInt(42)
	var rv runtime.Int = v // should be same type

	if rv != 42 {
		t.Fatalf("expected 42, got %d", rv)
	}
}

func TestTypeAlias_Float(t *testing.T) {
	var v values.Float = values.NewFloat(3.14)
	var rv runtime.Float = v

	if rv != 3.14 {
		t.Fatalf("expected 3.14, got %f", rv)
	}
}

func TestTypeAlias_String(t *testing.T) {
	var v values.String = values.NewString("hello")
	var rv runtime.String = v

	if rv != "hello" {
		t.Fatalf("expected \"hello\", got %q", rv)
	}
}

func TestTypeAlias_Boolean(t *testing.T) {
	var v values.Boolean = values.NewBoolean(true)
	var rv runtime.Boolean = v

	if !bool(rv) {
		t.Fatal("expected true")
	}
}

// --- Constants ---

func TestConstants(t *testing.T) {
	if values.True != true {
		t.Fatal("True should be true")
	}

	if values.False != false {
		t.Fatal("False should be false")
	}

	if values.ZeroInt != 0 {
		t.Fatal("ZeroInt should be 0")
	}

	if values.ZeroFloat != 0.0 {
		t.Fatal("ZeroFloat should be 0.0")
	}

	if values.EmptyString != "" {
		t.Fatal("EmptyString should be empty")
	}

	if values.SpaceString != " " {
		t.Fatal("SpaceString should be a space")
	}

	if values.None == nil {
		t.Fatal("NoneV2 should not be nil")
	}
}

// --- Constructor functions ---

func TestNewArray(t *testing.T) {
	arr := values.NewArray(5)
	if arr == nil {
		t.Fatal("expected non-nil array")
	}
}

func TestNewArrayWith(t *testing.T) {
	arr := values.NewArrayWith(values.NewInt(1), values.NewInt(2), values.NewInt(3))
	if arr == nil {
		t.Fatal("expected non-nil array")
	}
}

func TestEmptyArray(t *testing.T) {
	arr := values.EmptyArray()
	if arr == nil {
		t.Fatal("expected non-nil array")
	}
}

func TestNewObject(t *testing.T) {
	obj := values.NewObject()
	if obj == nil {
		t.Fatal("expected non-nil object")
	}
}

func TestNewBinary(t *testing.T) {
	b := values.NewBinary([]byte{1, 2, 3})
	if b == nil {
		t.Fatal("expected non-nil binary")
	}
}

// --- Parse ---

func TestParse(t *testing.T) {
	tests := []struct {
		input interface{}
		check func(core.Value) bool
		name  string
	}{
		{"int", 42, func(v core.Value) bool { return v.String() == "42" }},
		{"string", "hello", func(v core.Value) bool { return v.String() == "hello" }},
		{"bool", true, func(v core.Value) bool { return v.String() == "true" }},
		{"float", 3.14, func(v core.Value) bool { return v.String() == "3.14" }},
		{"nil", nil, func(v core.Value) bool { return v.String() == "" }},
	}

	for _, tc := range tests {
		v := values.Parse(tc.input)
		if v == nil {
			t.Fatalf("%s: expected non-nil value", tc.name)
		}

		if !tc.check(v) {
			t.Fatalf("%s: check failed, got %q", tc.name, v.String())
		}
	}
}

// --- Conversion helpers ---

func TestToBoolean(t *testing.T) {
	if values.ToBoolean(core.WrapValue(values.NewInt(1))) != true {
		t.Fatal("expected true for int 1")
	}

	if values.ToBoolean(core.WrapValue(values.NewInt(0))) != false {
		t.Fatal("expected false for int 0")
	}

	if values.ToBoolean(core.WrapValue(values.NewString(""))) != false {
		t.Fatal("expected false for empty string")
	}

	if values.ToBoolean(core.WrapValue(values.NewString("x"))) != true {
		t.Fatal("expected true for non-empty string")
	}
}

func TestToInt(t *testing.T) {
	result := values.ToInt(core.WrapValue(values.NewFloat(3.7)))
	if result != 3 {
		t.Fatalf("expected 3, got %d", result)
	}

	result = values.ToInt(core.WrapValue(values.NewString("42")))
	if result != 42 {
		t.Fatalf("expected 42, got %d", result)
	}
}

func TestToFloat(t *testing.T) {
	result := values.ToFloat(core.WrapValue(values.NewInt(5)))
	if result != 5.0 {
		t.Fatalf("expected 5.0, got %f", result)
	}
}

func TestToString(t *testing.T) {
	result := values.ToString(core.WrapValue(values.NewInt(42)))
	if result != "42" {
		t.Fatalf("expected \"42\", got %q", result)
	}
}

func TestToIntDefault(t *testing.T) {
	// Valid conversion
	result := values.ToIntDefault(core.WrapValue(values.NewString("10")), values.Int(99))
	if result != 10 {
		t.Fatalf("expected 10, got %d", result)
	}

	// Invalid conversion — should return default
	result = values.ToIntDefault(core.WrapValue(values.NewString("not_a_number")), values.Int(99))
	if result != 99 {
		t.Fatalf("expected 99 (default), got %d", result)
	}
}

// --- ObjectProperty wrapper ---

func TestNewObjectWith(t *testing.T) {
	obj := values.NewObjectWith(
		values.NewObjectProperty("name", values.NewString("ferret")),
		values.NewObjectProperty("version", values.NewInt(2)),
	)

	if obj == nil {
		t.Fatal("expected non-nil object")
	}
}

func TestNewObjectProperty(t *testing.T) {
	prop := values.NewObjectProperty("key", values.NewString("value"))

	if prop.Name != "key" {
		t.Fatalf("expected name \"key\", got %q", prop.Name)
	}
}

// --- NaN ---

func TestNaN(t *testing.T) {
	nan := values.NaN()
	if values.IsNaN(nan) != true {
		t.Fatal("expected NaN")
	}
}
