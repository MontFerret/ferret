package internal

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func resolveFunctionName(ctx fql.IFunctionCallContext, aliases map[string]string) runtime.String {
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

		return runtime.NewString(name)
	}

	fn := ctx.FunctionName().GetText()

	if len(aliases) > 0 {
		if target, ok := aliases[fn]; ok && target != "" {
			// Bare calls should only use function aliases (e.g. USE NS::FN AS ALIAS).
			// Namespace aliases (e.g. USE NS AS ALIAS) are intended for qualified
			// calls such as ALIAS::FN and must not rewrite ALIAS().
			if strings.Contains(target, runtime.NamespaceSeparator) {
				return runtime.NewString(target)
			}
		}
	}

	name += fn

	return runtime.NewString(name)
}

func applyNamespaceAlias(ns string, aliases map[string]string) string {
	if ns == "" || len(aliases) == 0 {
		return ns
	}

	trimmed := strings.TrimSuffix(ns, runtime.NamespaceSeparator)
	if trimmed == "" {
		return ns
	}

	parts := strings.Split(trimmed, runtime.NamespaceSeparator)
	if len(parts) == 0 {
		return ns
	}

	target, ok := aliases[parts[0]]
	if !ok {
		return ns
	}

	target = strings.TrimSuffix(target, runtime.NamespaceSeparator)
	if target == "" {
		return ns
	}

	parts = parts[1:]
	if len(parts) == 0 {
		return target + runtime.NamespaceSeparator
	}

	return target + runtime.NamespaceSeparator + strings.Join(parts, runtime.NamespaceSeparator) + runtime.NamespaceSeparator
}
