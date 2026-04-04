package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func variableName(ctx *fql.VariableContext) (string, antlr.Token) {
	if ctx == nil {
		return "", nil
	}

	if id := ctx.Identifier(); id != nil {
		return id.GetText(), id.GetSymbol()
	}

	if id := ctx.SafeReservedWord(); id != nil {
		if prc, ok := id.(antlr.ParserRuleContext); ok {
			return id.GetText(), prc.GetStart()
		}

		return id.GetText(), nil
	}

	return "", nil
}

func findVariableRefs(node antlr.Tree, out *[]*fql.VariableContext) {
	if node == nil || out == nil {
		return
	}

	if v, ok := node.(*fql.VariableContext); ok {
		*out = append(*out, v)
	}

	for i := 0; i < node.GetChildCount(); i++ {
		findVariableRefs(node.GetChild(i), out)
	}
}

func findAssignmentRefs(node antlr.Tree, out *[]*fql.AssignmentStatementContext) {
	if node == nil || out == nil {
		return
	}

	if stmt, ok := node.(*fql.AssignmentStatementContext); ok {
		*out = append(*out, stmt)
	}

	for i := 0; i < node.GetChildCount(); i++ {
		findAssignmentRefs(node.GetChild(i), out)
	}
}

func findFunctionCallRefs(node antlr.Tree, out *[]*fql.FunctionCallContext) {
	if node == nil || out == nil {
		return
	}

	if _, ok := node.(*fql.FunctionDeclarationContext); ok {
		return
	}

	if call, ok := node.(*fql.FunctionCallContext); ok {
		*out = append(*out, call)
	}

	for i := 0; i < node.GetChildCount(); i++ {
		findFunctionCallRefs(node.GetChild(i), out)
	}
}

func addUDFCapture(captures map[string]core.UDFCapture, order *[]string, name string, storage core.BindingStorage) {
	if captures == nil || order == nil || name == "" {
		return
	}

	capture, exists := captures[name]
	if !exists {
		captures[name] = core.UDFCapture{
			Name:    name,
			Mutable: storage == core.BindingStorageCell,
			Storage: storage,
		}
		*order = append(*order, name)
		return
	}

	if storage == core.BindingStorageCell && capture.Storage != core.BindingStorageCell {
		capture.Storage = core.BindingStorageCell
		capture.Mutable = true
		captures[name] = capture
	}
}

func orderedUDFCaptures(captures map[string]core.UDFCapture, order []string) []core.UDFCapture {
	if len(order) == 0 {
		return nil
	}

	out := make([]core.UDFCapture, 0, len(order))
	for _, name := range order {
		capture, ok := captures[name]
		if !ok {
			continue
		}
		out = append(out, capture)
	}

	return out
}

func udfSetToSlice(set map[*core.UDFInfo]struct{}) []*core.UDFInfo {
	if len(set) == 0 {
		return nil
	}

	out := make([]*core.UDFInfo, 0, len(set))
	for fn := range set {
		out = append(out, fn)
	}

	return out
}
