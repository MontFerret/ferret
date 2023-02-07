package types_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func TestHelpers(t *testing.T) {
	Convey("Compare", t, func() {
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

		Convey("None", func() {
			So(types.Compare(types.None, types.None), ShouldEqual, 0)

			for _, t := range typesList[1:] {
				So(types.Compare(types.None, t), ShouldEqual, -1)
			}
		})

		Convey("Boolean", func() {
			for _, t := range typesList {
				switch t.ID() {
				case types.None.ID():
					So(types.Compare(types.Boolean, t), ShouldEqual, 1)
				case types.Boolean.ID():
					So(types.Compare(types.Boolean, t), ShouldEqual, 0)
				default:
					So(types.Compare(types.Boolean, t), ShouldEqual, -1)
				}
			}
		})

		Convey("Int", func() {
			for _, t := range typesList {
				switch t.ID() {
				case types.None.ID():
					So(types.Compare(types.Int, t), ShouldEqual, 1)
				case types.Boolean.ID():
					So(types.Compare(types.Int, t), ShouldEqual, 1)
				case types.Int.ID():
					So(types.Compare(types.Int, t), ShouldEqual, 0)
				default:
					So(types.Compare(types.Int, t), ShouldEqual, -1)
				}
			}
		})

		Convey("Float", func() {
			for _, t := range typesList {
				switch t.ID() {
				case types.None.ID():
					So(types.Compare(types.Float, t), ShouldEqual, 1)
				case types.Boolean.ID():
					So(types.Compare(types.Float, t), ShouldEqual, 1)
				case types.Int.ID():
					So(types.Compare(types.Float, t), ShouldEqual, 1)
				case types.Float.ID():
					So(types.Compare(types.Float, t), ShouldEqual, 0)
				default:
					So(types.Compare(types.Float, t), ShouldEqual, -1)
				}
			}
		})

		Convey("String", func() {
			for _, t := range typesList {
				switch t.ID() {
				case types.None.ID():
					So(types.Compare(types.String, t), ShouldEqual, 1)
				case types.Boolean.ID():
					So(types.Compare(types.String, t), ShouldEqual, 1)
				case types.Int.ID():
					So(types.Compare(types.String, t), ShouldEqual, 1)
				case types.Float.ID():
					So(types.Compare(types.String, t), ShouldEqual, 1)
				case types.String.ID():
					So(types.Compare(types.String, t), ShouldEqual, 0)
				default:
					So(types.Compare(types.String, t), ShouldEqual, -1)
				}
			}
		})

		Convey("DateTime", func() {
			for _, t := range typesList {
				switch t.ID() {
				case types.None.ID():
					So(types.Compare(types.DateTime, t), ShouldEqual, 1)
				case types.Boolean.ID():
					So(types.Compare(types.DateTime, t), ShouldEqual, 1)
				case types.Int.ID():
					So(types.Compare(types.DateTime, t), ShouldEqual, 1)
				case types.Float.ID():
					So(types.Compare(types.DateTime, t), ShouldEqual, 1)
				case types.String.ID():
					So(types.Compare(types.DateTime, t), ShouldEqual, 1)
				case types.DateTime.ID():
					So(types.Compare(types.DateTime, t), ShouldEqual, 0)
				default:
					So(types.Compare(types.DateTime, t), ShouldEqual, -1)
				}
			}
		})

		Convey("Array", func() {
			for _, t := range typesList {
				switch t.ID() {
				case types.None.ID():
					So(types.Compare(types.Array, t), ShouldEqual, 1)
				case types.Boolean.ID():
					So(types.Compare(types.Array, t), ShouldEqual, 1)
				case types.Int.ID():
					So(types.Compare(types.Array, t), ShouldEqual, 1)
				case types.Float.ID():
					So(types.Compare(types.Array, t), ShouldEqual, 1)
				case types.String.ID():
					So(types.Compare(types.Array, t), ShouldEqual, 1)
				case types.DateTime.ID():
					So(types.Compare(types.Array, t), ShouldEqual, 1)
				case types.Array.ID():
					So(types.Compare(types.Array, t), ShouldEqual, 0)
				default:
					So(types.Compare(types.Array, t), ShouldEqual, -1)
				}
			}
		})

		Convey("Object", func() {
			for _, t := range typesList {
				switch t.ID() {
				case types.None.ID():
					So(types.Compare(types.Object, t), ShouldEqual, 1)
				case types.Boolean.ID():
					So(types.Compare(types.Object, t), ShouldEqual, 1)
				case types.Int.ID():
					So(types.Compare(types.Object, t), ShouldEqual, 1)
				case types.Float.ID():
					So(types.Compare(types.Object, t), ShouldEqual, 1)
				case types.String.ID():
					So(types.Compare(types.Object, t), ShouldEqual, 1)
				case types.DateTime.ID():
					So(types.Compare(types.Object, t), ShouldEqual, 1)
				case types.Array.ID():
					So(types.Compare(types.Object, t), ShouldEqual, 1)
				case types.Object.ID():
					So(types.Compare(types.Object, t), ShouldEqual, 0)
				default:
					So(types.Compare(types.Object, t), ShouldEqual, -1)
				}
			}
		})

		Convey("Binary", func() {
			for _, t := range typesList {
				switch t.ID() {
				case types.Binary.ID():
					So(types.Compare(types.Binary, t), ShouldEqual, 0)
				default:
					So(types.Compare(types.Binary, t), ShouldEqual, 1)
				}
			}
		})
	})
}
