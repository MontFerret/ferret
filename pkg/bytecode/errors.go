package bytecode

import "errors"

var (
	// ErrInvalidProgram reports that a program is structurally invalid for
	// persistence or validation.
	ErrInvalidProgram = errors.New("bytecode: invalid program")
	// ErrInvalidInstruction reports that a program contains an invalid or
	// unsupported instruction or operand.
	ErrInvalidInstruction = errors.New("bytecode: invalid instruction")
	// ErrInvalidConstant reports that a persisted constant frame is malformed or
	// uses an unsupported type.
	ErrInvalidConstant = errors.New("bytecode: invalid constant")
)
