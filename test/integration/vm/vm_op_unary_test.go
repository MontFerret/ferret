package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/vm"

	gocontext "context"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnaryOperators(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN !TRUE", false),
		Case("RETURN NOT TRUE", false),
		Case("RETURN !FALSE", true),
		Case("RETURN -1", -1),
		Case("RETURN -1.1", -1.1),
		Case("RETURN +1", 1),
		Case("RETURN +1.1", 1.1),
		Case("LET v = 1 RETURN -v", -1),
		Case("LET v = 1.1 RETURN -v", -1.1),
		Case("LET v = -1 RETURN -v", 1),
		Case("LET v = -1.1 RETURN -v", 1.1),
		Case("LET v = -1 RETURN +v", -1),
		Case("LET v = -1.1 RETURN +v", -1.1),
	})

	Convey("RETURN { enabled: !val}", t, func() {
		c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))

		p1, err := c.Compile(file.NewAnonymousSource(`
			LET val = ""
			RETURN { enabled: !val }
		`))

		So(err, ShouldBeNil)

		vm1, err := vm.New(p1)
		So(err, ShouldBeNil)

		r1, err := vm1.Run(gocontext.Background(), nil)

		So(err, ShouldBeNil)

		out1, err := encodingjson.Default.Encode(r1.Root())
		So(r1.Close(), ShouldBeNil)
		So(err, ShouldBeNil)

		So(string(out1), ShouldEqual, `{"enabled":true}`)

		p2, err := c.Compile(file.NewAnonymousSource(`
			LET val = ""
			RETURN { enabled: !!val }
		`))

		So(err, ShouldBeNil)

		vm2, err := vm.New(p2)
		So(err, ShouldBeNil)

		r2, err := vm2.Run(gocontext.Background(), nil)

		So(err, ShouldBeNil)

		out2, err := encodingjson.Default.Encode(r2.Root())
		So(r2.Close(), ShouldBeNil)

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `{"enabled":false}`)
	})
}
