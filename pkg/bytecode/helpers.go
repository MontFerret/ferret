package bytecode

func IsProtectedCall(op Opcode) bool {
	switch op {
	case OpProtectedHCall, OpProtectedCall:
		return true
	default:
		return false
	}
}

func IsProtectedUdfCall(op Opcode) bool {
	switch op {
	case OpProtectedCall:
		return true
	default:
		return false
	}
}
