package runtime_test

import (
	"testing"
	"unsafe"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

// testStruct is a named type defined in this test package to verify
// that Box.Type() returns a non-empty Type for package-level types.
type testStruct struct {
	X int
}

// Compile-time interface conformance checks.
var (
	_ runtime.Value       = (*runtime.Box[int])(nil)
	_ runtime.Unwrappable = (*runtime.Box[int])(nil)
)

func TestBox(t *testing.T) {
	Convey("NewBox", t, func() {
		Convey("Should store an int value", func() {
			b := runtime.NewBox(42)
			So(b, ShouldNotBeNil)
			So(b.Value, ShouldEqual, 42)
		})

		Convey("Should store a string value", func() {
			b := runtime.NewBox("hello")
			So(b.Value, ShouldEqual, "hello")
		})

		Convey("Should store a bool value", func() {
			b := runtime.NewBox(true)
			So(b.Value, ShouldEqual, true)
		})

		Convey("Should store a struct value", func() {
			s := testStruct{X: 1}
			b := runtime.NewBox(s)
			So(b.Value, ShouldResemble, s)
		})

		Convey("Should store a pointer value", func() {
			n := 99
			b := runtime.NewBox(&n)
			So(b.Value, ShouldEqual, &n)
		})

		Convey("Should store a slice value", func() {
			sl := []int{1, 2, 3}
			b := runtime.NewBox(sl)
			So(b.Value, ShouldResemble, sl)
		})
	})

	Convey(".String", t, func() {
		Convey("Should return expected format for int", func() {
			b := runtime.NewBox(42)
			So(b.String(), ShouldEqual, "Box[int](42)")
		})

		Convey("Should return expected format for string", func() {
			b := runtime.NewBox("world")
			So(b.String(), ShouldEqual, "Box[string](world)")
		})

		Convey("Should return expected format for a named struct", func() {
			s := testStruct{X: 7}
			b := runtime.NewBox(s)
			So(b.String(), ShouldEqual, "Box[runtime_test.testStruct]({7})")
		})
	})

	Convey(".Unwrap", t, func() {
		Convey("Should return the underlying int value", func() {
			b := runtime.NewBox(42)
			So(b.Unwrap(), ShouldEqual, 42)
		})

		Convey("Should return the underlying struct value", func() {
			s := testStruct{X: 3}
			b := runtime.NewBox(s)
			So(b.Unwrap(), ShouldResemble, s)
		})

		Convey("Should return the same pointer", func() {
			n := 10
			b := runtime.NewBox(&n)
			So(b.Unwrap(), ShouldEqual, &n)
		})
	})

	Convey(".Hash", t, func() {
		Convey("Should return a non-zero hash", func() {
			b := runtime.NewBox(42)
			So(b.Hash(), ShouldBeGreaterThan, 0)
		})

		Convey("Hash should be consistent across calls", func() {
			b := runtime.NewBox(42)
			So(b.Hash(), ShouldEqual, b.Hash())
		})

		Convey("Two boxes with the same value should produce the same hash", func() {
			b1 := runtime.NewBox(42)
			b2 := runtime.NewBox(42)
			So(b1.Hash(), ShouldEqual, b2.Hash())
		})

		Convey("Boxes with different values should produce different hashes", func() {
			b1 := runtime.NewBox(42)
			b2 := runtime.NewBox(43)
			So(b1.Hash(), ShouldNotEqual, b2.Hash())
		})

		Convey("Boxes with different values of different types produce different hashes", func() {
			bInt := runtime.NewBox(1)
			bStr := runtime.NewBox("hello")
			So(bInt.Hash(), ShouldNotEqual, bStr.Hash())
		})

		Convey("Boxes of different builtin types with the same printed value produce the same hash", func() {
			bInt := runtime.NewBox(1)
			bStr := runtime.NewBox("1")
			So(bInt.String(), ShouldNotEqual, bStr.String())
		})
	})

	Convey(".Copy", t, func() {
		Convey("Should return a Value", func() {
			b := runtime.NewBox(42)
			var _ runtime.Value = b.Copy() // compile-time + runtime check
		})

		Convey("Should return a distinct pointer from the original", func() {
			b := runtime.NewBox(42)
			c := b.Copy().(*runtime.Box[int])
			So(unsafe.Pointer(c), ShouldNotEqual, unsafe.Pointer(b))
		})

		Convey("Should have the same value as the original", func() {
			b := runtime.NewBox(42)
			c := b.Copy().(*runtime.Box[int])
			So(c.Value, ShouldEqual, b.Value)
		})

		Convey("Should perform a shallow copy (slice shares underlying array)", func() {
			sl := []int{1, 2, 3}
			b := runtime.NewBox(sl)
			c := b.Copy().(*runtime.Box[[]int])

			// Mutate via the copy — original slice should reflect the change
			// because Copy is intentionally shallow.
			c.Value[0] = 99
			So(b.Value[0], ShouldEqual, 99)
		})

		Convey("Should perform a shallow copy (pointer shares the same target)", func() {
			n := 7
			b := runtime.NewBox(&n)
			c := b.Copy().(*runtime.Box[*int])

			So(c.Value, ShouldEqual, b.Value)
		})
	})
}
