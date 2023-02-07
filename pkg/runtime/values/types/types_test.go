package types_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type TestValue struct {
	t core.Type
}

func (v TestValue) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (v TestValue) Type() core.Type {
	return v.t
}

func (v TestValue) String() string {
	return ""
}

func (v TestValue) Compare(other core.Value) int64 {
	return 0
}

func (v TestValue) Unwrap() interface{} {
	return nil
}

func (v TestValue) Hash() uint64 {
	return 0
}

func (v TestValue) Copy() core.Value {
	return v
}

func TestType(t *testing.T) {
	Convey(".GetName", t, func() {
		So(types.None.String(), ShouldEqual, "none")
		So(types.Boolean.String(), ShouldEqual, "boolean")
		So(types.Int.String(), ShouldEqual, "int")
		So(types.Float.String(), ShouldEqual, "float")
		So(types.String.String(), ShouldEqual, "string")
		So(types.DateTime.String(), ShouldEqual, "date_time")
		So(types.Array.String(), ShouldEqual, "array")
		So(types.Object.String(), ShouldEqual, "object")
		So(types.Binary.String(), ShouldEqual, "binary")
	})

	Convey("==", t, func() {
		typesList := []core.Type{
			types.None,
			types.Boolean,
			types.Int,
			types.Float,
			types.String,
			types.DateTime,
			types.Array,
			types.Object,
			types.Binary,
		}

		valuesList := []core.Value{
			TestValue{types.None},
			TestValue{types.Boolean},
			TestValue{types.Int},
			TestValue{types.Float},
			TestValue{types.String},
			TestValue{types.DateTime},
			TestValue{types.Array},
			TestValue{types.Object},
			TestValue{types.Binary},
		}

		for i, t := range typesList {
			So(t == valuesList[i].Type(), ShouldBeTrue)
		}
	})
}
