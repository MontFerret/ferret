package internal

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type CallResolver struct {
	session *CompilationSession
}

func NewCallResolver(session *CompilationSession) *CallResolver {
	return &CallResolver{session: session}
}

func (r *CallResolver) ResolveFunctionName(ctx fql.IFunctionCallContext) runtime.String {
	var name string
	funcNS := ctx.Namespace()
	nsText := ""

	if funcNS != nil {
		nsText = funcNS.GetText()
	}

	if nsText != "" {
		ns := nsText

		if len(r.session.Program.UseAliases) > 0 {
			ns = r.applyNamespaceAlias(ns)
		}

		name += ns
		name += ctx.FunctionName().GetText()

		return runtime.NewString(name)
	}

	fn := ctx.FunctionName().GetText()

	if len(r.session.Program.UseAliases) > 0 {
		if target, ok := r.session.Program.UseAliases[fn]; ok && target != "" {
			if strings.Contains(target, runtime.NamespaceSeparator) {
				return runtime.NewString(target)
			}
		}
	}

	name += fn

	return runtime.NewString(name)
}

func (r *CallResolver) ResolveLocalFunctionName(ctx fql.IFunctionCallContext) (string, bool) {
	if ctx == nil {
		return "", false
	}

	if ns := ctx.Namespace(); ns != nil && ns.GetText() != "" {
		return "", false
	}

	fnCtx := ctx.FunctionName()
	if fnCtx == nil {
		return "", false
	}

	name := fnCtx.GetText()
	if name == "" {
		return "", false
	}

	if len(r.session.Program.UseAliases) > 0 {
		if target, ok := r.session.Program.UseAliases[name]; ok && target != "" {
			if strings.Contains(target, runtime.NamespaceSeparator) {
				return "", false
			}
		}
	}

	return name, true
}

func (r *CallResolver) ResolveUDF(ctx fql.IFunctionCallContext) (*core.UDFInfo, bool) {
	if r == nil || r.session == nil || r.session.Program.UDFs == nil || r.session.Function.UDFScope == nil {
		return nil, false
	}

	name, ok := r.ResolveLocalFunctionName(ctx)
	if !ok {
		return nil, false
	}

	return r.session.Program.UDFs.Resolve(name, r.session.Function.UDFScope)
}

func (r *CallResolver) applyNamespaceAlias(ns string) string {
	if ns == "" || len(r.session.Program.UseAliases) == 0 {
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

	target, ok := r.session.Program.UseAliases[parts[0]]
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
