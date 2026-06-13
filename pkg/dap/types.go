package dap

import (
	"bufio"
	"io"

	godap "github.com/google/go-dap"

	ferret "github.com/MontFerret/ferret/v2"
)

type (
	// Config supplies the adapter transport and optional logging sink.
	Config struct {
		In    io.Reader
		Out   io.Writer
		Log   io.Writer
		Trace bool
	}

	launchArguments struct {
		Program     string   `json:"program"`
		Cwd         string   `json:"cwd,omitempty"`
		Args        []string `json:"args,omitempty"`
		StopOnEntry bool     `json:"stopOnEntry,omitempty"`
	}

	requestEnvelope struct {
		godap.ProtocolMessage
		Command string `json:"command"`
	}

	requestMetadata struct {
		Initialize        initializeRequestMetadata
		BreakpointColumns []bool
	}

	initializeRequestMetadata struct {
		LinesStartAt1   *bool
		ColumnsStartAt1 *bool
	}

	breakpointKey struct {
		Line   int
		Column int
	}

	variableHandleKind uint8

	variableHandle struct {
		ScopeVariables []ferret.DebugVariable
		ValueReference ferret.DebugValueReference
		Kind           variableHandleKind
	}

	frameState struct {
		Frame ferret.DebugFrame
		ID    int
	}

	adapterState struct {
		BreakpointsBySource map[string]map[breakpointKey]ferret.DebugBreakpoint
		FrameIndexes        map[int]int
		Handles             map[int]variableHandle
		Frames              []frameState
	}
)

const (
	variableHandleScope variableHandleKind = iota + 1
	variableHandleValue
)

func newReader(input io.Reader) *bufio.Reader {
	if reader, ok := input.(*bufio.Reader); ok {
		return reader
	}

	return bufio.NewReader(input)
}
