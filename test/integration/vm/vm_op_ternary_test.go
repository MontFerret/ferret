package vm_test

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/MontFerret/ferret/test/integration/base"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func TestTernaryOperator(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case("RETURN 1 < 2 ? 3 : 4", 3),
		Case("RETURN 1 > 2 ? 3 : 4", 4),
		Case("RETURN 2 ? : 4", 2),
		Case("LET foo = TRUE RETURN foo ? TRUE : FALSE", true),
		Case("LET foo = FALSE RETURN foo ? TRUE : FALSE", false),
		CaseArray("FOR i IN [1, 2, 3, 4, 5, 6] RETURN i < 3 ? i * 3 : i * 2", []any{3, 6, 6, 8, 10, 12}),
		CaseArray(`FOR i IN [NONE, 2, 3, 4, 5, 6] RETURN i ? : i`, []any{nil, 2, 3, 4, 5, 6}),
		Case(`RETURN 0 && true ? "1" : "some"`, "some"),
		Case(`RETURN length([]) > 0 && true ? "1" : "some"`, "some"),
	})

	Convey("Should compile ternary operator with default values", t, func() {
		vals := []string{
			"0",
			"0.0",
			"''",
			"NONE",
			"FALSE",
		}

		c := compiler.New()

		for _, val := range vals {
			p, err := c.Compile(fmt.Sprintf(`
			FOR i IN [%s, 1, 2, 3]
				RETURN i ? i * 2 : 'no value'
		`, val))

			So(err, ShouldBeNil)

			out, err := Run(p)

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `["no value",2,4,6]`)
		}
	})
}
