package vm_test

import (
	"fmt"
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

		Case(
			`
LET name = "World"
RETURN `+"`"+`Hello ${name}!`+"`"+`
		`, "Hello World!", "Should interpolate template literals"),

		Case(
			`
RETURN `+"`"+`sum=${1 + 2}`+"`"+`
		`, "sum=3", "Should interpolate expressions inside template literals"),

		Case(
			`
RETURN `+"`"+`${1 + 2}`+"`"+`
		`, "3", "Template literals should coerce expressions to strings"),

		CaseObject(
			`
RETURN { `+"`"+`foo${1}`+"`"+`: 2 }
		`, map[string]any{"foo1": 2}, "Template literals should work as property names"),
	})
}
