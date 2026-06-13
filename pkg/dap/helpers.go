package dap

import (
	"strings"

	godap "github.com/google/go-dap"

	ferret "github.com/MontFerret/ferret/v2"
)

func splitVariables(variables []ferret.DebugVariable) (locals []ferret.DebugVariable, params []ferret.DebugVariable) {
	for _, variable := range variables {
		if variable.Param {
			params = append(params, variable)
			continue
		}

		locals = append(locals, variable)
	}

	return locals, params
}

func sliceVariables[T any](values []T, start, count int) []T {
	if start < 0 {
		start = 0
	}
	if start >= len(values) {
		return nil
	}
	if count <= 0 || start+count > len(values) {
		count = len(values) - start
	}

	return values[start : start+count]
}

func unsupportedBreakpointMessage(breakpoint godap.SourceBreakpoint) string {
	var messages []string

	if breakpoint.Condition != "" {
		messages = append(messages, "Conditional breakpoints are not supported yet.")
	}
	if breakpoint.HitCondition != "" {
		messages = append(messages, "Hit-count breakpoints are not supported yet.")
	}
	if breakpoint.LogMessage != "" {
		messages = append(messages, "Logpoints are not supported yet.")
	}

	return strings.Join(messages, " ")
}

func breakpointIDs(ids []ferret.DebugBreakpointID) []int {
	if len(ids) == 0 {
		return nil
	}

	out := make([]int, 0, len(ids))
	for _, id := range ids {
		out = append(out, int(id))
	}

	return out
}
