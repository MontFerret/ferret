package artifact

import "errors"

var (
	// ErrInvalidArtifact reports a loader or artifact state error not covered by
	// a more specific sentinel.
	ErrInvalidArtifact = errors.New("bytecode artifact: invalid artifact")
	// ErrInvalidHeader reports that the artifact header is malformed or
	// inconsistent with the payload framing.
	ErrInvalidHeader = errors.New("bytecode artifact: invalid header")
	// ErrInvalidPayload reports that the selected payload format could not decode
	// the serialized program.
	ErrInvalidPayload = errors.New("bytecode artifact: invalid payload")
	// ErrInvalidMagic reports that the artifact header does not start with the
	// Ferret artifact magic bytes.
	ErrInvalidMagic = errors.New("bytecode artifact: invalid magic")
	// ErrUnsupportedSchema reports that the artifact schema version is not
	// understood by the current loader.
	ErrUnsupportedSchema = errors.New("bytecode artifact: unsupported schema version")
	// ErrUnknownFormat reports that an artifact references an unknown or
	// unregistered payload format.
	ErrUnknownFormat = errors.New("bytecode artifact: unknown format")
	// ErrIncompatibleISA reports that an artifact or payload requires a different
	// bytecode ISA version than the current runtime.
	ErrIncompatibleISA = errors.New("bytecode artifact: incompatible ISA")
)
