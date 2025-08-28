package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
	. "github.com/smartystreets/goconvey/convey"
)

func TestErrorSpan(t *testing.T) {
	Convey("ErrorSpan constructors", t, func() {
		Convey("NewErrorSpan should create ErrorSpan with all fields", func() {
			span := file.Span{Start: 0, End: 10}
			label := "test label"
			main := true

			result := NewErrorSpan(span, label, main)

			So(result.Span, ShouldEqual, span)
			So(result.Label, ShouldEqual, label)
			So(result.Main, ShouldEqual, main)
		})

		Convey("NewMainErrorSpan should create main ErrorSpan", func() {
			span := file.Span{Start: 0, End: 10}
			label := "main error"

			result := NewMainErrorSpan(span, label)

			So(result.Span, ShouldEqual, span)
			So(result.Label, ShouldEqual, label)
			So(result.Main, ShouldBeTrue)
		})

		Convey("NewSecondaryErrorSpan should create non-main ErrorSpan", func() {
			span := file.Span{Start: 5, End: 15}
			label := "secondary error"

			result := NewSecondaryErrorSpan(span, label)

			So(result.Span, ShouldEqual, span)
			So(result.Label, ShouldEqual, label)
			So(result.Main, ShouldBeFalse)
		})
	})
}