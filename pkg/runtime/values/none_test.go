package values_test

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNone(t *testing.T) {
	Convey("None", t, func() {
		Convey(".Compare", func() {
			valuesList := []core.Value{
				values.None,
				values.True,
				values.ZeroInt,
				values.ZeroFloat,
				values.EmptyString,
				values.ZeroDateTime,
				values.NewArray(1),
				values.NewObject(),
				values.NewBinary(make([]byte, 0)),
			}

			So(values.None.Compare(values.None), ShouldEqual, 0)

			for _, v := range valuesList[1:] {
				So(values.None.Compare(v), ShouldEqual, -1)
			}
		})
	})
}
