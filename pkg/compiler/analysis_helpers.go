package compiler

import (
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func cloneCall(call Call) Call {
	call.ArgumentSpans = append([]source.Span(nil), call.ArgumentSpans...)

	return call
}

func cloneDiagnostics(values []*diagnostics.Diagnostic) []*diagnostics.Diagnostic {
	out := make([]*diagnostics.Diagnostic, len(values))
	for i, value := range values {
		out[i] = cloneDiagnostic(value)
	}

	return out
}

func cloneDiagnostic(value *diagnostics.Diagnostic) *diagnostics.Diagnostic {
	if value == nil {
		return nil
	}

	out := *value
	out.Spans = append([]diagnostics.ErrorSpan(nil), value.Spans...)

	if value.Cause != nil {
		if diagnostic, ok := value.Cause.(*diagnostics.Diagnostic); ok {
			out.Cause = cloneDiagnostic(diagnostic)
		} else {
			out.Cause = errors.New(value.Cause.Error())
		}
	}

	out.Source = value.Source

	return &out
}

func spanContains(span source.Span, offset int) bool {
	return span.Start >= 0 && span.End > span.Start && offset >= span.Start && offset < span.End
}

func narrowerSpan(candidate, current source.Span) bool {
	if current.End <= current.Start {
		return true
	}

	return candidate.End-candidate.Start < current.End-current.Start
}

func semanticNamespace(kind SymbolKind) string {
	switch kind {
	case SymbolKindUDF:
		return "function"
	case SymbolKindNamespaceAlias:
		return "namespace"
	case SymbolKindBindParameter:
		return "parameter"
	default:
		return "binding"
	}
}

func analysisError(analysis *Analysis) error {
	if analysis == nil {
		return nil
	}

	values := analysis.Diagnostics()
	if len(values) == 0 {
		return nil
	}

	if len(values) == 1 {
		return values[0]
	}

	return diagnostics.NewDiagnosticsOf(values)
}
