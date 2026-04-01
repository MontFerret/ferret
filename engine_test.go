package ferret

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	formatjson "github.com/MontFerret/ferret/v2/pkg/bytecode/format/json"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/goccy/go-json"
)

const (
	coverageValidQuery   = "RETURN 1"
	coverageInvalidQuery = "LET"
)

type coverageModule struct {
	registerFn func(boot Bootstrap) error
}

func (m coverageModule) Name() string {
	return "coverage-module"
}

func (m coverageModule) Register(boot Bootstrap) error {
	if m.registerFn == nil {
		return nil
	}

	return m.registerFn(boot)
}

func coverageVarFn(context.Context, ...runtime.Value) (runtime.Value, error) {
	return runtime.None, nil
}

func mustNewEngine(t *testing.T, setters ...Option) *Engine {
	t.Helper()

	eng, err := New(setters...)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	return eng
}

func mustCompilePlan(t *testing.T, eng *Engine, query string) *Plan {
	t.Helper()

	plan, err := eng.Compile(context.Background(), source.NewAnonymous(query))
	if err != nil {
		t.Fatalf("failed to compile query %q: %v", query, err)
	}

	return plan
}

func mustNewSession(t *testing.T, plan *Plan, setters ...SessionOption) *Session {
	t.Helper()

	session, err := plan.NewSession(context.Background(), setters...)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	return session
}

func mustMarshalArtifact(t *testing.T, query string, opts ...artifact.Options) []byte {
	t.Helper()

	prog, err := compiler.New().Compile(source.NewAnonymous(query))
	if err != nil {
		t.Fatalf("failed to compile query %q: %v", query, err)
	}

	marshalOpts := artifact.Options{}
	if len(opts) > 0 {
		marshalOpts = opts[0]
	}

	data, err := artifact.Marshal(prog, marshalOpts)
	if err != nil {
		t.Fatalf("failed to marshal artifact: %v", err)
	}

	return data
}

func TestEngineNewReturnsOptionError(t *testing.T) {
	t.Parallel()

	eng, err := New(WithModules())
	if err != nil {
		t.Fatalf("expected New to succeed, got: %v", err)
	}

	if eng == nil {
		t.Fatal("expected engine to be non-nil on successful construction")
	}
}

func TestEngineNewReturnsModuleRegistrationError(t *testing.T) {
	t.Parallel()

	registerErr := errors.New("module register failed")
	closeCalled := false

	eng, err := New(WithModules(
		coverageModule{
			registerFn: func(boot Bootstrap) error {
				boot.Hooks().Engine().OnClose(func() error {
					closeCalled = true
					return nil
				})

				return nil
			},
		},
		coverageModule{
			registerFn: func(Bootstrap) error {
				return registerErr
			},
		},
	))
	if err == nil {
		t.Fatal("expected New to fail on module registration error")
	}

	if eng != nil {
		t.Fatal("expected engine to be nil on module registration error")
	}

	if !closeCalled {
		t.Fatal("expected close hooks to run on module registration error")
	}

	if !errors.Is(err, registerErr) {
		t.Fatalf("expected registration error to be preserved, got: %v", err)
	}
}

func TestEngineNewJoinsModuleRegistrationAndCloseHookErrors(t *testing.T) {
	t.Parallel()

	registerErr := errors.New("module register failed")
	closeErr := errors.New("close hook failed")

	_, err := New(WithModules(
		coverageModule{
			registerFn: func(boot Bootstrap) error {
				boot.Hooks().Engine().OnClose(func() error {
					return closeErr
				})

				return nil
			},
		},
		coverageModule{
			registerFn: func(Bootstrap) error {
				return registerErr
			},
		},
	))
	if err == nil {
		t.Fatal("expected New to fail when module registration and close hooks fail")
	}

	if !errors.Is(err, registerErr) {
		t.Fatalf("expected error to include module registration error, got: %v", err)
	}

	if !errors.Is(err, closeErr) {
		t.Fatalf("expected error to include close hook error, got: %v", err)
	}

	if !strings.Contains(err.Error(), "close hooks") {
		t.Fatalf("expected error to include close hooks label, got: %v", err)
	}
}

func TestEngineNewJoinsHostBuildAndCloseHookErrors(t *testing.T) {
	t.Parallel()

	closeErr := errors.New("close hook failed")

	_, err := New(WithModules(coverageModule{
		registerFn: func(boot Bootstrap) error {
			boot.Hooks().Engine().OnClose(func() error {
				return closeErr
			})

			boot.Host().Library().Function().A0().Add("COVER_DUP", func(context.Context) (runtime.Value, error) {
				return runtime.None, nil
			})
			boot.Host().Library().Function().A0().Add("COVER_DUP", func(context.Context) (runtime.Value, error) {
				return runtime.None, nil
			})

			return nil
		},
	}))
	if err == nil {
		t.Fatal("expected New to fail when host build and close hooks fail")
	}

	if !errors.Is(err, closeErr) {
		t.Fatalf("expected error to include close hook failure, got: %v", err)
	}

	if !strings.Contains(err.Error(), "close hooks") {
		t.Fatalf("expected error to include close hooks label, got: %v", err)
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected error to include host build failure details, got: %v", err)
	}
}

func TestEngineNewReturnsInitHookErrorAndRunsCleanup(t *testing.T) {
	t.Parallel()

	initErr := errors.New("init hook failed")
	closeCalled := false

	_, err := New(
		WithEngineInitHook(func() error {
			return initErr
		}),
		WithEngineCloseHook(func() error {
			closeCalled = true
			return nil
		}),
	)
	if err == nil {
		t.Fatal("expected New to fail when init hook fails")
	}

	if !closeCalled {
		t.Fatal("expected close hooks to run after init hook failure")
	}

	if !errors.Is(err, initErr) {
		t.Fatalf("expected error to include init hook failure, got: %v", err)
	}

	if !strings.Contains(err.Error(), "init hooks") {
		t.Fatalf("expected error to include init hooks label, got: %v", err)
	}
}

func TestEngineCompileReturnsBeforeHookError(t *testing.T) {
	t.Parallel()

	beforeErr := errors.New("before compile failed")
	afterCalled := false
	eng := mustNewEngine(t,
		WithBeforeCompileHook(func(context.Context) error {
			return beforeErr
		}),
		WithAfterCompileHook(func(context.Context, error) error {
			afterCalled = true
			return nil
		}),
	)

	plan, err := eng.Compile(context.Background(), source.NewAnonymous(coverageValidQuery))
	if err == nil {
		t.Fatal("expected compile to fail on before hook error")
	}

	if plan != nil {
		t.Fatal("expected plan to be nil when before hook fails")
	}

	if !errors.Is(err, beforeErr) {
		t.Fatalf("expected before hook error to be preserved, got: %v", err)
	}

	if !strings.Contains(err.Error(), "before compile hooks") {
		t.Fatalf("expected before compile hooks label, got: %v", err)
	}

	if afterCalled {
		t.Fatal("after compile hooks must not run when before compile hook fails")
	}
}

func TestEngineCompileReturnsCompilerErrorWhenAfterHooksSucceed(t *testing.T) {
	t.Parallel()

	afterCalled := false
	var compileErrSeen error

	eng := mustNewEngine(t, WithAfterCompileHook(func(_ context.Context, err error) error {
		afterCalled = true
		compileErrSeen = err
		return nil
	}))

	plan, err := eng.Compile(context.Background(), source.NewAnonymous(coverageInvalidQuery))
	if err == nil {
		t.Fatal("expected compile to fail for invalid query")
	}

	if plan != nil {
		t.Fatal("expected plan to be nil when compilation fails")
	}

	if !afterCalled {
		t.Fatal("expected after compile hooks to run even on compile failure")
	}

	if compileErrSeen == nil {
		t.Fatal("expected after compile hook to receive compile error")
	}

	if !errors.Is(err, compileErrSeen) {
		t.Fatalf("expected returned error to preserve compile error, got: %v", err)
	}

	if strings.Contains(err.Error(), "after compile hooks") {
		t.Fatalf("did not expect after compile hook label when after hook succeeds, got: %v", err)
	}
}

func TestEngineCompileReturnsAfterHookErrorOnSuccess(t *testing.T) {
	t.Parallel()

	afterErr := errors.New("after compile failed")
	var seenErr error
	eng := mustNewEngine(t, WithAfterCompileHook(func(_ context.Context, err error) error {
		seenErr = err
		return afterErr
	}))

	plan, err := eng.Compile(context.Background(), source.NewAnonymous(coverageValidQuery))
	if err == nil {
		t.Fatal("expected compile to fail when after hook fails")
	}

	if plan != nil {
		t.Fatal("expected plan to be nil when after hook fails")
	}

	if seenErr != nil {
		t.Fatalf("expected successful compile to pass nil run error to after hook, got: %v", seenErr)
	}

	if !errors.Is(err, afterErr) {
		t.Fatalf("expected after hook error to be preserved, got: %v", err)
	}

	if !strings.Contains(err.Error(), "after compile hooks") {
		t.Fatalf("expected after compile hooks label, got: %v", err)
	}
}

func TestEngineCompileJoinsCompileAndAfterHookErrors(t *testing.T) {
	t.Parallel()

	afterErr := errors.New("after compile failed")
	var compileErrSeen error

	eng := mustNewEngine(t, WithAfterCompileHook(func(_ context.Context, err error) error {
		compileErrSeen = err
		return afterErr
	}))

	plan, err := eng.Compile(context.Background(), source.NewAnonymous(coverageInvalidQuery))
	if err == nil {
		t.Fatal("expected compile to fail for invalid query and after hook error")
	}

	if plan != nil {
		t.Fatal("expected plan to be nil when compile/after hook fail")
	}

	if compileErrSeen == nil {
		t.Fatal("expected after hook to receive compile error")
	}

	if !errors.Is(err, afterErr) {
		t.Fatalf("expected joined error to include after hook error, got: %v", err)
	}

	if !strings.Contains(err.Error(), compileErrSeen.Error()) {
		t.Fatalf("expected joined error to include compile error details, got: %v", err)
	}

	if !strings.Contains(err.Error(), "after compile hooks") {
		t.Fatalf("expected joined error to include after compile hooks label, got: %v", err)
	}
}

func TestEngineRunReturnsCompileErrorWithoutPlanClose(t *testing.T) {
	t.Parallel()

	planCloseCalled := false
	eng := mustNewEngine(t, WithPlanCloseHook(func() error {
		planCloseCalled = true
		return nil
	}))

	result, err := eng.Run(context.Background(), source.NewAnonymous(coverageInvalidQuery))
	if err == nil {
		t.Fatal("expected run to fail when compile fails")
	}

	if result != nil {
		t.Fatal("expected nil result when compile fails")
	}

	if planCloseCalled {
		t.Fatal("plan close hooks should not run when plan was never created")
	}
}

func TestEngineLoadCreatesExecutablePlan(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t)
	data := mustMarshalArtifact(t, "RETURN 1 + 1")

	plan, err := eng.Load(data)
	if err != nil {
		t.Fatalf("expected load to succeed, got %v", err)
	}

	session := mustNewSession(t, plan)
	out, err := session.Run(context.Background())
	if err != nil {
		t.Fatalf("expected loaded plan to run, got %v", err)
	}

	if got, want := string(out.Content), "2"; got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestEngineLoadUsesConfiguredProgramLoader(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithProgramLoader(
		artifact.NewLoader(
			artifact.RegisteredFormat{ID: artifact.FormatJSON, Format: formatjson.Default},
		),
	))

	data := mustMarshalArtifact(t, "RETURN 1")

	plan, err := eng.Load(data)
	if err == nil {
		t.Fatal("expected load to fail when configured loader cannot decode default format")
	}

	if plan != nil {
		t.Fatal("expected plan to be nil when loader fails")
	}

	if !errors.Is(err, artifact.ErrUnknownFormat) {
		t.Fatalf("expected ErrUnknownFormat, got %v", err)
	}
}

func TestEngineLoadDoesNotRunCompileHooks(t *testing.T) {
	t.Parallel()

	beforeCalled := false
	afterCalled := false
	eng := mustNewEngine(t,
		WithBeforeCompileHook(func(context.Context) error {
			beforeCalled = true
			return errors.New("before compile should not run")
		}),
		WithAfterCompileHook(func(context.Context, error) error {
			afterCalled = true
			return errors.New("after compile should not run")
		}),
	)

	data := mustMarshalArtifact(t, "RETURN 1")

	plan, err := eng.Load(data)
	if err != nil {
		t.Fatalf("expected load to succeed without compile hooks, got %v", err)
	}

	if plan == nil {
		t.Fatal("expected load to return a plan")
	}

	if beforeCalled {
		t.Fatal("before compile hook must not run on load")
	}

	if afterCalled {
		t.Fatal("after compile hook must not run on load")
	}
}

func TestEngineLoadFailureDoesNotRunPlanCloseHooks(t *testing.T) {
	t.Parallel()

	planCloseCalled := false
	eng := mustNewEngine(t, WithPlanCloseHook(func() error {
		planCloseCalled = true
		return nil
	}))

	plan, err := eng.Load([]byte("bad artifact"))
	if err == nil {
		t.Fatal("expected load to fail for invalid artifact")
	}

	if plan != nil {
		t.Fatal("expected plan to be nil when load fails")
	}

	if planCloseCalled {
		t.Fatal("plan close hooks should not run when plan was never created")
	}
}

func TestEngineLoadPlanRunsPlanCloseHooks(t *testing.T) {
	t.Parallel()

	planCloseCalled := false
	eng := mustNewEngine(t, WithPlanCloseHook(func() error {
		planCloseCalled = true
		return nil
	}))

	data := mustMarshalArtifact(t, "RETURN 1")
	plan, err := eng.Load(data)
	if err != nil {
		t.Fatalf("expected load to succeed, got %v", err)
	}

	if err := plan.Close(); err != nil {
		t.Fatalf("expected close to succeed, got %v", err)
	}

	if !planCloseCalled {
		t.Fatal("expected plan close hook to run for loaded plan")
	}
}

func TestEngineNewReturnsProgramLoaderOptionError(t *testing.T) {
	t.Parallel()

	eng, err := New(WithProgramLoader(nil))
	if err == nil {
		t.Fatal("expected New to fail for nil program loader")
	}

	if eng != nil {
		t.Fatal("expected engine to be nil on program loader option error")
	}

	if !strings.Contains(err.Error(), "program loader cannot be nil") {
		t.Fatalf("expected program loader validation error, got: %v", err)
	}
}

func TestEngineClose(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t)
	if err := eng.Close(); err != nil {
		t.Fatalf("expected close without hooks to succeed, got: %v", err)
	}

	closeErr := errors.New("engine close failed")
	eng = mustNewEngine(t, WithEngineCloseHook(func() error {
		return closeErr
	}))

	err := eng.Close()
	if err == nil {
		t.Fatal("expected close to fail when close hook fails")
	}

	if !errors.Is(err, closeErr) {
		t.Fatalf("expected close error to include hook error, got: %v", err)
	}

	if !strings.Contains(err.Error(), "close hooks") {
		t.Fatalf("expected close hooks label, got: %v", err)
	}
}

func TestEngineParams(t *testing.T) {
	t.Parallel()

	eng, err := New(WithParams(map[string]any{
		"param1": "value1",
	}))

	if err != nil {
		t.Fatalf("expected New to succeed with valid params, got: %v", err)
	}

	if eng == nil {
		t.Fatal("expected engine to be non-nil on successful construction")
	}

	out, err := eng.Run(context.Background(), source.NewAnonymous("RETURN @param1"))

	if err != nil {
		t.Fatalf("expected run to succeed, got: %v", err)
	}

	var result any

	if err := json.Unmarshal(out.Content, &result); err != nil {
		t.Fatal("expected output to be valid JSON")
	}

	str, ok := result.(string)
	if !ok {
		t.Fatalf("expected result to be a string, got: %T", result)
	}

	if str != "value1" {
		t.Fatalf("expected run to return value1, got: %q", str)
	}
}

func TestEngineParam(t *testing.T) {
	t.Parallel()

	eng, err := New(WithParam("param1", "value1"))

	if err != nil {
		t.Fatalf("expected New to succeed with valid params, got: %v", err)
	}

	if eng == nil {
		t.Fatal("expected engine to be non-nil on successful construction")
	}

	out, err := eng.Run(context.Background(), source.NewAnonymous("RETURN @param1"))

	if err != nil {
		t.Fatalf("expected run to succeed, got: %v", err)
	}

	var result any

	if err := json.Unmarshal(out.Content, &result); err != nil {
		t.Fatal("expected output to be valid JSON")
	}

	str, ok := result.(string)
	if !ok {
		t.Fatalf("expected result to be a string, got: %T", result)
	}

	if str != "value1" {
		t.Fatalf("expected run to return value1, got: %q", str)
	}
}

func TestEngineRuntimeParams(t *testing.T) {
	t.Parallel()

	rtp, err := runtime.NewParamsFrom(map[string]any{
		"param1": "value1",
	})

	if err != nil {
		t.Fatalf("expected runtime.NewParamsFrom to succeed, got: %v", err)
	}

	eng, err := New(WithRuntimeParams(rtp))

	if err != nil {
		t.Fatalf("expected New to succeed with valid params, got: %v", err)
	}

	if eng == nil {
		t.Fatal("expected engine to be non-nil on successful construction")
	}

	out, err := eng.Run(context.Background(), source.NewAnonymous("RETURN @param1"))

	if err != nil {
		t.Fatalf("expected run to succeed, got: %v", err)
	}

	var result any

	if err := json.Unmarshal(out.Content, &result); err != nil {
		t.Fatal("expected output to be valid JSON")
	}

	str, ok := result.(string)
	if !ok {
		t.Fatalf("expected result to be a string, got: %T", result)
	}

	if str != "value1" {
		t.Fatalf("expected run to return value1, got: %q", str)
	}
}

func TestEngineRuntimeParam(t *testing.T) {
	t.Parallel()

	eng, err := New(WithRuntimeParam("param1", runtime.NewString("value1")))

	if err != nil {
		t.Fatalf("expected New to succeed with valid params, got: %v", err)
	}

	if eng == nil {
		t.Fatal("expected engine to be non-nil on successful construction")
	}

	out, err := eng.Run(context.Background(), source.NewAnonymous("RETURN @param1"))

	if err != nil {
		t.Fatalf("expected run to succeed, got: %v", err)
	}

	var result any

	if err := json.Unmarshal(out.Content, &result); err != nil {
		t.Fatal("expected output to be valid JSON")
	}

	str, ok := result.(string)
	if !ok {
		t.Fatalf("expected result to be a string, got: %T", result)
	}

	if str != "value1" {
		t.Fatalf("expected run to return value1, got: %q", str)
	}
}
