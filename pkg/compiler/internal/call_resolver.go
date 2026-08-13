package internal

import (
	"strings"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	CallResolver struct {
		session *CompilationSession
	}

	resolvedCallKind uint8

	resolvedCall struct {
		Function *core.UDFInfo
		Name     string
		Identity string
		Kind     resolvedCallKind
	}
)

const (
	resolvedCallUDF resolvedCallKind = iota + 1
	resolvedCallBuiltin
	resolvedCallHost
)

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
			if r.session.Program.Semantics != nil {
				r.recordNamespaceAliasReference(funcNS, nsText)
			}

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
				if r.session.Program.Semantics != nil && ctx.FunctionName().GetStart() != nil {
					start := ctx.FunctionName().GetStart().GetStart()
					r.session.Program.Semantics.RecordNamespaceAliasReference(fn, source.Span{Start: start, End: start + utf8.RuneCountInString(fn)})
				}

				return runtime.NewString(target)
			}
		}
	}

	name += fn

	return runtime.NewString(name)
}

func (r *CallResolver) recordNamespaceAliasReference(ctx fql.INamespaceContext, namespace string) {
	if r == nil || r.session == nil || r.session.Program.Semantics == nil || ctx == nil || ctx.GetStart() == nil {
		return
	}

	trimmed := strings.TrimSuffix(namespace, runtime.NamespaceSeparator)
	parts := strings.Split(trimmed, runtime.NamespaceSeparator)
	if len(parts) == 0 || parts[0] == "" {
		return
	}

	alias := parts[0]
	if _, ok := r.session.Program.UseAliases[alias]; !ok {
		return
	}

	start := ctx.GetStart().GetStart()
	r.session.Program.Semantics.RecordNamespaceAliasReference(alias, source.Span{Start: start, End: start + utf8.RuneCountInString(alias)})
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

	return r.ResolveUDFInScope(ctx, r.session.Function.UDFScope)
}

func (r *CallResolver) resolveCall(ctx fql.IFunctionCallContext, name runtime.String) resolvedCall {
	nameStr := name.String()
	namespaced := strings.Contains(nameStr, runtime.NamespaceSeparator)

	if ctx != nil {
		if ns := ctx.Namespace(); ns != nil && ns.GetText() != "" {
			namespaced = true
		}
	}

	if !namespaced {
		if fn, ok := r.ResolveUDF(ctx); ok {
			return resolvedCall{
				Function: fn,
				Name:     fn.DisplayName,
				Kind:     resolvedCallUDF,
			}
		}

		builtin := strings.ToUpper(nameStr)
		switch builtin {
		case runtimeLength, runtimeTypename, runtimeWait:
			identity := ""
			if r.session.Program.Semantics != nil {
				identity = runtime.NormalizeRegisteredName(builtin)
			}

			return resolvedCall{
				Name:     builtin,
				Identity: identity,
				Kind:     resolvedCallBuiltin,
			}
		}
	}

	identity := ""
	if r.session.Program.Semantics != nil {
		identity = runtime.NormalizeRegisteredName(nameStr)
	}

	return resolvedCall{
		Name:     nameStr,
		Identity: identity,
		Kind:     resolvedCallHost,
	}
}

// ResolveUDFInScope resolves an unqualified source call against an explicit lexical UDF scope.
func (r *CallResolver) ResolveUDFInScope(ctx fql.IFunctionCallContext, scope *core.UDFScope) (*core.UDFInfo, bool) {
	if r == nil || r.session == nil || r.session.Program.UDFs == nil || scope == nil {
		return nil, false
	}

	name, ok := r.ResolveLocalFunctionName(ctx)
	if !ok {
		return nil, false
	}

	return r.session.Program.UDFs.Resolve(name, scope)
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
