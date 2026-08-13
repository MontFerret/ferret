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

func bindingIDFromRule(ctx antlr.ParserRuleContext) core.BindingID {
	if ctx == nil || ctx.GetStart() == nil || ctx.GetStart().GetStart() < 0 {
		return core.InvalidBindingID
	}

	return core.BindingID(ctx.GetStart().GetStart() + 1)
}

func addUDFCapture(captures map[core.BindingID]core.UDFCapture, order *[]core.BindingID, incoming core.UDFCapture) bool {
	if captures == nil || order == nil || incoming.ID == core.InvalidBindingID || incoming.Name == "" {
		return false
	}

	incoming.Mutable = incoming.Storage == core.BindingStorageCell
	capture, exists := captures[incoming.ID]
	if !exists {
		captures[incoming.ID] = incoming
		*order = append(*order, incoming.ID)

		return true
	}

	changed := false
	if incoming.Storage == core.BindingStorageCell && capture.Storage != core.BindingStorageCell {
		capture.Storage = core.BindingStorageCell
		capture.Mutable = true
		changed = true
	}

	if incoming.Visible && !capture.Visible {
		capture.Visible = true
		changed = true
	}

	if changed {
		captures[incoming.ID] = capture
	}

	return changed
}

func orderedUDFCaptures(captures map[core.BindingID]core.UDFCapture, order []core.BindingID) []core.UDFCapture {
	if len(order) == 0 {
		return nil
	}

	out := make([]core.UDFCapture, 0, len(order))
	for _, id := range order {
		capture, ok := captures[id]
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
