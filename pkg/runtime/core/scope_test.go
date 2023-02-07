package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestScope(t *testing.T) {
	Convey(".SetVariable", t, func() {
		Convey("Should set a new variable", func() {
			rs, cf := core.NewRootScope()

			So(cf, ShouldNotBeNil)

			err := rs.SetVariable("foo", values.NewString("bar"))

			So(err, ShouldBeNil)
		})

		Convey("Should return an error when a variable is already defined", func() {
			rs, cf := core.NewRootScope()

			So(cf, ShouldNotBeNil)

			err := rs.SetVariable("foo", values.NewString("bar"))
			So(err, ShouldBeNil)

			err = rs.SetVariable("foo", values.NewString("bar"))
			So(err, ShouldHaveSameTypeAs, core.ErrNotUnique)
		})
	})

	Convey(".GetVariable", t, func() {
		Convey("Should set and get a variable", func() {
			rs, cf := core.NewRootScope()

			So(cf, ShouldNotBeNil)

			err := rs.SetVariable("foo", values.NewString("bar"))
			So(err, ShouldBeNil)

			v, err := rs.GetVariable("foo")

			So(err, ShouldBeNil)
			So(v, ShouldEqual, "bar")
		})

		Convey("Should return an error when variable is not defined", func() {
			rs, cf := core.NewRootScope()

			So(cf, ShouldNotBeNil)

			_, err := rs.GetVariable("foo")

			So(err, ShouldNotBeNil)
		})
	})

	Convey(".HasVariable", t, func() {
		Convey("Should return TRUE when a variable exists", func() {
			rs, cf := core.NewRootScope()

			So(cf, ShouldNotBeNil)

			err := rs.SetVariable("foo", values.NewString("bar"))
			So(err, ShouldBeNil)

			exists := rs.HasVariable("foo")

			So(exists, ShouldBeTrue)
		})

		Convey("Should return FALSE when a variable exists", func() {
			rs, cf := core.NewRootScope()

			So(cf, ShouldNotBeNil)

			exists := rs.HasVariable("foo")

			So(exists, ShouldBeFalse)
		})
	})

	Convey(".Fork", t, func() {
		Convey("Should create a nested scope", func() {
			Convey("Should set a variable only in a child scope", func() {
				rs, cf := core.NewRootScope()
				So(cf, ShouldNotBeNil)

				cs := rs.Fork()
				cs.SetVariable("foo", values.NewString("bar"))

				exists := rs.HasVariable("foo")

				So(exists, ShouldBeFalse)
			})

			Convey("Should return a variable defined only in a child scope", func() {
				rs, cf := core.NewRootScope()
				So(cf, ShouldNotBeNil)

				cs := rs.Fork()
				err := cs.SetVariable("foo", values.NewString("bar"))
				So(err, ShouldBeNil)

				v, err := cs.GetVariable("foo")

				So(err, ShouldBeNil)
				So(v, ShouldEqual, "bar")
			})

			Convey("Should return a variable defined only in a parent scope", func() {
				rs, cf := core.NewRootScope()
				So(cf, ShouldNotBeNil)

				cs := rs.Fork()
				err := cs.SetVariable("foo", values.NewString("bar"))
				So(err, ShouldBeNil)

				err = rs.SetVariable("faz", values.NewString("qaz"))
				So(err, ShouldBeNil)

				v, err := cs.GetVariable("faz")

				So(err, ShouldBeNil)
				So(v, ShouldEqual, "qaz")
			})

			Convey("Should set a new variable with a same name defined in a parent scope", func() {
				rs, cf := core.NewRootScope()
				So(cf, ShouldNotBeNil)

				err := rs.SetVariable("foo", values.NewString("bar"))
				So(err, ShouldBeNil)

				cs := rs.Fork()
				err = cs.SetVariable("foo", values.NewString("faz"))
				So(err, ShouldBeNil)

				rsV, err := rs.GetVariable("foo")
				So(err, ShouldBeNil)

				csV, err := cs.GetVariable("foo")
				So(err, ShouldBeNil)

				So(csV, ShouldNotEqual, rsV)
			})
		})
	})
}

func BenchmarkScope(b *testing.B) {
	root, _ := core.NewRootScope()

	for n := 0; n < b.N; n++ {
		root.Fork()
	}
}

type (
	TestCloserType struct{}

	TestCloserValue struct {
		closed bool
	}
)

func (tct TestCloserType) ID() int64 {
	return 99
}

func (tct TestCloserType) String() string {
	return "TestCloser"
}

func (tct TestCloserType) Equals(other core.Type) bool {
	return other.ID() == tct.ID()
}

func (tc *TestCloserValue) MarshalJSON() ([]byte, error) {
	return nil, core.ErrNotImplemented
}

func (tc *TestCloserValue) Type() core.Type {
	return TestCloserType{}
}

func (tc *TestCloserValue) String() string {
	return ""
}

func (tc *TestCloserValue) Compare(other core.Value) int64 {
	return 0
}

func (tc *TestCloserValue) Unwrap() interface{} {
	return tc
}

func (tc *TestCloserValue) Hash() uint64 {
	return 0
}

func (tc *TestCloserValue) Copy() core.Value {
	return &TestCloserValue{}
}

func (tc *TestCloserValue) Close() error {
	if tc.closed {
		return core.Error(core.ErrInvalidOperation, "already closed")
	}

	tc.closed = true

	return nil
}

func TestCloseFunc(t *testing.T) {
	Convey("Should close root scope and close all io.Closer values", t, func() {
		rs, cf := core.NewRootScope()

		tc := &TestCloserValue{}

		rs.SetVariable("disposable", tc)
		So(tc.closed, ShouldBeFalse)

		err := cf()
		So(err, ShouldBeNil)

		So(tc.closed, ShouldBeTrue)
	})

	Convey("Should return error if it's already closed", t, func() {
		rs, cf := core.NewRootScope()

		tc := &TestCloserValue{}

		rs.SetVariable("disposable", tc)
		So(tc.closed, ShouldBeFalse)

		err := cf()
		So(err, ShouldBeNil)

		So(tc.closed, ShouldBeTrue)

		err = cf()
		So(err, ShouldHaveSameTypeAs, core.ErrInvalidOperation)
	})
}
