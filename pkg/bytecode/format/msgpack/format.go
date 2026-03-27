package msgpack

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	vmmsgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/internal/persist"
)

type Format struct{}

// Default is the built-in MessagePack bytecode payload format.
var Default = Format{}

func (Format) Name() string {
	return "msgpack"
}

func (Format) Marshal(program *bytecode.Program) ([]byte, error) {
	frame, err := persist.FromProgram(program)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	enc := vmmsgpack.GetEncoder()
	enc.ResetWriter(&buf)
	defer vmmsgpack.PutEncoder(enc)

	if err := enc.Encode(frame); err != nil {
		return nil, fmt.Errorf("bytecode msgpack format: marshal payload: %w", err)
	}

	return buf.Bytes(), nil
}

func (Format) Unmarshal(data []byte) (*bytecode.Program, error) {
	dec := vmmsgpack.GetDecoder()
	dec.ResetReader(bytes.NewReader(data))
	defer vmmsgpack.PutDecoder(dec)

	var frame persist.ProgramFrame
	if err := dec.Decode(&frame); err != nil {
		return nil, fmt.Errorf("bytecode msgpack format: unmarshal payload: %w", err)
	}

	if _, err := dec.PeekCode(); err == nil {
		return nil, fmt.Errorf("bytecode msgpack format: multiple root values")
	} else if !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("bytecode msgpack format: check payload tail: %w", err)
	}

	return persist.ToProgram(frame)
}
