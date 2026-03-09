package benchmarks_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type udfArityCase struct {
	name   string
	params string
	args   string
}

var udfArityCases = []udfArityCase{
	{name: "A0", params: "", args: ""},
	{name: "A1", params: "a", args: "1"},
	{name: "A2", params: "a, b", args: "1, 2"},
	{name: "A3", params: "a, b, c", args: "1, 2, 3"},
	{name: "A4", params: "a, b, c, d", args: "1, 2, 3, 4"},
	{name: "A6", params: "a, b, c, d, e, f", args: "1, 2, 3, 4, 5, 6"},
}

func udfReturnExpr(params string) string {
	if params == "" {
		return "1"
	}

	parts := strings.Split(params, ",")
	return strings.TrimSpace(parts[0])
}

func udfTopLevelQuery(a udfArityCase) string {
	expr := udfReturnExpr(a.params)

	return fmt.Sprintf(`
FUNC f(%s) => %s
RETURN f(%s)
`, a.params, expr, a.args)
}

func udfNestedQuery(a udfArityCase) string {
	expr := udfReturnExpr(a.params)

	return fmt.Sprintf(`
FUNC outer(%s) (
  FUNC inner(%s) (
    RETURN %s
  )
  LET v = inner(%s)
  RETURN v
)
RETURN outer(%s)
`, a.params, a.params, expr, a.args, a.args)
}

func hostTopLevelQuery(a udfArityCase) string {
	return fmt.Sprintf(`
RETURN TEST(%s)
`, a.args)
}

func hostOption(a udfArityCase) vm.EnvironmentOption {
	return withBuilder(func(b *runtime.FunctionsBuilder) {
		switch a.name {
		case "A0":
			b.A0().Add("TEST", func(ctx context.Context) (runtime.Value, error) {
				return runtime.True, nil
			})
		case "A1":
			b.A1().Add("TEST", func(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
				return runtime.True, nil
			})
		case "A2":
			b.A2().Add("TEST", func(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
				return runtime.True, nil
			})
		case "A3":
			b.A3().Add("TEST", func(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
				return runtime.True, nil
			})
		case "A4":
			b.A4().Add("TEST", func(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
				return runtime.True, nil
			})
		default:
			b.Var().Add("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
				return runtime.True, nil
			})
		}
	})
}

func BenchmarkUdfCalls(b *testing.B) {
	for _, c := range udfArityCases {
		c := c

		b.Run(fmt.Sprintf("UDF/TopLevel/%s/O0", c.name), func(b *testing.B) {
			RunBenchmarkO0(b, udfTopLevelQuery(c))
		})
		b.Run(fmt.Sprintf("UDF/TopLevel/%s/O1", c.name), func(b *testing.B) {
			RunBenchmarkO1(b, udfTopLevelQuery(c))
		})
		b.Run(fmt.Sprintf("UDF/Nested/%s/O0", c.name), func(b *testing.B) {
			RunBenchmarkO0(b, udfNestedQuery(c))
		})
		b.Run(fmt.Sprintf("UDF/Nested/%s/O1", c.name), func(b *testing.B) {
			RunBenchmarkO1(b, udfNestedQuery(c))
		})
	}
}

func BenchmarkUdfCalls_HostBaseline(b *testing.B) {
	for _, c := range udfArityCases {
		c := c

		b.Run(fmt.Sprintf("Host/TopLevel/%s/O0", c.name), func(b *testing.B) {
			RunBenchmarkO0(b, hostTopLevelQuery(c), hostOption(c))
		})

		b.Run(fmt.Sprintf("Host/TopLevel/%s/O1", c.name), func(b *testing.B) {
			RunBenchmarkO1(b, hostTopLevelQuery(c), hostOption(c))
		})
	}
}
