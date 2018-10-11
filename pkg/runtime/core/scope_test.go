package core_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

func TestScope(t *testing.T) {
	Convey("Should match", t, func() {
		rs, cf := core.NewRootScope()

		So(cf, ShouldNotBeNil)

		s := core.NewScope(rs)

		So(s.HasVariable("a"), ShouldBeFalse)

		s.SetVariable("a", values.NewString("test"))

		So(s.HasVariable("a"), ShouldBeTrue)

		v, err := s.GetVariable("a")

		So(err, ShouldBeNil)
		So(v, ShouldEqual, "test")

		c := s.Fork()

		So(c.HasVariable("a"), ShouldBeTrue)

		cv, err := c.GetVariable("a")

		So(err, ShouldBeNil)
		So(cv, ShouldEqual, "test")
	})
}

func TestScopeTraversing(t *testing.T) {
	Convey("Should match", t, func() {
		rs, cf := core.NewRootScope()
		So(cf, ShouldNotBeNil)

		s := core.NewScope(rs)

		rs.SetVariable("a", values.NewString("test"))
		v, err := s.GetVariable("a")

		// root traversal should work
		So(err, ShouldBeNil)
		So(v, ShouldEqual, "test")

		s.SetVariable("b", values.NewString("test2"))
		s2 := core.NewScope(rs)
		_, err = s2.GetVariable("b")

		// child traversal should fail
		So(err, ShouldNotBeNil)

		v2, err := s2.GetVariable("a")

		// root traversal should work
		So(err, ShouldBeNil)
		So(v2, ShouldEqual, "test")

	})
}
