package core_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
)

// TODO: Move the tests to "values" module

func TestTypeString(t *testing.T) {
	Convey("The string representation of the type should match this type", t, func() {
		So(core.Type(0).String(), ShouldEqual, "none")
		So(core.Type(1).String(), ShouldEqual, "boolean")
		So(core.Type(2).String(), ShouldEqual, "int")
		So(core.Type(3).String(), ShouldEqual, "float")
		So(core.Type(4).String(), ShouldEqual, "string")
		So(core.Type(5).String(), ShouldEqual, "datetime")
		So(core.Type(6).String(), ShouldEqual, "array")
		So(core.Type(7).String(), ShouldEqual, "object")
		So(core.Type(8).String(), ShouldEqual, "binary")
	})
}

func TestIsTypeOf(t *testing.T) {
	Convey("Check type by value", t, func() {

		So(core.IsTypeOf(values.None, core.NoneType), ShouldBeTrue)
		So(core.IsTypeOf(values.NewBoolean(true), core.BooleanType), ShouldBeTrue)
		So(core.IsTypeOf(values.NewInt(1), core.IntType), ShouldBeTrue)
		So(core.IsTypeOf(values.NewFloat(1.1), core.FloatType), ShouldBeTrue)
		So(core.IsTypeOf(values.NewString("test"), core.StringType), ShouldBeTrue)
		So(core.IsTypeOf(values.NewDateTime(time.Now()), core.DateTimeType), ShouldBeTrue)
		So(core.IsTypeOf(values.NewArray(1), core.ArrayType), ShouldBeTrue)
		So(core.IsTypeOf(values.NewObject(), core.ObjectType), ShouldBeTrue)
		So(core.IsTypeOf(values.NewBinary([]byte{}), core.BinaryType), ShouldBeTrue)
	})
}

func TestValidateType(t *testing.T) {
	Convey("Value should match type", t, func() {

		So(core.ValidateType(values.None, core.NoneType), ShouldBeNil)
		So(core.ValidateType(values.NewBoolean(true), core.BooleanType), ShouldBeNil)
		So(core.ValidateType(values.NewInt(1), core.IntType), ShouldBeNil)
		So(core.ValidateType(values.NewFloat(1.1), core.FloatType), ShouldBeNil)
		So(core.ValidateType(values.NewString("test"), core.StringType), ShouldBeNil)
		So(core.ValidateType(values.NewDateTime(time.Now()), core.DateTimeType), ShouldBeNil)
		So(core.ValidateType(values.NewArray(1), core.ArrayType), ShouldBeNil)
		So(core.ValidateType(values.NewObject(), core.ObjectType), ShouldBeNil)
		So(core.ValidateType(values.NewBinary([]byte{}), core.BinaryType), ShouldBeNil)
	})

	Convey("Value should not match type", t, func() {
		So(core.ValidateType(values.None, core.BooleanType), ShouldBeError)
		So(core.ValidateType(values.NewBoolean(true), core.IntType, core.NoneType), ShouldBeError)
		So(core.ValidateType(values.NewInt(1), core.NoneType), ShouldBeError)
		So(core.ValidateType(values.NewFloat(1.1), core.StringType), ShouldBeError)
		So(core.ValidateType(values.NewString("test"), core.IntType, core.FloatType), ShouldBeError)
		So(core.ValidateType(values.NewDateTime(time.Now()), core.BooleanType), ShouldBeError)
		So(core.ValidateType(values.NewArray(1), core.StringType), ShouldBeError)
		So(core.ValidateType(values.NewObject(), core.BooleanType), ShouldBeError)
		So(core.ValidateType(values.NewBinary([]byte{}), core.NoneType), ShouldBeError)
	})
}
