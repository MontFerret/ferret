package asm

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// labelOrAddr returns a label name if one exists for the given address; otherwise just the number.
func labelOrAddr(pos int, labels map[int]string) string {
	if label, ok := labels[pos]; ok {
		return label
	}

	return fmt.Sprintf("%d", pos)
}

// constValue renders the constant at a given index from the program.
func constValue(p *vm.Program, idx int) string {
	if idx >= 0 && idx < len(p.Constants) {
		constant := p.Constants[idx]

		if runtime.IsNumber(constant) {
			return fmt.Sprintf("%d", constant)
		}

		return fmt.Sprintf("%q", constant.String())
	}

	return "<invalid>"
}

// formatLocation returns a line/col comment if available for the given instruction.
func formatLocation(p *vm.Program, ip int) string {
	if ip < len(p.Locations) {
		loc := p.Locations[ip]
		return fmt.Sprintf("; line %d col %d", loc.Line, loc.Column)
	}

	return ""
}

// formatParams generates comments mapping register indices to parameter names.
func formatParams(p *vm.Program) []string {
	lines := []string{}

	for i, name := range p.Params {
		lines = append(lines, fmt.Sprintf("; param R%d = %s", i, name))
	}

	return lines
}
