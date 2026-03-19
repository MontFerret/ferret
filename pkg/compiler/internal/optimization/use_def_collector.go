package optimization

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

type useDefCollector struct {
	uses []int
	defs []int
}

func (c *useDefCollector) addUse(op bytecode.Operand) {
	if op != bytecode.NoopOperand && op.IsRegister() {
		c.uses = append(c.uses, op.Register())
	}
}

func (c *useDefCollector) addDef(op bytecode.Operand) {
	if op != bytecode.NoopOperand && op.IsRegister() {
		c.defs = append(c.defs, op.Register())
	}
}

func (c *useDefCollector) addRangeUses(start, count int) {
	if count <= 0 || start <= 0 {
		return
	}

	for reg := start; reg < start+count; reg++ {
		c.uses = append(c.uses, reg)
	}
}
