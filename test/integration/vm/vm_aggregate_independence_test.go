package vm_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/format"
	jsonformat "github.com/MontFerret/ferret/v2/pkg/bytecode/format/json"
	msgpackformat "github.com/MontFerret/ferret/v2/pkg/bytecode/format/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestCollectAggregateIndependentOfMath(t *testing.T) {
	custom := runtime.NewLibrary()
	custom.Function().A1().Add("custom_count", func(ctx context.Context, v runtime.Value) (runtime.Value, error) {
		return v.(runtime.List).Length(ctx)
	})
	custom.Namespace("custom").Function().A1().Add("sum", func(context.Context, runtime.Value) (runtime.Value, error) { return runtime.Int(99), nil })
	tests := []struct{ name, query, want string }{
		{"case insensitive", `RETURN FOR v IN [1,none,3] COLLECT g=1 AGGREGATE a=sUm(v),b=aVeRaGe(v) RETURN [a,b]`, `[[4,2]]`},
		{"generic grouped", `RETURN FOR v IN [1,"2",none,true,3] COLLECT g=1 AGGREGATE a=SUM(v),b=AVERAGE(v) RETURN [g,a,b]`, `[[1,4,2]]`},
		{"generic extrema", `RETURN FOR v IN [1,"2",none,true,3] COLLECT g=1 AGGREGATE a=MIN(v),b=MAX(v) RETURN [a,b]`, `[[1,3]]`},
		{"fused grouped", `RETURN FOR v IN [1,"2",none,true,3] COLLECT g=1 AGGREGATE a=SUM(v),b=AVERAGE(v),c=COUNT(v),d=MIN(v),e=MAX(v) RETURN [a,b,c,d,e]`, `[[4,2,5,1,3]]`},
		{"global", `RETURN FOR v IN [1,"2",none,true,3] COLLECT AGGREGATE a=SUM(v),b=AVERAGE(v),c=COUNT(v),d=MIN(v),e=MAX(v) RETURN [a,b,c,d,e]`, `[[4,2,5,1,3]]`},
		{"generic global", `RETURN FOR v IN [1,"2",none,true,3] COLLECT AGGREGATE a=SUM(v),b=AVERAGE(v),c=custom_count(v) RETURN [a,b,c]`, `[[4,2,5]]`},
		{"protected", `RETURN FOR v IN [1,none,3] COLLECT g=1 AGGREGATE a=SUM(v)?,b=AVERAGE(v)? RETURN [a,b]`, `[[4,2]]`},
		{"global protected", `RETURN FOR v IN [1,none,3] COLLECT AGGREGATE a=SUM(v)? RETURN a`, `[4]`},
		{"generic into", `RETURN FOR v IN [1,none,3] COLLECT g=1 AGGREGATE a=SUM(v) INTO rows RETURN [a,LENGTH(rows)]`, `[[4,3]]`},
		{"generic global into", `RETURN FOR v IN [1,none,3] COLLECT AGGREGATE a=SUM(v),b=custom_count(v) INTO rows RETURN [a,b,LENGTH(rows)]`, `[[4,3,3]]`},
		{"global empty", `RETURN FOR v IN [] COLLECT AGGREGATE a=SUM(v),b=custom_count(v) RETURN [a,b]`, `[[null,null]]`},
		{"grouped empty", `RETURN FOR v IN [] COLLECT g=1 AGGREGATE a=SUM(v) RETURN a`, `[]`},
		{"generic nonnumeric", `RETURN FOR v IN [none,"2",true] COLLECT g=1 AGGREGATE a=SUM(v),b=AVERAGE(v) RETURN [a,b]`, `[[0,0]]`},
		{"fused nonnumeric", `RETURN FOR v IN [none,"2",true] COLLECT g=1 AGGREGATE a=SUM(v),b=AVERAGE(v),c=MIN(v),d=MAX(v) RETURN [a,b,c,d]`, `[[0,0,null,null]]`},
		{"custom global name", `RETURN FOR v IN [1,2] COLLECT AGGREGATE a=custom::sum(v) RETURN a`, `[99]`},
		{"custom grouped name", `RETURN FOR v IN [1,2] COLLECT g=1 AGGREGATE a=custom::sum(v),b=SUM(v),c=COUNT(v) RETURN [a,b,c]`, `[[99,3,2]]`},
	}

	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		for _, test := range tests {
			t.Run(level.String()+"/"+test.name, func(t *testing.T) {
				program, err := c.Compile(t.Context(), source.New(test.name, test.query))
				if err != nil {
					t.Fatal(err)
				}

				for _, fn := range program.Functions.Host {
					if !strings.HasPrefix(fn.Name, "custom") {
						t.Fatalf("query built-in depends on public function %+v", fn)
					}
				}

				programs := []*bytecode.Program{program}
				for _, codec := range []format.Format{jsonformat.Default, msgpackformat.Default} {
					encoded, err := codec.Marshal(program)
					if err != nil {
						t.Fatal(err)
					}

					decoded, err := codec.Unmarshal(encoded)
					if err != nil {
						t.Fatal(err)
					}

					programs = append(programs, decoded)
				}

				for _, program := range programs {
					out, err := spec.Run(program, vm.WithNamespace(custom))
					if err != nil {
						t.Fatal(err)
					}

					var got, want any
					if err := json.Unmarshal(out, &got); err != nil {
						t.Fatal(err)
					}

					if err := json.Unmarshal([]byte(test.want), &want); err != nil {
						t.Fatal(err)
					}

					if !reflect.DeepEqual(got, want) {
						t.Fatalf("got %s, want %s", out, test.want)
					}
				}
			})
		}
	}
}

func TestNamespacedMathAggregatesRemainStrict(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		for _, query := range []string{
			`RETURN FOR v IN [1,"2",3] COLLECT AGGREGATE a=math::sum(v) RETURN a`,
			`RETURN FOR v IN [1,"2",3] COLLECT g=1 AGGREGATE a=math::sum(v),b=math::min(v),c=math::max(v) RETURN [a,b,c]`,
			`USE math AS m RETURN FOR v IN [1,"2",3] COLLECT AGGREGATE a=m::sum(v) RETURN a`,
			`RETURN math::percentile([1,2],50,"interpolation")`,
			`RETURN math::percentile([],101)`,
			`RETURN math::percentile([1,2],-1)`,
		} {
			t.Run(fmt.Sprintf("%s/%s", level, query), func(t *testing.T) {
				program, err := c.Compile(t.Context(), source.New("strict", query))
				if err != nil {
					t.Fatal(err)
				}

				if _, err := spec.Run(program, vm.WithNamespace(stdlib.New())); err == nil {
					t.Fatal("expected canonical validation or arity error")
				}
			})
		}
	}
}
