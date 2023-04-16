package runtime_v2_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime_v2"
)

func TestProgram(t *testing.T) {
	Convey("Disassemble", t, func() {
		Convey("Constants", t, func() {
			Convey("Should return a string", func() {
				p := runtime_v2.NewProgram([]runtime_v2.Opcode{
					runtime_v2.OpConstant,
				}, []core.Value{
					values.String("test"),
				})

				out := p.Disassemble()

				So(out, ShouldEqual, "const 0\n")
			})
		})
	})
}
