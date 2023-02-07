package core_test

import (
	"context"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestValidateArgs(t *testing.T) {
	Convey("Should match", t, func() {
		a := []core.Value{values.NewInt(1), values.NewInt(2)}

		e := core.ValidateArgs(a, 1, 2)
		So(e, ShouldBeNil)

		e = core.ValidateArgs(a, 3, 4)
		So(e, ShouldNotBeNil)
	})
}

func TestFunctions(t *testing.T) {

	fnTrue := func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return values.True, nil
	}

	Convey(".Set", t, func() {

		Convey("Should set function by name", func() {
			fns := core.NewFunctions()
			fname := "F"

			fns.Set(fname, fnTrue)

			names := fns.Names()

			So(names, ShouldHaveLength, 1)
			So(names[0], ShouldEqual, fname)
		})

		Convey("Should set function by name at uppercase", func() {
			fns := core.NewFunctions()
			fname := "f"

			fns.Set(fname, fnTrue)

			names := fns.Names()

			So(names, ShouldHaveLength, 1)
			So(names[0], ShouldEqual, strings.ToUpper(fname))
		})

		Convey("Should set when Functions created not by NewFunctions", func() {
			fns := core.Functions{}
			fname := "F"

			fns.Set(fname, fnTrue)

			names := fns.Names()

			So(names, ShouldHaveLength, 1)
			So(names[0], ShouldEqual, fname)
		})
	})

	Convey(".Get", t, func() {

		Convey("Should get function by name", func() {
			fns := core.NewFunctions()
			fname := "F"
			fns.Set(fname, fnTrue)

			fn, exists := fns.Get(fname)

			So(exists, ShouldBeTrue)
			So(fn, ShouldEqual, fnTrue)
		})

		Convey("Should get function by name at uppercase", func() {
			fns := core.NewFunctions()
			fname := "f"
			fns.Set(fname, fnTrue)

			fn, exists := fns.Get(strings.ToUpper(fname))

			So(exists, ShouldBeTrue)
			So(fn, ShouldEqual, fnTrue)
		})

		Convey("Should not panic when Functions created not by NewFunctions", func() {
			fns := core.Functions{}

			fn, exists := fns.Get("f")

			So(exists, ShouldBeFalse)
			So(fn, ShouldBeNil)
		})
	})

	Convey(".Unset", t, func() {

		Convey("Should unset function by name", func() {
			fns := core.NewFunctions()
			fname := "F"
			fns.Set(fname, fnTrue)

			fns.Unset(fname)

			So(fns.Names(), ShouldHaveLength, 0)
		})

		Convey("Should get function by name at uppercase", func() {
			fns := core.NewFunctions()
			fname := "f"
			fns.Set(fname, fnTrue)

			fns.Unset(strings.ToUpper(fname))

			So(fns.Names(), ShouldHaveLength, 0)
		})

		Convey("Should not panic when Functions created not by NewFunctions", func() {
			fns := core.Functions{}
			fname := "F"
			fns.Set(fname, fnTrue)

			fns.Unset(fname)

			So(fns.Names(), ShouldHaveLength, 0)
		})
	})

	Convey(".Names", t, func() {

		Convey("Should return name", func() {
			fns := core.NewFunctions()
			fname := "F"
			fns.Set(fname, fnTrue)

			names := fns.Names()

			So(names, ShouldHaveLength, 1)
			So(names, ShouldContain, fname)
		})

		Convey("Should return name at uppercase", func() {
			fns := core.NewFunctions()
			fname := "f"
			fns.Set(fname, fnTrue)

			names := fns.Names()

			So(names, ShouldHaveLength, 1)
			So(names, ShouldContain, strings.ToUpper(fname))
		})

		Convey("Should not panic when Functions created not by NewFunctions", func() {
			fns := core.Functions{}
			So(fns.Names(), ShouldHaveLength, 0)
		})
	})
}
