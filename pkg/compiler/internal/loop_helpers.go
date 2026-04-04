package internal

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func aggregateKind(name runtime.String) (bytecode.AggregateKind, bool) {
	fn := strings.ToUpper(name.String())
	if strings.Contains(fn, runtime.NamespaceSeparator) {
		parts := strings.Split(fn, runtime.NamespaceSeparator)
		fn = parts[len(parts)-1]
	}

	switch fn {
	case "COUNT":
		return bytecode.AggregateCount, true
	case "SUM":
		return bytecode.AggregateSum, true
	case "MIN":
		return bytecode.AggregateMin, true
	case "MAX":
		return bytecode.AggregateMax, true
	case "AVERAGE":
		return bytecode.AggregateAverage, true
	default:
		return 0, false
	}
}

func sortDirection(dir antlr.TerminalNode) runtime.SortDirection {
	if dir == nil {
		return runtime.SortDirectionAsc
	}

	if strings.ToLower(dir.GetText()) == "desc" {
		return runtime.SortDirectionDesc
	}

	return runtime.SortDirectionAsc
}
