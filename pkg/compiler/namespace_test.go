package compiler_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNamespaceBuilder(t *testing.T) {
	Convey("Namespaces", t, func() {
		Convey("Should return an error when a function name contains NS separator", func() {
			c := compiler.New()
			err := c.RegisterFunction("FOO::SPY", func(ctx context.Context, args ...core.Value) (value core.Value, e error) {
				return values.None, nil
			})

			So(err, ShouldNotBeNil)
		})

		Convey("Should successfully register a name within a namespace", func() {
			c := compiler.New()
			err := c.Namespace("FOO").RegisterFunction("SPY", func(ctx context.Context, args ...core.Value) (value core.Value, e error) {
				return values.None, nil
			})

			So(err, ShouldBeNil)

			funcs := c.RegisteredFunctions()

			var exists bool

			for _, name := range funcs {
				exists = name == "FOO::SPY"

				if exists {
					break
				}
			}

			So(exists, ShouldBeTrue)
		})

		Convey("Root namespace should return all registered functions", func() {
			c := compiler.New()
			err := c.Namespace("FOO").RegisterFunction("SPY", func(ctx context.Context, args ...core.Value) (value core.Value, e error) {
				return values.None, nil
			})

			So(err, ShouldBeNil)

			funcs := c.RegisteredFunctions()

			So(len(funcs), ShouldBeGreaterThan, 1)
		})

		Convey("Namespace should return all registered functions", func() {
			c := compiler.New()
			err := c.Namespace("FOO").RegisterFunction("SPY", func(ctx context.Context, args ...core.Value) (value core.Value, e error) {
				return values.None, nil
			})

			So(err, ShouldBeNil)

			err = c.Namespace("FOO").Namespace("UTILS").RegisterFunction("SPY", func(ctx context.Context, args ...core.Value) (value core.Value, e error) {
				return values.None, nil
			})

			So(err, ShouldBeNil)

			funcs := c.Namespace("FOO").RegisteredFunctions()

			So(funcs, ShouldHaveLength, 2)

			funcs2 := c.Namespace("FOO").Namespace("UTILS").RegisteredFunctions()

			So(funcs2, ShouldHaveLength, 1)
		})

		Convey("Namespace should return an error if namespace name is incorrect", func() {
			c := compiler.New()
			noop := func(ctx context.Context, args ...core.Value) (value core.Value, e error) {
				return values.None, nil
			}
			err := c.Namespace("FOO::").RegisterFunction("SPY", noop)

			So(err, ShouldNotBeNil)

			err = c.Namespace("@F").RegisterFunction("SPY", noop)

			So(err, ShouldNotBeNil)
		})
	})
}
