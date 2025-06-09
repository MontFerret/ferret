package vm_test

import (
	"fmt"
	. "github.com/MontFerret/ferret/test/integration/base"
	"testing"
)

func TestString(t *testing.T) {
	RunUseCases(t, []UseCase{
		Case(
			`
			RETURN "
FOO
BAR
"
		`, "\nFOO\nBAR\n", "Should be possible to use multi line string"),

		CaseJSON(
			fmt.Sprintf(`
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
`, "`", "`"), `<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<title>GetTitle</title>
		</head>
		<body>
			Hello world
		</body>
		</html>`, "Should be possible to use multi line string with nested strings using backtick"),

		CaseJSON(
			fmt.Sprintf(`
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
`, "´", "´"),
			`<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<title>GetTitle</title>
		</head>
		<body>
			Hello world
		</body>
		</html>`, "Should be possible to use multi line string with nested strings using tick"),
	})
}
