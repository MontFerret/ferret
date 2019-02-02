package types_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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
	Convey(".Name", t, func() {
		So(types.None.Name(), ShouldEqual, "none")
		So(types.Boolean.Name(), ShouldEqual, "boolean")
		So(types.Int.Name(), ShouldEqual, "int")
		So(types.Float.Name(), ShouldEqual, "float")
		So(types.String.Name(), ShouldEqual, "string")
		So(types.DateTime.Name(), ShouldEqual, "date_time")
		So(types.Array.Name(), ShouldEqual, "array")
		So(types.Object.Name(), ShouldEqual, "object")
		So(types.Binary.Name(), ShouldEqual, "binary")
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
