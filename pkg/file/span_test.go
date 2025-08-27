package file

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpan(t *testing.T) {
	Convey("Span struct", t, func() {
		Convey("Should create span with start and end", func() {
			span := Span{Start: 5, End: 10}

			So(span.Start, ShouldEqual, 5)
			So(span.End, ShouldEqual, 10)
		})

		Convey("Should handle zero values", func() {
			span := Span{}

			So(span.Start, ShouldEqual, 0)
			So(span.End, ShouldEqual, 0)
		})

		Convey("Should handle negative values", func() {
			span := Span{Start: -1, End: -5}

			So(span.Start, ShouldEqual, -1)
			So(span.End, ShouldEqual, -5)
		})
	})
}
