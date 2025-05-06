package runtime_test

import (
	"context"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateArgs(t *testing.T) {
	Convey("Should match", t, func() {
		a := []runtime.Value{runtime.NewInt(1), runtime.NewInt(2)}

		e := runtime.ValidateArgs(a, 1, 2)
		So(e, ShouldBeNil)

		e = runtime.ValidateArgs(a, 3, 4)
		So(e, ShouldNotBeNil)
	})
}

func TestFunctions(t *testing.T) {

	fnTrue := func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
		return runtime.True, nil
	}

	Convey(".Set", t, func() {

		Convey("Should set function by name", func() {
			fns := runtime.NewFunctions()
			fname := "F"

			fns.Set(fname, fnTrue)

			names := fns.Names()

			So(names, ShouldHaveLength, 1)
			So(names[0], ShouldEqual, fname)
		})

		Convey("Should set function by name at uppercase", func() {
			fns := runtime.NewFunctions()
			fname := "f"

			fns.Set(fname, fnTrue)

			names := fns.Names()

			So(names, ShouldHaveLength, 1)
			So(names[0], ShouldEqual, strings.ToUpper(fname))
		})

		Convey("Should set when Functions created not by NewFunctions", func() {
			fns := runtime.Functions{}
			fname := "F"

			fns.Set(fname, fnTrue)

			names := fns.Names()

			So(names, ShouldHaveLength, 1)
			So(names[0], ShouldEqual, fname)
		})
	})

	Convey(".Get", t, func() {

		Convey("Should get function by name", func() {
			fns := runtime.NewFunctions()
			fname := "F"
			fns.Set(fname, fnTrue)

			fn, exists := fns.Get(fname)

			So(exists, ShouldBeTrue)
			So(fn, ShouldEqual, fnTrue)
		})

		Convey("Should get function by name at uppercase", func() {
			fns := runtime.NewFunctions()
			fname := "f"
			fns.Set(fname, fnTrue)

			fn, exists := fns.Get(strings.ToUpper(fname))

			So(exists, ShouldBeTrue)
			So(fn, ShouldEqual, fnTrue)
		})

		Convey("Should not panic when Functions created not by NewFunctions", func() {
			fns := runtime.Functions{}

			fn, exists := fns.Get("f")

			So(exists, ShouldBeFalse)
			So(fn, ShouldBeNil)
		})
	})

	Convey(".Unset", t, func() {

		Convey("Should unset function by name", func() {
			fns := runtime.NewFunctions()
			fname := "F"
			fns.Set(fname, fnTrue)

			fns.Unset(fname)

			So(fns.Names(), ShouldHaveLength, 0)
		})

		Convey("Should get function by name at uppercase", func() {
			fns := runtime.NewFunctions()
			fname := "f"
			fns.Set(fname, fnTrue)

			fns.Unset(strings.ToUpper(fname))

			So(fns.Names(), ShouldHaveLength, 0)
		})

		Convey("Should not panic when Functions created not by NewFunctions", func() {
			fns := runtime.Functions{}
			fname := "F"
			fns.Set(fname, fnTrue)

			fns.Unset(fname)

			So(fns.Names(), ShouldHaveLength, 0)
		})
	})

	Convey(".Names", t, func() {

		Convey("Should return name", func() {
			fns := runtime.NewFunctions()
			fname := "F"
			fns.Set(fname, fnTrue)

			names := fns.Names()

			So(names, ShouldHaveLength, 1)
			So(names, ShouldContain, fname)
		})

		Convey("Should return name at uppercase", func() {
			fns := runtime.NewFunctions()
			fname := "f"
			fns.Set(fname, fnTrue)

			names := fns.Names()

			So(names, ShouldHaveLength, 1)
			So(names, ShouldContain, strings.ToUpper(fname))
		})

		Convey("Should not panic when Functions created not by NewFunctions", func() {
			fns := runtime.Functions{}
			So(fns.Names(), ShouldHaveLength, 0)
		})
	})
}
