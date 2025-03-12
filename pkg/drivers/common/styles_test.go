package common_test

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	. "github.com/smartystreets/goconvey/convey"
)

type style struct {
	raw   string
	name  core.String
	value core.Value
}

func TestDeserializeStyles(t *testing.T) {
	Convey("DeserializeStyles", t, func() {
		styles := []style{
			{
				raw:   "min-height: 1.15",
				name:  "min-height",
				value: core.NewFloat(1.15),
			},
			{
				raw:   "background-color: #4A154B",
				name:  "background-color",
				value: core.NewString("#4A154B"),
			},
			{
				raw:   "font-size:26pt",
				name:  "font-size",
				value: core.NewString("26pt"),
			},
			{
				raw:   "page-break-after:avoid",
				name:  "page-break-after",
				value: core.NewString("avoid"),
			},
			{
				raw:   `font-family: Arial,"Helvetica Neue",Helvetica,sans-serif`,
				name:  "font-family",
				value: core.NewString(`Arial,"Helvetica Neue",Helvetica,sans-serif`),
			},
			{
				raw:   "color: black",
				name:  "color",
				value: core.NewString("black"),
			},
			{
				raw:   "display: inline-block",
				name:  "display",
				value: core.NewString("inline-block"),
			},
			{
				raw:   "min-width: 50",
				name:  "min-width",
				value: core.NewFloat(50),
			},
		}

		Convey("Should parse a single style", func() {
			for _, s := range styles {
				out, err := common.DeserializeStyles(core.NewString(s.raw))

				So(err, ShouldBeNil)
				So(out, ShouldNotBeNil)

				value, exists := out.Get(s.name)

				So(bool(exists), ShouldBeTrue)

				So(value.Compare(s.value) == 0, ShouldBeTrue)
			}
		})

		Convey("Should parse multiple styles", func() {
			var buff bytes.Buffer

			for _, s := range styles {
				buff.WriteString(s.raw)
				buff.WriteString("; ")
			}

			out, err := common.DeserializeStyles(core.NewString(buff.String()))

			So(err, ShouldBeNil)
			So(out, ShouldNotBeNil)
			So(int(out.Length()), ShouldEqual, len(styles))

			for _, s := range styles {
				value, exists := out.Get(s.name)

				So(bool(exists), ShouldBeTrue)

				So(value.Compare(s.value) == 0, ShouldBeTrue)
			}
		})
	})
}
