package internal

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func loadConstant(ctx *CompilerContext, value runtime.Value) bytecode.Operand {
	reg := ctx.Registers.Allocate()
	ctx.Emitter.EmitLoadConst(reg, ctx.Symbols.AddConstant(value))
	ctx.Types.Set(reg, valueTypeFromRuntime(value))

	return reg
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

func getFunctionName(ctx fql.IFunctionCallContext, aliases map[string]string) runtime.String {
	var name string
	funcNS := ctx.Namespace()
	nsText := ""

	if funcNS != nil {
		nsText = funcNS.GetText()
	}

	if nsText != "" {
		ns := nsText

		if len(aliases) > 0 {
			ns = applyNamespaceAlias(ns, aliases)
		}

		name += ns
		name += ctx.FunctionName().GetText()

		return runtime.NewString(strings.ToUpper(name))
	}

	fn := ctx.FunctionName().GetText()

	if len(aliases) > 0 {
		if target, ok := aliases[strings.ToUpper(fn)]; ok && target != "" {
			// Bare calls should only use function aliases (e.g. USE NS::FN AS ALIAS).
			// Namespace aliases (e.g. USE NS AS ALIAS) are intended for qualified
			// calls such as ALIAS::FN and must not rewrite ALIAS().
			if strings.Contains(strings.ToUpper(target), runtime.NamespaceSeparator) {
				return runtime.NewString(strings.ToUpper(target))
			}
		}
	}

	name += fn

	return runtime.NewString(strings.ToUpper(name))
}

func applyNamespaceAlias(ns string, aliases map[string]string) string {
	if ns == "" || len(aliases) == 0 {
		return ns
	}

	upper := strings.ToUpper(ns)
	trimmed := strings.TrimSuffix(upper, runtime.NamespaceSeparator)
	if trimmed == "" {
		return upper
	}

	parts := strings.Split(trimmed, runtime.NamespaceSeparator)
	if len(parts) == 0 {
		return upper
	}

	target, ok := aliases[parts[0]]
	if !ok {
		return upper
	}

	target = strings.TrimSuffix(strings.ToUpper(target), runtime.NamespaceSeparator)
	if target == "" {
		return upper
	}

	parts = parts[1:]
	if len(parts) == 0 {
		return target + runtime.NamespaceSeparator
	}

	return target + runtime.NamespaceSeparator + strings.Join(parts, runtime.NamespaceSeparator) + runtime.NamespaceSeparator
}
