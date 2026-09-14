package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestCollectAggregateReductionEmission(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		c, err := compiler.New(compiler.WithOptimizationLevel(level))
		if err != nil {
			t.Fatal(err)
		}

		for _, test := range []struct {
			name, query   string
			reduce, hosts int
			protected     bool
		}{
			{"grouped", `RETURN FOR v IN [1,2] COLLECT g=1 AGGREGATE a=SUM(v),b=AVERAGE(v) RETURN [a,b]`, 2, 0, false},
			{"mixed global", `RETURN FOR v IN [1,2] COLLECT AGGREGATE a=SUM(v),b=MEDIAN(v) RETURN [a,b]`, 1, 1, false},
			{"protected", `RETURN FOR v IN [1,2] COLLECT AGGREGATE a=SUM(v)? RETURN a`, 1, 0, true},
			{"namespaced", `RETURN FOR v IN [1,2] COLLECT AGGREGATE a=math::sum(v) RETURN a`, 0, 1, false},
			{"namespaced grouped", `RETURN FOR v IN [1,2] COLLECT g=1 AGGREGATE a=math::sum(v),b=math::min(v),c=math::max(v) RETURN [a,b,c]`, 0, 3, false},
		} {
			t.Run(level.String()+"/"+test.name, func(t *testing.T) {
				program, err := c.Compile(t.Context(), source.New(test.name, test.query))
				if err != nil {
					t.Fatal(err)
				}

				count := 0
				for pc, inst := range program.Bytecode {
					if inst.Opcode != bytecode.OpAggregateReduce {
						continue
					}

					count++
					kind := bytecode.AggregateKind(inst.Operands[2])
					if kind != bytecode.AggregateSum && kind != bytecode.AggregateAverage {
						t.Fatalf("invalid reduction kind after optimization: %d", kind)
					}

					if test.protected {
						caught := false
						for _, region := range program.CatchTable {
							if pc >= region[0] && pc <= region[1] {
								caught = true
								if region[2] <= pc {
									t.Fatal("recovery must skip the failed reduction")
								}
							}
						}

						if !caught {
							t.Fatal("protected reduction lacks recovery region")
						}
					}
				}

				if count != test.reduce || len(program.Functions.Host) != test.hosts {
					t.Fatalf("reductions=%d hosts=%v, want %d/%d", count, program.Functions.Host, test.reduce, test.hosts)
				}

				if err := bytecode.ValidateProgram(program); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}
