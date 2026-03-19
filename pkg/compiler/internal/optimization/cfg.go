package optimization

// ControlFlowGraph represents the control flow structure of a bytecode program
type ControlFlowGraph struct {
	Entry  *BasicBlock   // Entry block (first instruction)
	Exit   *BasicBlock   // Exit block (virtual block representing program exit)
	Blocks []*BasicBlock // All basic blocks in the program
}
