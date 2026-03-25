package vm_test

import (
	"fmt"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	"testing"
)

func TestString(t *testing.T) {
	RunSpecs(t, []Spec{
		S(
			`
			RETURN "
FOO
BAR
"
		`, "\nFOO\nBAR\n", "Should be possible to use multi line string"),

		JSON(
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

		JSON(
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

		S(
			`
LET name = "World"
RETURN `+"`"+`Hello ${name}!`+"`"+`
		`, "Hello World!", "Should interpolate template literals"),

		S(
			`
RETURN `+"`"+`sum=${1 + 2}`+"`"+`
		`, "sum=3", "Should interpolate expressions inside template literals"),

		S(
			`
RETURN `+"`"+`${1 + 2}`+"`"+`
		`, "3", "Template literals should coerce expressions to strings"),

		Object(
			`
RETURN { `+"`"+`foo${1}`+"`"+`: 2 }
		`, map[string]any{"foo1": 2}, "Template literals should work as property names"),

		S(
			fmt.Sprintf("\nRETURN %scost=\\${1}%s\n", "`", "`"),
			"cost=${1}",
			"Should allow escaping interpolation start in template literals"),

		S(
			fmt.Sprintf("\nRETURN %suse \\`backtick\\`%s\n", "`", "`"),
			"use `backtick`",
			"Should allow escaped backticks in template literals"),
	})
}
