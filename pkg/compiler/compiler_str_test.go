package compiler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString(t *testing.T) {
	Convey("Should be possible to use multi line string", t, func() {
		out := compiler.New().
			MustCompile(`
			RETURN "
FOO
BAR
"
		`).
			MustRun(context.Background())

		So(string(out), ShouldEqual, `"\nFOO\nBAR\n"`)
	})

	Convey("Should be possible to use multi line string with nested strings", t, func() {
		out := compiler.New().
			MustCompile(fmt.Sprintf(`
RETURN %s<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<title>GetTitle</title>
		</head>
		<body>
			Hello world
		</body>
		</html>%s
`, "`", "`")).
			MustRun(context.Background())

		out, err := json.Marshal(`<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<title>GetTitle</title>
		</head>
		<body>
			Hello world
		</body>
		</html>`)

		So(err, ShouldBeNil)

		So(string(out), ShouldEqual, string(out))
	})
}

func BenchmarkStringLiteral(b *testing.B) {
	p := compiler.New().MustCompile(`
			RETURN "
FOO
BAR
"
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
