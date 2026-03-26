package inspect

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

func HasAggregatePlan(program *bytecode.Program) bool {
	return len(program.Metadata.AggregatePlans) > 0
}
