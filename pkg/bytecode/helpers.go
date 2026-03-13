package bytecode

func IsProtectedCall(op Opcode) bool {
	switch op {
	case OpProtectedHCall, OpProtectedCall:
		return true
	default:
		return false
	}
}
