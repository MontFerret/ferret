package dap

import (
	"bufio"
	"encoding/json"
	"fmt"

	godap "github.com/google/go-dap"
)

func readRequest(reader *bufio.Reader) (godap.Message, requestMetadata, error) {
	content, err := godap.ReadBaseMessage(reader)
	if err != nil {
		return nil, requestMetadata{}, err
	}

	message, err := decodeRequest(content)
	if err != nil {
		return nil, requestMetadata{}, err
	}

	metadata, err := decodeRequestMetadata(content)
	if err != nil {
		return nil, requestMetadata{}, err
	}

	return message, metadata, nil
}

func decodeRequest(content []byte) (godap.Message, error) {
	var envelope requestEnvelope
	if err := json.Unmarshal(content, &envelope); err != nil {
		return nil, err
	}

	if envelope.Type != "request" {
		return nil, fmt.Errorf("unexpected protocol message type %q", envelope.Type)
	}

	switch envelope.Command {
	case "initialize":
		return unmarshalRequest(content, &godap.InitializeRequest{})
	case "launch":
		return unmarshalRequest(content, &godap.LaunchRequest{})
	case "setBreakpoints":
		return unmarshalRequest(content, &godap.SetBreakpointsRequest{})
	case "configurationDone":
		return unmarshalRequest(content, &godap.ConfigurationDoneRequest{})
	case "threads":
		return unmarshalRequest(content, &godap.ThreadsRequest{})
	case "stackTrace":
		return unmarshalRequest(content, &godap.StackTraceRequest{})
	case "scopes":
		return unmarshalRequest(content, &godap.ScopesRequest{})
	case "variables":
		return unmarshalRequest(content, &godap.VariablesRequest{})
	case "continue":
		return unmarshalRequest(content, &godap.ContinueRequest{})
	case "next":
		return unmarshalRequest(content, &godap.NextRequest{})
	case "stepIn":
		return unmarshalRequest(content, &godap.StepInRequest{})
	case "stepOut":
		return unmarshalRequest(content, &godap.StepOutRequest{})
	case "pause":
		return unmarshalRequest(content, &godap.PauseRequest{})
	case "evaluate":
		return unmarshalRequest(content, &godap.EvaluateRequest{})
	case "disconnect":
		return unmarshalRequest(content, &godap.DisconnectRequest{})
	case "setExceptionBreakpoints":
		return unmarshalRequest(content, &godap.SetExceptionBreakpointsRequest{})
	default:
		return unmarshalRequest(content, &godap.Request{})
	}
}

func unmarshalRequest(content []byte, message godap.Message) (godap.Message, error) {
	if err := json.Unmarshal(content, message); err != nil {
		return nil, err
	}

	return message, nil
}

func decodeRequestMetadata(content []byte) (requestMetadata, error) {
	var raw struct {
		Arguments struct {
			LinesStartAt1   *bool `json:"linesStartAt1"`
			ColumnsStartAt1 *bool `json:"columnsStartAt1"`
			Breakpoints     []struct {
				Column *int `json:"column"`
			} `json:"breakpoints"`
		} `json:"arguments"`
	}
	if err := json.Unmarshal(content, &raw); err != nil {
		return requestMetadata{}, err
	}

	metadata := requestMetadata{
		Initialize: initializeRequestMetadata{
			LinesStartAt1:   raw.Arguments.LinesStartAt1,
			ColumnsStartAt1: raw.Arguments.ColumnsStartAt1,
		},
		BreakpointColumns: make([]bool, len(raw.Arguments.Breakpoints)),
	}
	for index, breakpoint := range raw.Arguments.Breakpoints {
		metadata.BreakpointColumns[index] = breakpoint.Column != nil
	}

	return metadata, nil
}
