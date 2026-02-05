package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/vm"

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

		p1 := c.MustCompile(file.NewAnonymousSource(`
			LET val = ""
			RETURN { enabled: !val }
		`))

		v1, err := vm.New(p1).Run(gocontext.Background(), nil)

		So(err, ShouldBeNil)

		out1, err := v1.MarshalJSON()
		So(err, ShouldBeNil)

		So(string(out1), ShouldEqual, `{"enabled":true}`)

		p2 := c.MustCompile(file.NewAnonymousSource(`
			LET val = ""
			RETURN { enabled: !!val }
		`))

		v2, err := vm.New(p2).Run(gocontext.Background(), nil)

		So(err, ShouldBeNil)

		out2, err := v2.MarshalJSON()

		So(err, ShouldBeNil)
		So(string(out2), ShouldEqual, `{"enabled":false}`)
	})
}
