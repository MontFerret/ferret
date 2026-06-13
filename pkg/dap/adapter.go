package dap

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	godap "github.com/google/go-dap"

	ferret "github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

// Adapter serves one DAP session over standard DAP framing.
type Adapter struct {
	state             adapterState
	inputCloser       io.Closer
	out               io.Writer
	log               io.Writer
	closeErr          error
	engine            *ferret.Engine
	session           *ferret.DebugSession
	plan              *ferret.Plan
	reader            *bufio.Reader
	pathFormat        string
	programPath       string
	cwd               string
	handlers          sync.WaitGroup
	seq               int
	nextFrameID       int
	nextVariableID    int
	shutdownOnce      sync.Once
	resourceCloseOnce sync.Once
	mu                sync.Mutex
	outMu             sync.Mutex
	launchHandled     bool
	terminated        bool
	resourcesClosed   bool
	adapterClosed     bool
	linesStartAt1     bool
	terminating       bool
	stopOnEntry       bool
	trace             bool
	running           bool
	started           bool
	launched          bool
	columnsStartAt1   bool
	initialized       bool
}

// New creates a DAP adapter over the provided transport.
func New(config Config) (*Adapter, error) {
	if config.In == nil {
		return nil, wrapStateError("input is required")
	}
	if config.Out == nil {
		return nil, wrapStateError("output is required")
	}

	adapter := &Adapter{
		reader:          newReader(config.In),
		out:             config.Out,
		log:             config.Log,
		trace:           config.Trace,
		pathFormat:      "path",
		linesStartAt1:   true,
		columnsStartAt1: true,
		nextFrameID:     1,
		nextVariableID:  1,
		state: adapterState{
			BreakpointsBySource: make(map[string]map[breakpointKey]ferret.DebugBreakpoint),
			FrameIndexes:        make(map[int]int),
			Handles:             make(map[int]variableHandle),
		},
	}

	if config.Log == nil {
		adapter.log = io.Discard
	}
	if closer, ok := config.In.(io.Closer); ok {
		adapter.inputCloser = closer
	}

	return adapter, nil
}

// Serve processes DAP requests until EOF, disconnect, or context cancellation.
func (a *Adapter) Serve(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	stopCancellation := context.AfterFunc(ctx, func() {
		if a.inputCloser != nil {
			_ = a.inputCloser.Close()
		}
	})
	defer stopCancellation()

	for {
		if ctx.Err() != nil {
			_ = a.shutdown()
			a.handlers.Wait()
			return a.shutdown()
		}

		message, metadata, err := readRequest(a.reader)
		if err != nil {
			if ctx.Err() != nil {
				_ = a.shutdown()
				a.handlers.Wait()
				return a.shutdown()
			}
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || a.isAdapterClosed() {
				_ = a.shutdown()
				a.handlers.Wait()
				return a.shutdown()
			}

			if a.hasSession() {
				_ = a.sendTerminated()
			}
			shutdownErr := a.shutdown()
			a.handlers.Wait()
			return errors.Join(err, shutdownErr)
		}

		a.traceMessage("recv", message)
		if a.isAdapterClosed() {
			a.handlers.Wait()
			return a.shutdown()
		}

		a.dispatch(ctx, message, metadata)
		if a.isAdapterClosed() {
			a.handlers.Wait()
			return a.shutdown()
		}
	}
}

func (a *Adapter) dispatch(ctx context.Context, message godap.Message, metadata requestMetadata) {
	switch request := message.(type) {
	case *godap.InitializeRequest:
		a.handleInitialize(request, metadata.Initialize)
	case *godap.LaunchRequest:
		a.handleLaunch(ctx, request)
	case *godap.SetBreakpointsRequest:
		a.handleSetBreakpoints(request, metadata.BreakpointColumns)
	case *godap.ConfigurationDoneRequest:
		a.handleConfigurationDone(ctx, request)
	case *godap.ThreadsRequest:
		a.handleThreads(request)
	case *godap.StackTraceRequest:
		a.handleStackTrace(request)
	case *godap.ScopesRequest:
		a.handleScopes(request)
	case *godap.VariablesRequest:
		a.handleVariables(request)
	case *godap.ContinueRequest:
		a.handleContinue(ctx, request)
	case *godap.NextRequest:
		a.handleNext(ctx, request)
	case *godap.StepInRequest:
		a.handleStepIn(ctx, request)
	case *godap.StepOutRequest:
		a.handleStepOut(ctx, request)
	case *godap.PauseRequest:
		a.handlePause(request)
	case *godap.EvaluateRequest:
		a.handleEvaluate(ctx, request)
	case *godap.DisconnectRequest:
		a.handleDisconnect(request)
	case *godap.SetExceptionBreakpointsRequest:
		a.handleSetExceptionBreakpoints(request)
	case *godap.Request:
		a.sendErrorResponse(request, fmt.Sprintf("request %q is not supported", request.Command))
	default:
		a.logf("ignoring unexpected DAP message type %T", message)
	}
}

func (a *Adapter) handleInitialize(request *godap.InitializeRequest, metadata initializeRequestMetadata) {
	a.mu.Lock()
	if a.initialized {
		a.mu.Unlock()
		a.sendErrorResponse(request, "initialize request was already handled")
		return
	}

	a.initialized = true
	a.linesStartAt1 = true
	if metadata.LinesStartAt1 != nil {
		a.linesStartAt1 = *metadata.LinesStartAt1
	}
	a.columnsStartAt1 = true
	if metadata.ColumnsStartAt1 != nil {
		a.columnsStartAt1 = *metadata.ColumnsStartAt1
	}
	if request.Arguments.PathFormat != "" {
		a.pathFormat = request.Arguments.PathFormat
	}
	a.mu.Unlock()

	response := &godap.InitializeResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
		Body: godap.Capabilities{
			SupportsConfigurationDoneRequest: true,
			SupportsEvaluateForHovers:        true,
		},
	}

	_ = a.send(response)
}

func (a *Adapter) handleLaunch(ctx context.Context, request *godap.LaunchRequest) {
	a.mu.Lock()
	if !a.initialized {
		a.mu.Unlock()
		a.sendErrorResponse(request, "initialize must be completed before launch")
		return
	}
	if a.launchHandled {
		a.mu.Unlock()
		a.sendErrorResponse(request, "launch has already been handled")
		return
	}
	a.mu.Unlock()

	var args launchArguments
	if err := json.Unmarshal(request.Arguments, &args); err != nil {
		a.sendErrorResponse(request, fmt.Sprintf("invalid launch arguments: %v", err))
		return
	}
	if strings.TrimSpace(args.Program) == "" {
		a.sendErrorResponse(request, "launch program is required")
		return
	}

	processCwd, err := os.Getwd()
	if err != nil {
		a.sendErrorResponse(request, fmt.Sprintf("resolve working directory: %v", err))
		return
	}

	resolvedCwd := processCwd
	if strings.TrimSpace(args.Cwd) != "" {
		resolvedCwd, err = resolveLaunchPath(strings.TrimSpace(args.Cwd), processCwd)
		if err != nil {
			a.sendErrorResponse(request, fmt.Sprintf("invalid launch cwd: %v", err))
			return
		}
	}

	programPath, err := resolveLaunchPath(strings.TrimSpace(args.Program), resolvedCwd)
	if err != nil {
		a.sendErrorResponse(request, fmt.Sprintf("invalid launch program: %v", err))
		return
	}

	src, err := source.Read(programPath)
	if err != nil {
		a.sendErrorResponse(request, fmt.Sprintf("read program: %v", err))
		return
	}

	engine, err := ferret.New(
		ferret.WithFSRoot(resolvedCwd),
		ferret.WithFSReadOnly(),
	)
	if err != nil {
		a.sendErrorResponse(request, formatError(err))
		return
	}

	plan, err := engine.CompileDebug(ctx, src)
	if err != nil {
		_ = engine.Close()
		a.sendErrorResponse(request, formatError(err))
		return
	}

	session, err := plan.NewDebugSession(ctx)
	if err != nil {
		_ = plan.Close()
		_ = engine.Close()
		a.sendErrorResponse(request, formatError(err))
		return
	}

	a.mu.Lock()
	a.cwd = resolvedCwd
	a.programPath = programPath
	a.stopOnEntry = args.StopOnEntry
	a.engine = engine
	a.plan = plan
	a.session = session
	a.launchHandled = true
	a.launched = true
	a.resourcesClosed = false
	a.terminating = false
	a.terminated = false
	a.state.BreakpointsBySource = make(map[string]map[breakpointKey]ferret.DebugBreakpoint)
	a.clearPausedStateLocked()
	a.mu.Unlock()

	response := &godap.LaunchResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
	}
	if err := a.send(response); err != nil {
		return
	}

	_ = a.send(&godap.InitializedEvent{
		Event: godap.Event{
			ProtocolMessage: godap.ProtocolMessage{Type: "event"},
			Event:           "initialized",
		},
	})
}

func (a *Adapter) handleSetBreakpoints(request *godap.SetBreakpointsRequest, breakpointColumns []bool) {
	session, launched, running, programPath, cwd := a.snapshotBreakpointsState()
	if !launched {
		a.sendErrorResponse(request, "launch must be completed before setBreakpoints")
		return
	}
	if running {
		a.sendErrorResponse(request, "breakpoints cannot be changed while execution is running")
		return
	}

	requested := request.Arguments.Breakpoints
	if len(requested) == 0 && len(request.Arguments.Lines) > 0 {
		requested = make([]godap.SourceBreakpoint, 0, len(request.Arguments.Lines))
		for _, line := range request.Arguments.Lines {
			requested = append(requested, godap.SourceBreakpoint{Line: line})
		}
	}

	sourcePath := request.Arguments.Source.Path
	if strings.TrimSpace(sourcePath) == "" {
		response := &godap.SetBreakpointsResponse{
			Response: godap.Response{
				ProtocolMessage: godap.ProtocolMessage{Type: "response"},
				RequestSeq:      request.Seq,
				Success:         true,
				Command:         request.Command,
			},
			Body: godap.SetBreakpointsResponseBody{
				Breakpoints: a.invalidBreakpoints(requested, "Invalid breakpoint source."),
			},
		}
		_ = a.send(response)
		return
	}

	normalized, matched, err := normalizeSourcePath(sourcePath, cwd, programPath)
	if err != nil || !matched {
		message := "Invalid breakpoint source."
		if err != nil {
			message = fmt.Sprintf("Invalid breakpoint source: %v", err)
		}
		response := &godap.SetBreakpointsResponse{
			Response: godap.Response{
				ProtocolMessage: godap.ProtocolMessage{Type: "response"},
				RequestSeq:      request.Seq,
				Success:         true,
				Command:         request.Command,
			},
			Body: godap.SetBreakpointsResponseBody{
				Breakpoints: a.invalidBreakpoints(requested, message),
			},
		}
		_ = a.send(response)
		return
	}

	existing := a.breakpointsForSource(normalized)
	desired := make(map[breakpointKey]ferret.DebugBreakpoint)
	responseBreakpoints := make([]godap.Breakpoint, 0, len(requested))

	for index, breakpoint := range requested {
		line := a.clientLineToFerret(breakpoint.Line)
		columnPresent := index < len(breakpointColumns) && breakpointColumns[index]
		column := a.clientColumnToFerret(breakpoint.Column, columnPresent)

		if message := unsupportedBreakpointMessage(breakpoint); message != "" {
			responseBreakpoints = append(responseBreakpoints, godap.Breakpoint{
				Verified: false,
				Message:  message,
				Line:     breakpoint.Line,
				Column:   breakpoint.Column,
			})
			continue
		}

		key := breakpointKey{Line: line, Column: column}
		if desiredBreakpoint, ok := desired[key]; ok {
			responseBreakpoints = append(responseBreakpoints, a.adaptBreakpoint(desiredBreakpoint))
			continue
		}
		if existingBreakpoint, ok := existing[key]; ok {
			desired[key] = existingBreakpoint
			delete(existing, key)
			responseBreakpoints = append(responseBreakpoints, a.adaptBreakpoint(existingBreakpoint))
			continue
		}

		created, err := session.SetBreakpointAt(
			ferret.DebugSourceLocation{File: normalized, Line: line, Column: column},
			ferret.DebugBreakpointOptions{BindingMode: ferret.DebugBreakpointBindNextExecutableInFunction},
		)
		if err != nil {
			responseBreakpoints = append(responseBreakpoints, godap.Breakpoint{
				Verified: false,
				Message:  formatError(err),
				Line:     breakpoint.Line,
				Column:   breakpoint.Column,
			})
			continue
		}

		desired[key] = created
		responseBreakpoints = append(responseBreakpoints, a.adaptBreakpoint(created))
	}

	for _, breakpoint := range existing {
		if err := session.DeleteBreakpoint(breakpoint.ID); err != nil {
			a.sendErrorResponse(request, formatError(err))
			return
		}
	}

	a.setBreakpointsForSource(normalized, desired)

	response := &godap.SetBreakpointsResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
		Body: godap.SetBreakpointsResponseBody{
			Breakpoints: responseBreakpoints,
		},
	}

	_ = a.send(response)
}

func (a *Adapter) handleConfigurationDone(ctx context.Context, request *godap.ConfigurationDoneRequest) {
	session, started, launched := a.snapshotLaunchState()
	if !launched {
		a.sendErrorResponse(request, "launch must be completed before configurationDone")
		return
	}
	if started {
		a.sendErrorResponse(request, "configurationDone has already been handled")
		return
	}

	event, err := session.Start(ctx)
	if err != nil {
		a.sendErrorResponse(request, formatError(err))
		return
	}

	a.mu.Lock()
	a.started = true
	a.mu.Unlock()

	response := &godap.ConfigurationDoneResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
	}
	if err := a.send(response); err != nil {
		return
	}

	if event.Reason == ferret.DebugReasonEntry && !a.stopOnEntry {
		if hitBreakpointIDs := a.entryBreakpointIDs(event.Location); len(hitBreakpointIDs) > 0 {
			a.publishStop(session, event, "breakpoint", hitBreakpointIDs)
			return
		}

		a.mu.Lock()
		a.running = true
		a.clearPausedStateLocked()
		a.mu.Unlock()

		a.startExecution(func() (*ferret.DebugEvent, error) {
			return session.Continue(ctx)
		}, session)
		return
	}

	a.handleExecutionEvent(session, event)
}

func (a *Adapter) handleThreads(request *godap.ThreadsRequest) {
	if !a.hasSession() {
		a.sendErrorResponse(request, "launch must be completed before threads")
		return
	}

	response := &godap.ThreadsResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
		Body: godap.ThreadsResponseBody{
			Threads: []godap.Thread{{Id: 1, Name: "main"}},
		},
	}

	_ = a.send(response)
}

func (a *Adapter) handleStackTrace(request *godap.StackTraceRequest) {
	frames, ok := a.framesForStoppedState(request.Arguments.ThreadId)
	if !ok {
		a.sendErrorResponse(request, "stackTrace requires a stopped thread")
		return
	}

	start := request.Arguments.StartFrame
	levels := request.Arguments.Levels
	sliced := sliceVariables(frames, start, levels)
	stackFrames := make([]godap.StackFrame, 0, len(sliced))
	for _, frame := range sliced {
		stackFrames = append(stackFrames, godap.StackFrame{
			Id:     frame.ID,
			Name:   frame.Frame.Name,
			Source: a.sourceDescriptor(frame.Frame.Location.File),
			Line:   a.ferretLineToClient(frame.Frame.Location.Line),
			Column: a.ferretColumnToClient(frame.Frame.Location.Column),
		})
	}

	response := &godap.StackTraceResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
		Body: godap.StackTraceResponseBody{
			StackFrames: stackFrames,
			TotalFrames: len(frames),
		},
	}

	_ = a.send(response)
}

func (a *Adapter) handleScopes(request *godap.ScopesRequest) {
	session, frame, ok := a.frameForRequest(request.Arguments.FrameId)
	if !ok {
		a.sendErrorResponse(request, "invalid frame ID")
		return
	}

	scopes := make([]godap.Scope, 0, 2)
	if frameIndex := a.frameIndex(request.Arguments.FrameId); frameIndex > 0 {
		handleID := a.storeScopeHandle(nil)
		scopes = append(scopes, godap.Scope{
			Name:               "Locals",
			PresentationHint:   "locals",
			VariablesReference: handleID,
			Expensive:          false,
			Source:             a.sourceDescriptor(frame.Location.File),
			Line:               a.ferretLineToClient(frame.Location.Line),
			Column:             a.ferretColumnToClient(frame.Location.Column),
		})
	} else {
		variables, err := session.Locals()
		if err != nil {
			a.sendErrorResponse(request, formatError(err))
			return
		}

		locals, params := splitVariables(variables)
		localsHandle := a.storeScopeHandle(locals)
		scopes = append(scopes, godap.Scope{
			Name:               "Locals",
			PresentationHint:   "locals",
			VariablesReference: localsHandle,
			NamedVariables:     len(locals),
			Expensive:          false,
			Source:             a.sourceDescriptor(frame.Location.File),
			Line:               a.ferretLineToClient(frame.Location.Line),
			Column:             a.ferretColumnToClient(frame.Location.Column),
		})

		if len(params) > 0 {
			paramsHandle := a.storeScopeHandle(params)
			scopes = append(scopes, godap.Scope{
				Name:               "Parameters",
				PresentationHint:   "arguments",
				VariablesReference: paramsHandle,
				NamedVariables:     len(params),
				Expensive:          false,
				Source:             a.sourceDescriptor(frame.Location.File),
				Line:               a.ferretLineToClient(frame.Location.Line),
				Column:             a.ferretColumnToClient(frame.Location.Column),
			})
		}
	}

	response := &godap.ScopesResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
		Body: godap.ScopesResponseBody{Scopes: scopes},
	}

	_ = a.send(response)
}

func (a *Adapter) handleVariables(request *godap.VariablesRequest) {
	session, handle, ok := a.variableHandleForRequest(request.Arguments.VariablesReference)
	if !ok {
		a.sendErrorResponse(request, "invalid variables reference")
		return
	}

	var variables []ferret.DebugVariable
	switch handle.Kind {
	case variableHandleScope:
		if request.Arguments.Filter == "indexed" {
			variables = nil
		} else {
			variables = append([]ferret.DebugVariable(nil), handle.ScopeVariables...)
		}
	case variableHandleValue:
		var err error
		variables, err = session.Variables(handle.ValueReference)
		if err != nil {
			a.sendErrorResponse(request, formatError(err))
			return
		}
	default:
		a.sendErrorResponse(request, "invalid variables reference")
		return
	}

	variables = sliceVariables(variables, request.Arguments.Start, request.Arguments.Count)
	adapted := make([]godap.Variable, 0, len(variables))
	for _, variable := range variables {
		adapted = append(adapted, a.adaptVariable(variable))
	}

	response := &godap.VariablesResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
		Body: godap.VariablesResponseBody{Variables: adapted},
	}

	_ = a.send(response)
}

func (a *Adapter) handleContinue(ctx context.Context, request *godap.ContinueRequest) {
	a.handleResumeRequest(
		ctx,
		request,
		request.Arguments.ThreadId,
		&godap.ContinueResponse{
			Response: godap.Response{
				ProtocolMessage: godap.ProtocolMessage{Type: "response"},
				RequestSeq:      request.Seq,
				Success:         true,
				Command:         request.Command,
			},
			Body: godap.ContinueResponseBody{AllThreadsContinued: true},
		},
		func(session *ferret.DebugSession) (*ferret.DebugEvent, error) {
			return session.Continue(ctx)
		},
	)
}

func (a *Adapter) handleNext(ctx context.Context, request *godap.NextRequest) {
	a.handleResumeRequest(
		ctx,
		request,
		request.Arguments.ThreadId,
		&godap.NextResponse{
			Response: godap.Response{
				ProtocolMessage: godap.ProtocolMessage{Type: "response"},
				RequestSeq:      request.Seq,
				Success:         true,
				Command:         request.Command,
			},
		},
		func(session *ferret.DebugSession) (*ferret.DebugEvent, error) {
			return session.Next(ctx)
		},
	)
}

func (a *Adapter) handleStepIn(ctx context.Context, request *godap.StepInRequest) {
	a.handleResumeRequest(
		ctx,
		request,
		request.Arguments.ThreadId,
		&godap.StepInResponse{
			Response: godap.Response{
				ProtocolMessage: godap.ProtocolMessage{Type: "response"},
				RequestSeq:      request.Seq,
				Success:         true,
				Command:         request.Command,
			},
		},
		func(session *ferret.DebugSession) (*ferret.DebugEvent, error) {
			return session.Step(ctx)
		},
	)
}

func (a *Adapter) handleStepOut(ctx context.Context, request *godap.StepOutRequest) {
	a.handleResumeRequest(
		ctx,
		request,
		request.Arguments.ThreadId,
		&godap.StepOutResponse{
			Response: godap.Response{
				ProtocolMessage: godap.ProtocolMessage{Type: "response"},
				RequestSeq:      request.Seq,
				Success:         true,
				Command:         request.Command,
			},
		},
		func(session *ferret.DebugSession) (*ferret.DebugEvent, error) {
			return session.Out(ctx)
		},
	)
}

func (a *Adapter) handlePause(request *godap.PauseRequest) {
	session, ok := a.sessionForPause(request.Arguments.ThreadId)
	if !ok {
		a.sendErrorResponse(request, "pause requires a running thread")
		return
	}

	response := &godap.PauseResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
	}

	if err := a.send(response); err != nil {
		return
	}
	if err := session.Pause(); err != nil {
		a.handleAsyncFailure(err)
	}
}

func (a *Adapter) handleEvaluate(ctx context.Context, request *godap.EvaluateRequest) {
	switch request.Arguments.Context {
	case "repl", "watch", "hover":
	default:
		a.sendErrorResponse(request, fmt.Sprintf("evaluate context %q is not supported", request.Arguments.Context))
		return
	}

	session, topFrameID, ok := a.sessionForEvaluation(request.Arguments.FrameId)
	if !ok {
		a.sendErrorResponse(request, "evaluate requires the current top frame")
		return
	}
	if request.Arguments.FrameId != 0 && request.Arguments.FrameId != topFrameID {
		a.sendErrorResponse(request, "evaluate is only supported for the current top frame")
		return
	}

	value, err := session.Evaluate(ctx, request.Arguments.Expression)
	if err != nil {
		a.sendErrorResponse(request, formatError(err))
		return
	}

	response := &godap.EvaluateResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
		Body: godap.EvaluateResponseBody{
			Result:             value.Display,
			Type:               value.Type,
			VariablesReference: a.storeValueHandle(value.Reference),
		},
	}

	_ = a.send(response)
}

func (a *Adapter) handleDisconnect(request *godap.DisconnectRequest) {
	a.beginTermination()

	response := &godap.DisconnectResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
	}
	if err := a.send(response); err != nil {
		return
	}

	_ = a.sendTerminated()
	_ = a.shutdown()
}

func (a *Adapter) handleSetExceptionBreakpoints(request *godap.SetExceptionBreakpointsRequest) {
	response := &godap.SetExceptionBreakpointsResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      request.Seq,
			Success:         true,
			Command:         request.Command,
		},
	}

	_ = a.send(response)
}

func (a *Adapter) handleResumeRequest(
	ctx context.Context,
	request godap.RequestMessage,
	threadID int,
	response godap.Message,
	resume func(*ferret.DebugSession) (*ferret.DebugEvent, error),
) {
	session, ok := a.sessionForResume(threadID)
	if !ok {
		a.sendErrorResponse(request, "execution is not stopped")
		return
	}

	if err := a.send(response); err != nil {
		return
	}

	a.startExecution(func() (*ferret.DebugEvent, error) {
		return resume(session)
	}, session)
}

func (a *Adapter) startExecution(
	resume func() (*ferret.DebugEvent, error),
	session *ferret.DebugSession,
) {
	a.handlers.Add(1)
	go func() {
		defer a.handlers.Done()

		event, err := resume()
		if err != nil {
			a.handleAsyncFailure(err)
			return
		}

		a.handleExecutionEvent(session, event)
	}()
}

func (a *Adapter) handleExecutionEvent(session *ferret.DebugSession, event *ferret.DebugEvent) {
	if event == nil {
		a.handleAsyncFailure(wrapStateError("execution returned no event"))
		return
	}

	switch event.Reason {
	case ferret.DebugReasonCompleted, ferret.DebugReasonTerminated:
		a.mu.Lock()
		a.running = false
		a.clearPausedStateLocked()
		a.mu.Unlock()
		_ = a.sendTerminated()
		_ = a.closeResources()
	case ferret.DebugReasonEntry:
		a.publishStop(session, event, "entry", nil)
	case ferret.DebugReasonBreakpoint:
		a.publishStop(session, event, "breakpoint", breakpointIDs(event.HitBreakpointIDs))
	case ferret.DebugReasonStep:
		a.publishStop(session, event, "step", nil)
	case ferret.DebugReasonPause:
		a.publishStop(session, event, "pause", nil)
	case ferret.DebugReasonRuntimeError:
		a.publishStop(session, event, "exception", nil)
	default:
		a.handleAsyncFailure(fmt.Errorf("unsupported debug reason %q", event.Reason))
	}
}

func (a *Adapter) publishStop(session *ferret.DebugSession, event *ferret.DebugEvent, reason string, hitBreakpointIDs []int) {
	frames, err := session.Frames()
	if err != nil {
		a.handleAsyncFailure(err)
		return
	}

	a.mu.Lock()
	if a.terminating || a.terminated || a.adapterClosed {
		a.mu.Unlock()
		return
	}
	a.running = false
	a.storeFramesLocked(frames)
	a.mu.Unlock()

	body := godap.StoppedEventBody{
		Reason:            reason,
		ThreadId:          1,
		AllThreadsStopped: true,
		HitBreakpointIds:  hitBreakpointIDs,
	}
	if event.Error != nil {
		description := formatError(event.Error)
		body.Description = description
		body.Text = description
	}

	_ = a.sendWhileActive(&godap.StoppedEvent{
		Event: godap.Event{
			ProtocolMessage: godap.ProtocolMessage{Type: "event"},
			Event:           "stopped",
		},
		Body: body,
	})
}

func (a *Adapter) handleAsyncFailure(err error) {
	if err == nil || a.isTerminating() {
		return
	}

	a.logf("async execution failure: %s", formatError(err))
	a.mu.Lock()
	a.closeErr = errors.Join(a.closeErr, err)
	a.mu.Unlock()
	_ = a.sendTerminated()
	_ = a.shutdown()
}

func (a *Adapter) send(message godap.Message) error {
	a.outMu.Lock()
	err := a.writeMessageLocked(message)
	a.outMu.Unlock()
	a.handleWriteError(err)

	return err
}

func (a *Adapter) sendWhileActive(message godap.Message) error {
	a.outMu.Lock()
	a.mu.Lock()
	active := !a.terminating && !a.terminated && !a.adapterClosed
	a.mu.Unlock()
	if !active {
		a.outMu.Unlock()
		return nil
	}

	err := a.writeMessageLocked(message)
	a.outMu.Unlock()
	a.handleWriteError(err)

	return err
}

func (a *Adapter) writeMessageLocked(message godap.Message) error {
	a.seq++
	switch message := message.(type) {
	case godap.ResponseMessage:
		message.GetResponse().Seq = a.seq
		message.GetResponse().Type = "response"
	case godap.EventMessage:
		message.GetEvent().Seq = a.seq
		message.GetEvent().Type = "event"
	case godap.RequestMessage:
		message.GetRequest().Seq = a.seq
		message.GetRequest().Type = "request"
	}

	a.traceMessage("send", message)
	return godap.WriteProtocolMessage(a.out, message)
}

func (a *Adapter) handleWriteError(err error) {
	if err != nil {
		a.mu.Lock()
		a.closeErr = errors.Join(a.closeErr, err)
		a.mu.Unlock()
		_ = a.shutdown()
	}
}

func (a *Adapter) sendErrorResponse(request godap.RequestMessage, message string) {
	req := request.GetRequest()
	_ = a.send(&godap.ErrorResponse{
		Response: godap.Response{
			ProtocolMessage: godap.ProtocolMessage{Type: "response"},
			RequestSeq:      req.Seq,
			Success:         false,
			Command:         req.Command,
			Message:         message,
		},
		Body: godap.ErrorResponseBody{
			Error: &godap.ErrorMessage{
				Id:       1,
				Format:   message,
				ShowUser: true,
			},
		},
	})
}

func (a *Adapter) sendTerminated() error {
	a.mu.Lock()
	if a.terminated {
		a.mu.Unlock()
		return nil
	}
	a.terminating = true
	a.terminated = true
	a.running = false
	a.clearPausedStateLocked()
	a.mu.Unlock()

	return a.send(&godap.TerminatedEvent{
		Event: godap.Event{
			ProtocolMessage: godap.ProtocolMessage{Type: "event"},
			Event:           "terminated",
		},
	})
}

func (a *Adapter) beginTermination() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.terminating = true
	a.running = false
	a.clearPausedStateLocked()
}

func (a *Adapter) closeResources() error {
	a.resourceCloseOnce.Do(func() {
		a.mu.Lock()
		session := a.session
		plan := a.plan
		engine := a.engine
		a.session = nil
		a.plan = nil
		a.engine = nil
		a.launched = false
		a.running = false
		a.resourcesClosed = true
		a.clearPausedStateLocked()
		a.state.BreakpointsBySource = make(map[string]map[breakpointKey]ferret.DebugBreakpoint)
		a.mu.Unlock()

		var closeErr error
		if session != nil {
			closeErr = errors.Join(closeErr, session.Close())
		}
		if plan != nil {
			closeErr = errors.Join(closeErr, plan.Close())
		}
		if engine != nil {
			closeErr = errors.Join(closeErr, engine.Close())
		}

		a.mu.Lock()
		a.closeErr = errors.Join(a.closeErr, closeErr)
		a.mu.Unlock()
	})

	a.mu.Lock()
	defer a.mu.Unlock()

	return a.closeErr
}

func (a *Adapter) shutdown() error {
	a.shutdownOnce.Do(func() {
		a.mu.Lock()
		a.adapterClosed = true
		a.mu.Unlock()

		_ = a.closeResources()
		var inputErr error
		if a.inputCloser != nil {
			inputErr = a.inputCloser.Close()
		}

		a.mu.Lock()
		a.closeErr = errors.Join(a.closeErr, inputErr)
		a.mu.Unlock()
	})

	a.mu.Lock()
	defer a.mu.Unlock()

	return a.closeErr
}

func (a *Adapter) snapshotBreakpointsState() (*ferret.DebugSession, bool, bool, string, string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.session, a.launched, a.running, a.programPath, a.cwd
}

func (a *Adapter) snapshotLaunchState() (*ferret.DebugSession, bool, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.session, a.started, a.launched
}

func (a *Adapter) hasSession() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.launchHandled
}

func (a *Adapter) framesForStoppedState(threadID int) ([]frameState, bool) {
	if threadID != 1 {
		return nil, false
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running || len(a.state.Frames) == 0 {
		return nil, false
	}

	frames := append([]frameState(nil), a.state.Frames...)
	return frames, true
}

func (a *Adapter) frameForRequest(frameID int) (*ferret.DebugSession, ferret.DebugFrame, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	index, ok := a.state.FrameIndexes[frameID]
	if !ok || a.session == nil || a.running || index >= len(a.state.Frames) {
		return nil, ferret.DebugFrame{}, false
	}

	return a.session, a.state.Frames[index].Frame, true
}

func (a *Adapter) frameIndex(frameID int) int {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.state.FrameIndexes[frameID]
}

func (a *Adapter) variableHandleForRequest(reference int) (*ferret.DebugSession, variableHandle, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	handle, ok := a.state.Handles[reference]
	if !ok || a.session == nil || a.running {
		return nil, variableHandle{}, false
	}

	return a.session, handle, true
}

func (a *Adapter) sessionForResume(threadID int) (*ferret.DebugSession, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if threadID != 1 || a.session == nil || !a.started || a.running || len(a.state.Frames) == 0 {
		return nil, false
	}

	session := a.session
	a.running = true
	a.clearPausedStateLocked()

	return session, true
}

func (a *Adapter) sessionForPause(threadID int) (*ferret.DebugSession, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if threadID != 1 || a.session == nil || !a.running {
		return nil, false
	}

	return a.session, true
}

func (a *Adapter) sessionForEvaluation(frameID int) (*ferret.DebugSession, int, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.session == nil || a.running || len(a.state.Frames) == 0 {
		return nil, 0, false
	}

	topFrameID := a.state.Frames[0].ID
	return a.session, topFrameID, true
}

func (a *Adapter) breakpointsForSource(path string) map[breakpointKey]ferret.DebugBreakpoint {
	a.mu.Lock()
	defer a.mu.Unlock()

	current := a.state.BreakpointsBySource[path]
	out := make(map[breakpointKey]ferret.DebugBreakpoint, len(current))
	for key, breakpoint := range current {
		out[key] = breakpoint
	}

	return out
}

func (a *Adapter) setBreakpointsForSource(path string, breakpoints map[breakpointKey]ferret.DebugBreakpoint) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(breakpoints) == 0 {
		delete(a.state.BreakpointsBySource, path)
		return
	}

	a.state.BreakpointsBySource[path] = breakpoints
}

func (a *Adapter) entryBreakpointIDs(location ferret.DebugLocation) []int {
	a.mu.Lock()
	defer a.mu.Unlock()

	current := a.state.BreakpointsBySource[a.programPath]
	if len(current) == 0 {
		return nil
	}

	ids := make([]int, 0, len(current))
	for _, breakpoint := range current {
		if breakpoint.Bound && breakpoint.Line == location.Line && breakpoint.Column == location.Column {
			ids = append(ids, int(breakpoint.ID))
		}
	}

	sort.Ints(ids)
	return ids
}

func (a *Adapter) adaptBreakpoint(breakpoint ferret.DebugBreakpoint) godap.Breakpoint {
	out := godap.Breakpoint{
		Id:       int(breakpoint.ID),
		Verified: breakpoint.Bound,
	}

	if breakpoint.Bound {
		out.Line = a.ferretLineToClient(breakpoint.Line)
		out.Column = a.ferretColumnToClient(breakpoint.Column)
		return out
	}

	out.Line = a.ferretLineToClient(breakpoint.RequestedLine)
	out.Column = a.ferretColumnToClient(breakpoint.RequestedColumn)
	out.Message = "Breakpoint could not be bound to executable code."
	return out
}

func (a *Adapter) invalidBreakpoints(requested []godap.SourceBreakpoint, message string) []godap.Breakpoint {
	out := make([]godap.Breakpoint, 0, len(requested))
	for _, breakpoint := range requested {
		out = append(out, godap.Breakpoint{
			Verified: false,
			Message:  message,
			Line:     breakpoint.Line,
			Column:   breakpoint.Column,
		})
	}

	return out
}

func (a *Adapter) adaptVariable(variable ferret.DebugVariable) godap.Variable {
	return godap.Variable{
		Name:               variable.Name,
		Value:              variable.Value.Display,
		Type:               variable.Value.Type,
		EvaluateName:       variable.Name,
		VariablesReference: a.storeValueHandle(variable.Value.Reference),
	}
}

func (a *Adapter) sourceDescriptor(path string) *godap.Source {
	if strings.TrimSpace(path) == "" {
		return nil
	}

	return &godap.Source{
		Name:            filepath.Base(path),
		Path:            encodeClientPath(path, a.pathFormat),
		SourceReference: 0,
	}
}

func (a *Adapter) storeFramesLocked(frames []ferret.DebugFrame) {
	a.state.Frames = make([]frameState, 0, len(frames))
	a.state.FrameIndexes = make(map[int]int, len(frames))
	a.state.Handles = make(map[int]variableHandle)

	for index, frame := range frames {
		id := a.nextFrameID
		a.nextFrameID++
		a.state.FrameIndexes[id] = index
		a.state.Frames = append(a.state.Frames, frameState{ID: id, Frame: frame})
	}
}

func (a *Adapter) clearPausedStateLocked() {
	a.state.Frames = nil
	a.state.FrameIndexes = make(map[int]int)
	a.state.Handles = make(map[int]variableHandle)
}

func (a *Adapter) storeScopeHandle(variables []ferret.DebugVariable) int {
	a.mu.Lock()
	defer a.mu.Unlock()

	id := a.nextVariableID
	a.nextVariableID++
	a.state.Handles[id] = variableHandle{
		ScopeVariables: append([]ferret.DebugVariable(nil), variables...),
		Kind:           variableHandleScope,
	}

	return id
}

func (a *Adapter) storeValueHandle(reference ferret.DebugValueReference) int {
	if !reference.Valid() {
		return 0
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	id := a.nextVariableID
	a.nextVariableID++
	a.state.Handles[id] = variableHandle{
		ValueReference: reference,
		Kind:           variableHandleValue,
	}

	return id
}

func (a *Adapter) ferretLineToClient(line int) int {
	if a.linesStartAt1 {
		return line
	}

	if line <= 0 {
		return 0
	}

	return line - 1
}

func (a *Adapter) ferretColumnToClient(column int) int {
	if column <= 0 {
		return 0
	}
	if a.columnsStartAt1 {
		return column
	}

	return column - 1
}

func (a *Adapter) clientLineToFerret(line int) int {
	if a.linesStartAt1 {
		return line
	}

	return line + 1
}

func (a *Adapter) clientColumnToFerret(column int, present bool) int {
	if !present {
		return 0
	}
	if a.columnsStartAt1 {
		return column
	}

	return column + 1
}

func (a *Adapter) traceMessage(direction string, message godap.Message) {
	if !a.trace {
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		a.logf("%s <marshal error: %v>", direction, err)
		return
	}

	a.logf("%s %s", direction, data)
}

func (a *Adapter) logf(format string, args ...any) {
	if a.log == nil {
		return
	}

	_, _ = fmt.Fprintf(a.log, format+"\n", args...)
}

func (a *Adapter) isAdapterClosed() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.adapterClosed
}

func (a *Adapter) isTerminating() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.terminating || a.terminated || a.adapterClosed
}
