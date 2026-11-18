package debugger

import (
	"errors"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type hostileDebugValue struct {
	infoCalls *int
	info      runtime.DebugInfo
	panicInfo bool
}

func (hostileDebugValue) String() string      { panic("String called") }
func (hostileDebugValue) Hash() uint64        { panic("Hash called") }
func (hostileDebugValue) Copy() runtime.Value { panic("Copy called") }
func (hostileDebugValue) Type() runtime.Type  { panic("Type called") }
func (v hostileDebugValue) DebugInfo() runtime.DebugInfo {
	if v.infoCalls != nil {
		(*v.infoCalls)++
	}
	if v.panicInfo {
		panic("DebugInfo called")
	}

	return v.info
}

func TestValueAccessRejectsOpaqueKeys(t *testing.T) {
	access := vm.NewDebugValueAccess()
	calls := 0
	key := hostileDebugValue{infoCalls: &calls}
	for _, value := range []runtime.Value{runtime.NewObject(), runtime.NewArray(0)} {
		if _, err := access.Lookup(value, key); !errors.Is(err, runtime.ErrInvalidArgument) || calls != 0 {
			t.Fatalf("lookup on %T with opaque key: err=%v, hook calls=%d", value, err, calls)
		}
	}
}

func TestFormatSemanticValues(t *testing.T) {
	access := vm.NewDebugValueAccess()
	pattern, err := runtime.NewRegexp("a+")
	if err != nil {
		t.Fatal(err)
	}

	date := runtime.NewDateTime(time.Date(2026, 2, 3, 4, 5, 6, 0, time.FixedZone(strings.Repeat("z", 10000), 3600)))
	for _, tc := range []struct {
		value         runtime.Value
		name, display string
	}{
		{nil, "None", "NONE"},
		{runtime.None, "None", "NONE"},
		{runtime.True, "Boolean", "true"},
		{runtime.Int(42), "Int", "42"},
		{runtime.Float(3.14), "Float", "3.14"},
		{runtime.NewDuration(1500 * time.Millisecond), "Duration", "1.5s"},
		{date, "DateTime", "2026-02-03T03:05:06Z"},
		{pattern, "Regexp", "<Regexp>"},
		{runtime.String("hello"), "String", `"hello"`},
		{runtime.Binary(make([]byte, 1<<20)), "Binary", "Binary(1048576)"},
		{runtime.NewArray(8), "Array", "Array(0)"},
		{runtime.NewObject(), "Object", "Object(0)"},
		{runtime.NewSizedArray(1000000), "Array", "Array(1000000)"},
		{(*runtime.Array)(nil), "Array", "<Array>"},
		{(*runtime.Object)(nil), "Object", "<Object>"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := inspector{values: access, format: DefaultFormatOptions()}
			if got := s.debugValue(tc.value); got != (Value{Type: tc.name, Display: tc.display}) {
				t.Fatalf("got %+v, want %s %s with reference zero", got, tc.name, tc.display)
			}
		})
	}
}

func TestFormatNativeScalarPointers(t *testing.T) {
	s := inspector{values: vm.NewDebugValueAccess(), format: DefaultFormatOptions()}
	boolean, integer, float := runtime.True, runtime.Int(42), runtime.Float(3.14)
	duration, date := runtime.NewDuration(time.Second), runtime.NewDateTime(time.Unix(42, 0))
	str, bin := runtime.String(strings.Repeat("é", 1<<20)), runtime.Binary(make([]byte, 4096))
	for _, tc := range []struct{ value, pointer, absent runtime.Value }{
		{boolean, &boolean, (*runtime.Boolean)(nil)},
		{integer, &integer, (*runtime.Int)(nil)},
		{float, &float, (*runtime.Float)(nil)},
		{duration, &duration, (*runtime.Duration)(nil)},
		{str, &str, (*runtime.String)(nil)},
		{date, &date, (*runtime.DateTime)(nil)},
		{bin, &bin, (*runtime.Binary)(nil)},
	} {
		want := s.debugValue(tc.value)
		if got := s.debugValue(tc.pointer); got != want {
			t.Errorf("%T = %+v, want %+v", tc.pointer, got, want)
		}

		if got := s.debugValue(tc.absent); got != (Value{Type: want.Type, Display: "<" + want.Type + ">"}) {
			t.Errorf("nil %T = %+v", tc.absent, got)
		}
	}

	if got := s.debugValue(&hostileDebugEvalValue{}); got != (Value{Type: "Resource", Display: "<Resource>"}) {
		t.Fatalf("opaque pointer = %+v", got)
	}
}

func TestBoundedText(t *testing.T) {
	for _, tc := range []struct {
		name   string
		input  string
		want   string
		max    int
		quoted bool
	}{
		{name: "zero", input: "value", max: 0, quoted: false, want: ""},
		{name: "negative", input: "value", max: -1, quoted: true, want: ""},
		{name: "one", input: "value", max: 1, quoted: false, want: "."},
		{name: "two", input: "value", max: 2, quoted: true, want: ".."},
		{name: "empty quoted", input: "", max: 2, quoted: true, want: `""`},
		{name: "short quoted", input: "x", max: 3, quoted: true, want: `"x"`},
		{name: "quoting overhead", input: "abc", max: 4, quoted: true, want: "..."},
		{name: "only marker fits", input: "abcd", max: 5, quoted: true, want: `"..."`},
		{name: "ascii", input: "product", max: 4, quoted: false, want: "p..."},
		{name: "rune exact fit", input: "🙂", max: 4, quoted: false, want: "🙂"},
		{name: "rune cut", input: "éclair", max: 5, quoted: false, want: "é..."},
		{name: "rune too wide", input: "🙂x", max: 3, quoted: false, want: "..."},
		{name: "escapes", input: "\n\"\\\x01", max: 12, quoted: true, want: `"\n\"\\\x01"`},
		{name: "escape boundary", input: "\x01\x01x", max: 8, quoted: true, want: `"..."`},
		{name: "quoted multibyte", input: "é🙂abcdef", max: 7, quoted: true, want: `"é..."`},
		{name: "visible invalid", input: "a\xffb", max: 7, quoted: true, want: `"a�b"`},
		{name: "large with invalid suffix", input: strings.Repeat("a", 4<<20) + "\xff", max: 8, quoted: true, want: `"aaa..."`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := renderText(tc.input, tc.max, tc.quoted)
			if got != tc.want || len(got) > max(0, tc.max) || !utf8.ValidString(got) {
				t.Fatalf("rendered %q, want %q", got, tc.want)
			}
		})
	}
}
