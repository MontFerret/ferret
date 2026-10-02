package benchmarks_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type compilerTailCallWorkload struct {
	name       string
	query      string
	forwarders int
	captures   int
	recoveries int
}

func BenchmarkCompilerTailCallElimination(b *testing.B) {
	queries := []compilerTailCallWorkload{
		{name: "NoUDF", query: `return @value + 1`},
		{name: "OrdinaryUDF", query: `func inc(value) => value + 1
return inc(1) + inc(2)`},
		{name: "TailChain", query: `func first(value) => second(value)
func second(value) => third(value)
func third(value) => value * 3
return first(2)`},
		{name: "Recursive", query: `func recur(value) => recur(NEXT(value))
return recur(256) on error return none`},
		{name: "Capture", query: compilerUdfLifecycleQuery},
		{name: "Recovery", query: `func target(value) => value / 0
func forward(value) => target(value) on error retry 2 or return value
return forward(3)`},
	}
	for _, forwarders := range []int{16, 64, 256} {
		queries = append(queries,
			buildCompilerTailCallWorkload("ManyUDF", forwarders, 0, 0),
			buildCompilerTailCallWorkload("CapturedUDF", forwarders, 8, 0),
			buildCompilerTailCallWorkload("RecoveryHeavy", forwarders, 8, 8),
		)
	}

	modes := []struct {
		name    string
		options []compiler.Option
		level   compiler.OptimizationLevel
	}{
		{"Default", nil, compiler.Full},
		{"None", []compiler.Option{compiler.WithOptimizationLevel(compiler.None)}, compiler.None},
		{"Basic", []compiler.Option{compiler.WithOptimizationLevel(compiler.Basic)}, compiler.Basic},
		{"Full", []compiler.Option{compiler.WithOptimizationLevel(compiler.Full)}, compiler.Full},
		{"Debug", []compiler.Option{compiler.WithOptimizationLevel(compiler.Full), compiler.WithDebugInfo()}, compiler.None},
	}
	tailsInLowering := compilerTailCallsInLowering(b)

	for _, query := range queries {
		for _, mode := range modes {
			b.Run(query.name+"/"+mode.name, func(b *testing.B) {
				instance := mustNewCompiler(b, mode.options...)
				src := source.New("tail-call-benchmark.fql", query.query)

				// Warm parser initialization and validate the workload before timing.
				program, err := instance.Compile(b.Context(), src)
				if err != nil {
					b.Fatal(err)
				}

				if program.Metadata.OptimizationLevel != int(mode.level) {
					b.Fatalf("effective optimization level = %d, want %v", program.Metadata.OptimizationLevel, mode.level)
				}

				if query.forwarders > 0 {
					checkCompilerTailCallWorkload(b, program, query, tailsInLowering || mode.level != compiler.None)
				}

				b.ReportAllocs()
				b.ResetTimer()

				for range b.N {
					if _, err := instance.Compile(b.Context(), src); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}

func buildCompilerTailCallWorkload(name string, forwarders, captures, recoveries int) compilerTailCallWorkload {
	var query strings.Builder
	for capture := range captures {
		fmt.Fprintf(&query, "var captured%d = @seed + %d\n", capture, capture)
	}

	for index := range forwarders {
		fmt.Fprintf(&query, "func f%d(value) {\n", index)
		if index == 0 {
			// A read-only VAR capture uses snapshots rather than cell handles.
			for capture := range captures {
				fmt.Fprintf(&query, "captured%d += value\n", capture)
			}
		}

		for recovery := range recoveries {
			fmt.Fprintf(&query, "let recovered%d = STEP(value, %d) on error return value\n", recovery, recovery)
		}

		target := fmt.Sprintf("f%d", index+1)
		if index+1 == forwarders {
			target = "leaf"
		}

		fmt.Fprintf(&query, "return %s(value", target)
		for capture := range captures {
			fmt.Fprintf(&query, " + captured%d", capture)
		}

		for recovery := range recoveries {
			fmt.Fprintf(&query, " + recovered%d", recovery)
		}

		query.WriteString(")\n}\n")
	}

	query.WriteString("func leaf(value) => value\nreturn f0(@value)")

	return compilerTailCallWorkload{
		name:       fmt.Sprintf("%s/UDFs%d", name, forwarders),
		query:      query.String(),
		forwarders: forwarders,
		captures:   captures,
		recoveries: recoveries,
	}
}

func compilerTailCallsInLowering(b *testing.B) bool {
	b.Helper()

	// The pre-refactor PR base emits tail calls even under None and Debug.
	// Detect that legacy behavior only for validation; the measured workloads
	// and compiler options are identical on both revisions.
	program, err := mustNewCompiler(b, compiler.WithOptimizationLevel(compiler.None)).Compile(b.Context(), source.NewAnonymous(`
func leaf(value) => value
func forward(value) => leaf(value)
return forward(@value)`))
	if err != nil {
		b.Fatal(err)
	}

	calls, tails := 0, 0
	for _, instruction := range program.Bytecode {
		switch instruction.Opcode {
		case bytecode.OpCall:
			calls++
		case bytecode.OpTailCall:
			tails++
		}
	}

	if calls == 1 && tails == 1 {
		return true
	}

	if calls != 2 || tails != 0 {
		b.Fatalf("unexpected None call/tail probe = %d/%d", calls, tails)
	}

	return false
}

func checkCompilerTailCallWorkload(b *testing.B, program *bytecode.Program, workload compilerTailCallWorkload, wantTail bool) {
	b.Helper()

	functions := program.Functions.UserDefined
	if len(functions) != workload.forwarders+1 {
		b.Fatalf("UDF count = %d, want %d", len(functions), workload.forwarders+1)
	}

	if got, want := len(program.CatchTable), workload.forwarders*workload.recoveries; got != want {
		b.Fatalf("catch count = %d, want %d", got, want)
	}

	wantOpcode := bytecode.OpCall
	if wantTail {
		wantOpcode = bytecode.OpTailCall
	}

	forwarders, cells := 0, 0
	for _, instruction := range program.Bytecode {
		if instruction.Opcode == bytecode.OpMakeCell {
			cells++
		}
	}

	if cells != workload.captures {
		b.Fatalf("created cells = %d, want %d borrowed outer captures", cells, workload.captures)
	}

	for _, function := range functions {
		if function.Name == "leaf" {
			continue
		}

		forwarders++
		if function.Params != workload.captures+1 {
			b.Fatalf("%s parameter/capture count = %d, want %d", function.Name, function.Params, workload.captures+1)
		}

		end := len(program.Bytecode)
		for _, next := range functions {
			if next.Entry > function.Entry && next.Entry < end {
				end = next.Entry
			}
		}

		calls := 0
		for pc := function.Entry; pc < end; pc++ {
			instruction := program.Bytecode[pc]
			if instruction.Opcode == bytecode.OpMakeCell {
				b.Fatalf("%s creates a local cell instead of borrowing its captures", function.Name)
			}

			if instruction.Opcode != bytecode.OpCall && instruction.Opcode != bytecode.OpTailCall {
				continue
			}

			calls++
			if instruction.Opcode != wantOpcode {
				b.Fatalf("%s candidate opcode = %v, want %v", function.Name, instruction.Opcode, wantOpcode)
			}

			for _, catch := range program.CatchTable {
				if catch[0] <= pc && pc <= catch[1] {
					b.Fatalf("%s candidate at %d is inside catch coverage", function.Name, pc)
				}
			}
		}

		if calls != 1 {
			b.Fatalf("%s candidate calls = %d, want 1", function.Name, calls)
		}
	}

	if forwarders != workload.forwarders {
		b.Fatalf("forwarding UDF count = %d, want %d", forwarders, workload.forwarders)
	}
}
