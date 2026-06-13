package dap

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	godap "github.com/google/go-dap"
)

const (
	testTimeout = 5 * time.Second
)

type adapterHarness struct {
	adapter *Adapter
	errCh   chan error
	in      *io.PipeWriter
	out     *bufio.Reader
}

func TestAdapterGoldenFlow(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "examples/debug/dap.fql", `LET items = [1, 2, 3]
FUNC double(value) (
  RETURN value * 2
)
LET mapped = (
  FOR item IN items
  RETURN double(item)
)
RETURN mapped`)

	harness := newAdapterHarness(t)
	defer harness.close(t)

	harness.send(t, initializeRequest(1, "path"))
	if _, ok := harness.read(t).(*godap.InitializeResponse); !ok {
		t.Fatal("expected initialize response")
	}

	harness.send(t, launchRequest(2, launchArguments{
		Program:     programPath,
		Cwd:         dir,
		StopOnEntry: true,
	}))
	if _, ok := harness.read(t).(*godap.LaunchResponse); !ok {
		t.Fatal("expected launch response")
	}
	if _, ok := harness.read(t).(*godap.InitializedEvent); !ok {
		t.Fatal("expected initialized event")
	}

	harness.send(t, setBreakpointsRequest(3, programPath, []godap.SourceBreakpoint{{Line: 5}}))
	breakpointsResponse, ok := harness.read(t).(*godap.SetBreakpointsResponse)
	if !ok {
		t.Fatal("expected setBreakpoints response")
	}
	if len(breakpointsResponse.Body.Breakpoints) != 1 || !breakpointsResponse.Body.Breakpoints[0].Verified || breakpointsResponse.Body.Breakpoints[0].Line != 5 {
		t.Fatalf("unexpected breakpoint response: %#v", breakpointsResponse.Body.Breakpoints)
	}

	harness.send(t, setExceptionBreakpointsRequest(4))
	if _, ok := harness.read(t).(*godap.SetExceptionBreakpointsResponse); !ok {
		t.Fatal("expected setExceptionBreakpoints response")
	}

	harness.send(t, configurationDoneRequest(5))
	if _, ok := harness.read(t).(*godap.ConfigurationDoneResponse); !ok {
		t.Fatal("expected configurationDone response")
	}
	entryStop, ok := harness.read(t).(*godap.StoppedEvent)
	if !ok || entryStop.Body.Reason != "entry" {
		t.Fatalf("unexpected entry stop: %#v", entryStop)
	}

	harness.send(t, continueRequest(6))
	if _, ok := harness.read(t).(*godap.ContinueResponse); !ok {
		t.Fatal("expected continue response")
	}
	breakpointStop, ok := harness.read(t).(*godap.StoppedEvent)
	if !ok || breakpointStop.Body.Reason != "breakpoint" || len(breakpointStop.Body.HitBreakpointIds) != 1 {
		t.Fatalf("unexpected breakpoint stop: %#v", breakpointStop)
	}

	harness.send(t, stackTraceRequest(7, 1))
	stackTrace, ok := harness.read(t).(*godap.StackTraceResponse)
	if !ok || len(stackTrace.Body.StackFrames) == 0 {
		t.Fatalf("unexpected stack trace: %#v", stackTrace)
	}
	if stackTrace.Body.StackFrames[0].Line != 5 {
		t.Fatalf("expected breakpoint frame on line 5, got %#v", stackTrace.Body.StackFrames[0])
	}

	harness.send(t, scopesRequest(8, stackTrace.Body.StackFrames[0].Id))
	scopes, ok := harness.read(t).(*godap.ScopesResponse)
	if !ok || len(scopes.Body.Scopes) == 0 {
		t.Fatalf("unexpected scopes response: %#v", scopes)
	}

	harness.send(t, variablesRequest(9, scopes.Body.Scopes[0].VariablesReference))
	locals, ok := harness.read(t).(*godap.VariablesResponse)
	if !ok || len(locals.Body.Variables) == 0 {
		t.Fatalf("unexpected locals response: %#v", locals)
	}

	var itemsReference int
	for _, variable := range locals.Body.Variables {
		if variable.Name == "items" {
			itemsReference = variable.VariablesReference
			break
		}
	}
	if itemsReference == 0 {
		t.Fatalf("expected expandable items local, got %#v", locals.Body.Variables)
	}

	harness.send(t, variablesRequest(10, itemsReference))
	items, ok := harness.read(t).(*godap.VariablesResponse)
	if !ok || len(items.Body.Variables) != 3 {
		t.Fatalf("unexpected items response: %#v", items)
	}

	harness.send(t, continueRequest(11))
	if _, ok := harness.read(t).(*godap.ContinueResponse); !ok {
		t.Fatal("expected second continue response")
	}
	if terminated, ok := harness.read(t).(*godap.TerminatedEvent); !ok || terminated.Event.Event != "terminated" {
		t.Fatalf("unexpected termination event: %#v", terminated)
	}
}

func TestAdapterSetBreakpointsNormalizesPathsAndBindsWithinFunctions(t *testing.T) {
	dir := t.TempDir()
	relativeProgram := filepath.Join("queries", "bind.fql")
	programPath := writeScript(t, dir, relativeProgram, `LET seed = 1
FUNC add(a) (
  LET b = a + 1

  RETURN b
)

RETURN add(seed)`)

	harness := newAdapterHarness(t)
	defer harness.close(t)

	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)

	harness.send(t, launchRequest(2, launchArguments{
		Program: relativeProgram,
		Cwd:     dir,
	}))
	harness.read(t)
	harness.read(t)

	harness.send(t, setBreakpointsRequest(3, filepath.Base(programPath), []godap.SourceBreakpoint{{Line: 4}}))
	response, ok := harness.read(t).(*godap.SetBreakpointsResponse)
	if !ok {
		t.Fatal("expected setBreakpoints response")
	}
	if len(response.Body.Breakpoints) != 1 || !response.Body.Breakpoints[0].Verified || response.Body.Breakpoints[0].Line != 5 {
		t.Fatalf("unexpected normalized breakpoint response: %#v", response.Body.Breakpoints)
	}
}

func TestAdapterStepInAndStepOutExposeUDFFrames(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "udf.fql", `FUNC add(a) (
  LET b = a + 1
  RETURN b
)
LET x = add(2)
RETURN x`)

	harness := newAdapterHarness(t)
	defer harness.close(t)

	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)

	harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir, StopOnEntry: true}))
	harness.read(t)
	harness.read(t)

	harness.send(t, configurationDoneRequest(3))
	harness.read(t)
	harness.read(t)

	harness.send(t, stepInRequest(4))
	if _, ok := harness.read(t).(*godap.StepInResponse); !ok {
		t.Fatal("expected stepIn response")
	}
	stepStop, ok := harness.read(t).(*godap.StoppedEvent)
	if !ok || stepStop.Body.Reason != "step" {
		t.Fatalf("unexpected step stop: %#v", stepStop)
	}

	harness.send(t, stackTraceRequest(5, 1))
	stackTrace, ok := harness.read(t).(*godap.StackTraceResponse)
	if !ok || len(stackTrace.Body.StackFrames) != 2 {
		t.Fatalf("unexpected stack trace: %#v", stackTrace)
	}
	if stackTrace.Body.StackFrames[0].Name != "add" || stackTrace.Body.StackFrames[1].Name != "<main>" {
		t.Fatalf("unexpected UDF frames: %#v", stackTrace.Body.StackFrames)
	}

	harness.send(t, stepOutRequest(6))
	if _, ok := harness.read(t).(*godap.StepOutResponse); !ok {
		t.Fatal("expected stepOut response")
	}
	if stop, ok := harness.read(t).(*godap.StoppedEvent); !ok || stop.Body.Reason != "step" {
		t.Fatalf("unexpected stepOut stop: %#v", stop)
	}
}

func TestAdapterNextStepsOverCalls(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "next.fql", `FUNC add(a) (
  RETURN a + 1
)
LET x = add(2)
RETURN x`)

	harness := newAdapterHarness(t)
	defer harness.close(t)

	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)

	harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir, StopOnEntry: true}))
	harness.read(t)
	harness.read(t)

	harness.send(t, configurationDoneRequest(3))
	harness.read(t)
	harness.read(t)

	harness.send(t, nextRequest(4))
	if _, ok := harness.read(t).(*godap.NextResponse); !ok {
		t.Fatal("expected next response")
	}
	if stop, ok := harness.read(t).(*godap.StoppedEvent); !ok || stop.Body.Reason != "step" {
		t.Fatalf("unexpected next stop: %#v", stop)
	}

	harness.send(t, stackTraceRequest(5, 1))
	stackTrace, ok := harness.read(t).(*godap.StackTraceResponse)
	if !ok || len(stackTrace.Body.StackFrames) != 1 || stackTrace.Body.StackFrames[0].Line != 5 {
		t.Fatalf("unexpected next stack trace: %#v", stackTrace)
	}
}

func TestAdapterRuntimeErrorStopsAsException(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "error.fql", `LET x = 7
RETURN x / 0`)

	harness := newAdapterHarness(t)
	defer harness.close(t)

	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)

	harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir}))
	harness.read(t)
	harness.read(t)

	harness.send(t, configurationDoneRequest(3))
	if _, ok := harness.read(t).(*godap.ConfigurationDoneResponse); !ok {
		t.Fatal("expected configurationDone response")
	}
	stop, ok := harness.read(t).(*godap.StoppedEvent)
	if !ok || stop.Body.Reason != "exception" {
		t.Fatalf("unexpected exception stop: %#v", stop)
	}

	harness.send(t, continueRequest(4))
	harness.read(t)
	if _, ok := harness.read(t).(*godap.TerminatedEvent); !ok {
		t.Fatal("expected terminated event after runtime error continue")
	}
}

func TestAdapterInvalidEvaluateAndUnknownRequestReturnProtocolErrors(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "eval.fql", `LET x = 1
RETURN x`)

	harness := newAdapterHarness(t)
	defer harness.close(t)

	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)

	harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir, StopOnEntry: true}))
	harness.read(t)
	harness.read(t)

	harness.send(t, configurationDoneRequest(3))
	harness.read(t)
	harness.read(t)

	harness.send(t, stackTraceRequest(4, 1))
	stackTrace, ok := harness.read(t).(*godap.StackTraceResponse)
	if !ok || len(stackTrace.Body.StackFrames) == 0 {
		t.Fatalf("unexpected stack trace: %#v", stackTrace)
	}

	harness.send(t, evaluateRequest(5, "LENGTH([1])", stackTrace.Body.StackFrames[0].Id))
	if response, ok := harness.read(t).(*godap.ErrorResponse); !ok || response.Command != "evaluate" {
		t.Fatalf("expected evaluate error response, got %#v", response)
	}

	harness.send(t, customRequest(6, "unknownRequest"))
	if response, ok := harness.read(t).(*godap.ErrorResponse); !ok || response.Command != "unknownRequest" {
		t.Fatalf("expected unknown request error response, got %#v", response)
	}
}

func TestAdapterPauseAndDisconnectCleanup(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "pause.fql", `LET values = (
  FOR item IN 1..500000000
  RETURN item
)
RETURN values`)

	harness := newAdapterHarness(t)

	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)

	harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir}))
	harness.read(t)
	harness.read(t)

	harness.send(t, configurationDoneRequest(3))
	if _, ok := harness.read(t).(*godap.ConfigurationDoneResponse); !ok {
		t.Fatal("expected configurationDone response")
	}

	harness.send(t, pauseRequest(4))
	if _, ok := harness.read(t).(*godap.PauseResponse); !ok {
		t.Fatal("expected pause response")
	}
	if stop, ok := harness.read(t).(*godap.StoppedEvent); !ok || stop.Body.Reason != "pause" {
		t.Fatalf("unexpected pause stop: %#v", stop)
	}

	harness.send(t, disconnectRequest(5))
	if _, ok := harness.read(t).(*godap.DisconnectResponse); !ok {
		t.Fatal("expected disconnect response")
	}
	if _, ok := harness.read(t).(*godap.TerminatedEvent); !ok {
		t.Fatal("expected terminated event on disconnect")
	}

	harness.wait(t)
}

func TestAdapterRejectsStaleFramesAndVariablesAcrossStops(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "stale.fql", `LET items = [1]
LET value = 2
RETURN items`)

	harness := newAdapterHarness(t)
	defer harness.close(t)
	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)
	harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir, StopOnEntry: true}))
	harness.read(t)
	harness.read(t)
	harness.send(t, setBreakpointsRequest(3, programPath, []godap.SourceBreakpoint{{Line: 2}}))
	harness.read(t)
	harness.send(t, configurationDoneRequest(4))
	harness.read(t)
	harness.read(t)

	harness.send(t, stackTraceRequest(5, 1))
	firstStack := harness.read(t).(*godap.StackTraceResponse)
	firstFrameID := firstStack.Body.StackFrames[0].Id
	harness.send(t, scopesRequest(6, firstFrameID))
	firstScopes := harness.read(t).(*godap.ScopesResponse)
	firstVariablesReference := firstScopes.Body.Scopes[0].VariablesReference

	harness.send(t, continueRequest(7))
	harness.read(t)
	harness.read(t)
	harness.send(t, stackTraceRequest(8, 1))
	secondStack := harness.read(t).(*godap.StackTraceResponse)
	if secondStack.Body.StackFrames[0].Id == firstFrameID {
		t.Fatalf("frame ID was reused across stops: %d", firstFrameID)
	}

	harness.send(t, scopesRequest(9, firstFrameID))
	if response, ok := harness.read(t).(*godap.ErrorResponse); !ok || response.Command != "scopes" {
		t.Fatalf("expected stale frame error, got %#v", response)
	}
	harness.send(t, variablesRequest(10, firstVariablesReference))
	if response, ok := harness.read(t).(*godap.ErrorResponse); !ok || response.Command != "variables" {
		t.Fatalf("expected stale variables error, got %#v", response)
	}
}

func TestAdapterBreakpointsReuseDuplicatesAndReportUnsupportedFields(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "breakpoints.fql", `LET value = 1
RETURN value`)
	sourceURI := (&url.URL{Scheme: "file", Path: filepath.ToSlash(programPath)}).String()

	harness := newAdapterHarness(t)
	defer harness.close(t)
	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)
	harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir}))
	harness.read(t)
	harness.read(t)

	harness.send(t, setBreakpointsRequest(3, sourceURI, []godap.SourceBreakpoint{{Line: 2}, {Line: 2}}))
	duplicates := harness.read(t).(*godap.SetBreakpointsResponse)
	if len(duplicates.Body.Breakpoints) != 2 ||
		duplicates.Body.Breakpoints[0].Id == 0 ||
		duplicates.Body.Breakpoints[0].Id != duplicates.Body.Breakpoints[1].Id {
		t.Fatalf("duplicate breakpoints were not reused: %#v", duplicates.Body.Breakpoints)
	}

	harness.send(t, setBreakpointsRequest(4, programPath, []godap.SourceBreakpoint{{
		Line:         2,
		Condition:    "value == 1",
		HitCondition: "2",
		LogMessage:   "value",
	}}))
	unsupported := harness.read(t).(*godap.SetBreakpointsResponse)
	const wantMessage = "Conditional breakpoints are not supported yet. Hit-count breakpoints are not supported yet. Logpoints are not supported yet."
	if len(unsupported.Body.Breakpoints) != 1 ||
		unsupported.Body.Breakpoints[0].Verified ||
		unsupported.Body.Breakpoints[0].Message != wantMessage {
		t.Fatalf("unexpected unsupported breakpoint response: %#v", unsupported.Body.Breakpoints)
	}
}

func TestAdapterHonorsDefaultAndExplicitZeroBasedPositions(t *testing.T) {
	t.Run("omitted_base_settings_default_to_one", func(t *testing.T) {
		dir := t.TempDir()
		programPath := writeScript(t, dir, "default-base.fql", "RETURN 1")
		harness := newAdapterHarness(t)
		defer harness.close(t)

		harness.sendJSON(t, map[string]any{
			"seq": 1, "type": "request", "command": "initialize",
			"arguments": map[string]any{"adapterID": "ferret"},
		})
		harness.read(t)
		harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir}))
		harness.read(t)
		harness.read(t)
		harness.send(t, setBreakpointsRequest(3, programPath, []godap.SourceBreakpoint{{Line: 1}}))
		breakpoints := harness.read(t).(*godap.SetBreakpointsResponse)
		if !breakpoints.Body.Breakpoints[0].Verified || breakpoints.Body.Breakpoints[0].Line != 1 {
			t.Fatalf("default base settings were not one-based: %#v", breakpoints.Body.Breakpoints)
		}
		harness.send(t, configurationDoneRequest(4))
		harness.read(t)
		if stop := harness.read(t).(*godap.StoppedEvent); stop.Body.Reason != "breakpoint" {
			t.Fatalf("expected entry breakpoint with default bases, got %#v", stop)
		}
	})

	t.Run("explicit_zero_based_column_is_present", func(t *testing.T) {
		dir := t.TempDir()
		programPath := writeScript(t, dir, "zero-base.fql", "RETURN 1")
		harness := newAdapterHarness(t)
		defer harness.close(t)

		harness.sendJSON(t, map[string]any{
			"seq": 1, "type": "request", "command": "initialize",
			"arguments": map[string]any{
				"adapterID": "ferret", "linesStartAt1": false, "columnsStartAt1": false,
			},
		})
		harness.read(t)
		harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir}))
		harness.read(t)
		harness.read(t)
		harness.sendJSON(t, map[string]any{
			"seq": 3, "type": "request", "command": "setBreakpoints",
			"arguments": map[string]any{
				"source":      map[string]any{"path": programPath},
				"breakpoints": []map[string]any{{"line": 0, "column": 0}},
			},
		})
		breakpoints := harness.read(t).(*godap.SetBreakpointsResponse)
		if !breakpoints.Body.Breakpoints[0].Verified {
			t.Fatalf("zero-based breakpoint was not verified: %#v", breakpoints.Body.Breakpoints)
		}
		harness.adapter.mu.Lock()
		_, hasExplicitFirstColumn := harness.adapter.state.BreakpointsBySource[programPath][breakpointKey{Line: 1, Column: 1}]
		harness.adapter.mu.Unlock()
		if !hasExplicitFirstColumn {
			t.Fatal("explicit zero-based column was treated as omitted")
		}
	})
}

func TestAdapterRejectsUnsupportedEvaluateContextAndInspectionWhileRunning(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "running.fql", `LET values = (
  FOR item IN 1..500000000
  RETURN item
)
RETURN values`)
	harness := newAdapterHarness(t)

	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)
	harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir, StopOnEntry: true}))
	harness.read(t)
	harness.read(t)
	harness.send(t, configurationDoneRequest(3))
	harness.read(t)
	harness.read(t)

	unsupported := evaluateRequest(4, "1", 0)
	unsupported.Arguments.Context = "clipboard"
	harness.send(t, unsupported)
	if response, ok := harness.read(t).(*godap.ErrorResponse); !ok || response.Command != "evaluate" {
		t.Fatalf("expected unsupported evaluate context error, got %#v", response)
	}

	harness.send(t, continueRequest(5))
	harness.read(t)
	harness.send(t, stackTraceRequest(6, 1))
	if response, ok := harness.read(t).(*godap.ErrorResponse); !ok || response.Command != "stackTrace" {
		t.Fatalf("expected running stackTrace error, got %#v", response)
	}
	harness.send(t, variablesRequest(7, 1))
	if response, ok := harness.read(t).(*godap.ErrorResponse); !ok || response.Command != "variables" {
		t.Fatalf("expected running variables error, got %#v", response)
	}

	harness.send(t, disconnectRequest(8))
	harness.read(t)
	harness.read(t)
	harness.wait(t)
	if _, err := godap.ReadProtocolMessage(harness.out); !errors.Is(err, io.EOF) {
		t.Fatalf("expected no late stopped event after disconnect, got %v", err)
	}
}

func TestAdapterRejectsSecondLaunch(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "launch.fql", "RETURN 1")
	harness := newAdapterHarness(t)
	defer harness.close(t)

	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)
	harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir}))
	harness.read(t)
	harness.read(t)
	harness.send(t, launchRequest(3, launchArguments{Program: programPath, Cwd: dir}))
	if response, ok := harness.read(t).(*godap.ErrorResponse); !ok || response.Command != "launch" {
		t.Fatalf("expected second launch error, got %#v", response)
	}
}

func TestAdapterWriteFailureClosesLaunchedResources(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "write-error.fql", "RETURN 1")
	var input bytes.Buffer
	if err := godap.WriteProtocolMessage(&input, initializeRequest(1, "path")); err != nil {
		t.Fatal(err)
	}
	if err := godap.WriteProtocolMessage(&input, launchRequest(2, launchArguments{Program: programPath, Cwd: dir})); err != nil {
		t.Fatal(err)
	}

	writeErr := errors.New("write failed")
	output := &failAfterWriter{FailAfter: 4, Err: writeErr}
	adapter, err := New(Config{In: &input, Out: output})
	if err != nil {
		t.Fatal(err)
	}
	if err := adapter.Serve(context.Background()); !errors.Is(err, writeErr) {
		t.Fatalf("expected write failure, got %v", err)
	}

	adapter.mu.Lock()
	defer adapter.mu.Unlock()
	if !adapter.adapterClosed || !adapter.resourcesClosed || adapter.session != nil || adapter.plan != nil || adapter.engine != nil {
		t.Fatalf("write failure did not close launched resources: %#v", adapter)
	}
}

func TestAdapterTerminalProtocolFailureTerminatesAndCleansUp(t *testing.T) {
	dir := t.TempDir()
	programPath := writeScript(t, dir, "protocol-error.fql", "RETURN 1")
	harness := newAdapterHarness(t)

	harness.send(t, initializeRequest(1, "path"))
	harness.read(t)
	harness.send(t, launchRequest(2, launchArguments{Program: programPath, Cwd: dir}))
	harness.read(t)
	harness.read(t)
	harness.sendRaw(t, []byte("{"))
	if _, ok := harness.read(t).(*godap.TerminatedEvent); !ok {
		t.Fatal("expected terminated event after terminal protocol failure")
	}
	if err := harness.waitResult(t); err == nil {
		t.Fatal("expected terminal protocol error")
	}
	if _, err := godap.ReadProtocolMessage(harness.out); !errors.Is(err, io.EOF) {
		t.Fatalf("expected exactly one terminated event, got %v", err)
	}

	harness.adapter.mu.Lock()
	defer harness.adapter.mu.Unlock()
	if !harness.adapter.adapterClosed || !harness.adapter.resourcesClosed {
		t.Fatalf("terminal protocol failure did not clean up adapter: %#v", harness.adapter)
	}
}

func TestAdapterContextCancellationInterruptsInputRead(t *testing.T) {
	input, writer := io.Pipe()
	adapter, err := New(Config{In: input, Out: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- adapter.Serve(ctx)
	}()

	cancel()
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(testTimeout):
		t.Fatal("context cancellation did not interrupt input read")
	}
	if _, err := writer.Write([]byte("x")); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("expected input transport to close, got %v", err)
	}
}

func newAdapterHarness(t *testing.T) *adapterHarness {
	t.Helper()

	inReader, inWriter := io.Pipe()
	outReader, outWriter := io.Pipe()
	adapter, err := New(Config{
		In:  inReader,
		Out: outWriter,
		Log: io.Discard,
	})
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)
	go func() {
		err := adapter.Serve(context.Background())
		_ = outWriter.Close()
		errCh <- err
	}()

	return &adapterHarness{
		adapter: adapter,
		errCh:   errCh,
		in:      inWriter,
		out:     bufio.NewReader(outReader),
	}
}

func (h *adapterHarness) sendJSON(t *testing.T, message any) {
	t.Helper()

	data, err := json.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}
	if err := godap.WriteBaseMessage(h.in, data); err != nil {
		t.Fatal(err)
	}
}

func (h *adapterHarness) sendRaw(t *testing.T, data []byte) {
	t.Helper()

	if err := godap.WriteBaseMessage(h.in, data); err != nil {
		t.Fatal(err)
	}
}

func (h *adapterHarness) send(t *testing.T, message godap.Message) {
	t.Helper()

	if err := godap.WriteProtocolMessage(h.in, message); err != nil {
		t.Fatal(err)
	}
}

func (h *adapterHarness) read(t *testing.T) godap.Message {
	t.Helper()

	type result struct {
		message godap.Message
		err     error
	}

	resultCh := make(chan result, 1)
	go func() {
		message, err := godap.ReadProtocolMessage(h.out)
		resultCh <- result{message: message, err: err}
	}()

	select {
	case result := <-resultCh:
		if result.err != nil {
			t.Fatal(result.err)
		}

		return result.message
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for DAP message")
		return nil
	}
}

func (h *adapterHarness) close(t *testing.T) {
	t.Helper()

	_ = h.in.Close()
	h.wait(t)
}

func (h *adapterHarness) wait(t *testing.T) {
	t.Helper()

	if err := h.waitResult(t); err != nil {
		t.Fatal(err)
	}
}

func (h *adapterHarness) waitResult(t *testing.T) error {
	t.Helper()

	select {
	case err := <-h.errCh:
		return err
	case <-time.After(testTimeout):
		t.Fatal("timed out waiting for adapter shutdown")
		return nil
	}
}

func writeScript(t *testing.T, dir, relativePath, content string) string {
	t.Helper()

	path := filepath.Join(dir, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	return path
}

func initializeRequest(seq int, pathFormat string) *godap.InitializeRequest {
	return &godap.InitializeRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "initialize",
		},
		Arguments: godap.InitializeRequestArguments{
			AdapterID:       "ferret",
			LinesStartAt1:   true,
			ColumnsStartAt1: true,
			PathFormat:      pathFormat,
		},
	}
}

func launchRequest(seq int, args launchArguments) *godap.LaunchRequest {
	raw, _ := json.Marshal(args)
	return &godap.LaunchRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "launch",
		},
		Arguments: raw,
	}
}

func setBreakpointsRequest(seq int, path string, breakpoints []godap.SourceBreakpoint) *godap.SetBreakpointsRequest {
	return &godap.SetBreakpointsRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "setBreakpoints",
		},
		Arguments: godap.SetBreakpointsArguments{
			Source:      godap.Source{Path: path},
			Breakpoints: breakpoints,
		},
	}
}

func setExceptionBreakpointsRequest(seq int) *godap.SetExceptionBreakpointsRequest {
	return &godap.SetExceptionBreakpointsRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "setExceptionBreakpoints",
		},
		Arguments: godap.SetExceptionBreakpointsArguments{},
	}
}

func configurationDoneRequest(seq int) *godap.ConfigurationDoneRequest {
	return &godap.ConfigurationDoneRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "configurationDone",
		},
	}
}

func continueRequest(seq int) *godap.ContinueRequest {
	return &godap.ContinueRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "continue",
		},
		Arguments: godap.ContinueArguments{ThreadId: 1},
	}
}

func nextRequest(seq int) *godap.NextRequest {
	return &godap.NextRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "next",
		},
		Arguments: godap.NextArguments{ThreadId: 1},
	}
}

func stepInRequest(seq int) *godap.StepInRequest {
	return &godap.StepInRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "stepIn",
		},
		Arguments: godap.StepInArguments{ThreadId: 1},
	}
}

func stepOutRequest(seq int) *godap.StepOutRequest {
	return &godap.StepOutRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "stepOut",
		},
		Arguments: godap.StepOutArguments{ThreadId: 1},
	}
}

func pauseRequest(seq int) *godap.PauseRequest {
	return &godap.PauseRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "pause",
		},
		Arguments: godap.PauseArguments{ThreadId: 1},
	}
}

func stackTraceRequest(seq, threadID int) *godap.StackTraceRequest {
	return &godap.StackTraceRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "stackTrace",
		},
		Arguments: godap.StackTraceArguments{ThreadId: threadID},
	}
}

func scopesRequest(seq, frameID int) *godap.ScopesRequest {
	return &godap.ScopesRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "scopes",
		},
		Arguments: godap.ScopesArguments{FrameId: frameID},
	}
}

func variablesRequest(seq, reference int) *godap.VariablesRequest {
	return &godap.VariablesRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "variables",
		},
		Arguments: godap.VariablesArguments{VariablesReference: reference},
	}
}

func evaluateRequest(seq int, expression string, frameID int) *godap.EvaluateRequest {
	return &godap.EvaluateRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "evaluate",
		},
		Arguments: godap.EvaluateArguments{
			Expression: expression,
			FrameId:    frameID,
			Context:    "watch",
		},
	}
}

func disconnectRequest(seq int) *godap.DisconnectRequest {
	return &godap.DisconnectRequest{
		Request: godap.Request{
			ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
			Command:         "disconnect",
		},
	}
}

func customRequest(seq int, command string) *godap.Request {
	return &godap.Request{
		ProtocolMessage: godap.ProtocolMessage{Seq: seq, Type: "request"},
		Command:         command,
	}
}
