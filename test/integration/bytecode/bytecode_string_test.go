package bytecode_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestString(t *testing.T) {
	RunUseCases(t, []UseCase{
		ByteCodeCase(
			`
			RETURN "
FOO
BAR
"
		`, []vm.Instruction{
				I(vm.OpLoadConst, 1, C(0)),
				I(vm.OpMove, 0, R(1)),
				I(vm.OpReturn, 0),
			}, "Should be possible to use multi line string"),
		//
		//		CaseJSON(
		//			fmt.Sprintf(`
		//RETURN %s<!DOCTYPE html>
		//		<html lang="en">
		//		<head>
		//		<meta charset="UTF-8">
		//		<title>GetTitle</title>
		//		</head>
		//		<body>
		//			Hello world
		//		</body>
		//		</html>%s
		//`, "`", "`"), `<!DOCTYPE html>
		//		<html lang="en">
		//		<head>
		//		<meta charset="UTF-8">
		//		<title>GetTitle</title>
		//		</head>
		//		<body>
		//			Hello world
		//		</body>
		//		</html>`, "Should be possible to use multi line string with nested strings using backtick"),
		//
		//		CaseJSON(
		//			fmt.Sprintf(`
		//RETURN %s<!DOCTYPE html>
		//		<html lang="en">
		//		<head>
		//		<meta charset="UTF-8">
		//		<title>GetTitle</title>
		//		</head>
		//		<body>
		//			Hello world
		//		</body>
		//		</html>%s
		//`, "´", "´"),
		//			`<!DOCTYPE html>
		//		<html lang="en">
		//		<head>
		//		<meta charset="UTF-8">
		//		<title>GetTitle</title>
		//		</head>
		//		<body>
		//			Hello world
		//		</body>
		//		</html>`, "Should be possible to use multi line string with nested strings using tick"),
	})
}
