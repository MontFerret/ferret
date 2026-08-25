package ferret_test

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/MontFerret/api"
	apisource "github.com/MontFerret/api/source"
	ferret "github.com/MontFerret/ferret/v2"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	coreruntime "github.com/MontFerret/ferret/v2/pkg/runtime"
)

var (
	_ api.Runtime = (*ferret.Engine)(nil)
	_ api.Plan    = (*ferret.Plan)(nil)
	_ api.Session = (*ferret.Session)(nil)
)

type cleanupFailureValue struct {
	err error
}

func (v *cleanupFailureValue) Close() error {
	return v.err
}

func (v *cleanupFailureValue) MarshalJSON() ([]byte, error) {
	return []byte("1"), nil
}

func (v *cleanupFailureValue) String() string {
	return "cleanup-failure"
}

func (v *cleanupFailureValue) Hash() uint64 {
	return coreruntime.NewString(v.String()).Hash()
}

func (v *cleanupFailureValue) Copy() coreruntime.Value {
	return v
}

func TestUniversalRuntimeRunReturnsExactOutputAndCleansUp(t *testing.T) {
	t.Parallel()

	var planCloses atomic.Int32
	var sessionCloses atomic.Int32

	engine, err := ferret.New(
		ferret.WithPlanCloseHook(func() error {
			planCloses.Add(1)

			return nil
		}),
		ferret.WithSessionCloseHook(func() error {
			sessionCloses.Add(1)

			return nil
		}),
	)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	t.Cleanup(func() { _ = engine.Close() })

	var universal api.Runtime = engine
	output, err := universal.Run(t.Context(), apisource.File{
		Name:    "runtime-run.fql",
		Content: "RETURN @value",
	}, api.WithParam("value", "ok"))
	if err != nil {
		t.Fatalf("runtime run: %v", err)
	}

	if output.ContentType != encodingjson.ContentType {
		t.Fatalf("content type = %q, want %q", output.ContentType, encodingjson.ContentType)
	}

	if got, want := string(output.Content), `"ok"`; got != want {
		t.Fatalf("content = %q, want %q", got, want)
	}

	if got := sessionCloses.Load(); got != 1 {
		t.Fatalf("session close hooks = %d, want 1", got)
	}

	if got := planCloses.Load(); got != 1 {
		t.Fatalf("plan close hooks = %d, want 1", got)
	}
}

func TestUniversalPlanAndSessionAreReusableWithoutRecompilation(t *testing.T) {
	t.Parallel()

	var compiles atomic.Int32
	engine, err := ferret.New(ferret.WithBeforeCompileHook(func(context.Context) error {
		compiles.Add(1)

		return nil
	}))
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	t.Cleanup(func() { _ = engine.Close() })

	var universal api.Runtime = engine
	plan, err := universal.Compile(t.Context(), apisource.File{
		Name:    "reuse.fql",
		Content: "RETURN @value",
	})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	t.Cleanup(func() { _ = plan.Close() })

	for _, value := range []int{1, 2} {
		session, sessionErr := plan.NewSession(t.Context(), api.WithParam("value", value))
		if sessionErr != nil {
			t.Fatalf("new session: %v", sessionErr)
		}

		output, runErr := session.Run(t.Context())
		if runErr != nil {
			t.Fatalf("run session: %v", runErr)
		}

		if got, want := string(output.Content), strconv.Itoa(value); got != want {
			t.Fatalf("content = %q, want %q", got, want)
		}

		if closeErr := session.Close(); closeErr != nil {
			t.Fatalf("close session: %v", closeErr)
		}
	}

	session, err := plan.NewSession(t.Context(), api.WithParam("value", 3))
	if err != nil {
		t.Fatalf("new reusable session: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	for run := 0; run < 2; run++ {
		output, runErr := session.Run(t.Context())
		if runErr != nil {
			t.Fatalf("sequential run %d: %v", run+1, runErr)
		}

		if got := string(output.Content); got != "3" {
			t.Fatalf("sequential run %d content = %q, want 3", run+1, got)
		}
	}

	if got := compiles.Load(); got != 1 {
		t.Fatalf("compile hooks = %d, want 1", got)
	}
}

func TestUniversalSessionOptionsCancellationAndClosedStates(t *testing.T) {
	t.Parallel()

	engine, err := ferret.New()
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	t.Cleanup(func() { _ = engine.Close() })

	var universal api.Runtime = engine
	plan, err := universal.Compile(t.Context(), apisource.File{
		Name:    "options.fql",
		Content: "RETURN [@first, @second]",
	})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	session, err := plan.NewSession(
		t.Context(),
		nil,
		api.WithParams(map[string]any{"first": 1}),
		api.WithParam("second", 2),
		api.WithOutputContentType(encodingjson.ContentType),
	)
	if err != nil {
		t.Fatalf("new session: %v", err)
	}

	output, err := session.Run(t.Context())
	if err != nil {
		t.Fatalf("run: %v", err)
	}

	if got, want := string(output.Content), "[1,2]"; got != want {
		t.Fatalf("content = %q, want %q", got, want)
	}

	if err := session.Close(); err != nil {
		t.Fatalf("close session: %v", err)
	}

	closedOutput, err := session.Run(t.Context())
	if !errors.Is(err, coreruntime.ErrInvalidOperation) {
		t.Fatalf("closed session error = %v, want invalid operation", err)
	}

	if closedOutput.ContentType != "" || closedOutput.Content != nil {
		t.Fatalf("closed session output = %#v, want zero", closedOutput)
	}

	if err := plan.Close(); err != nil {
		t.Fatalf("close plan: %v", err)
	}

	closedSession, err := plan.NewSession(t.Context())
	if !errors.Is(err, coreruntime.ErrInvalidOperation) {
		t.Fatalf("closed plan error = %v, want invalid operation", err)
	}

	if closedSession != nil {
		_ = closedSession.Close()

		t.Fatalf("closed plan session = %T, want nil", closedSession)
	}

	cancelCtx, cancel := context.WithCancel(t.Context())
	cancel()

	cancelledOutput, err := universal.Run(cancelCtx, apisource.File{
		Name:    "cancel.fql",
		Content: "RETURN WAITFOR false TIMEOUT 10s EVERY 10s",
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled run error = %v, want context canceled", err)
	}

	if cancelledOutput.ContentType != "" || cancelledOutput.Content != nil {
		t.Fatalf("cancelled output = %#v, want zero", cancelledOutput)
	}
}

func TestUniversalSessionPreservesOutputWithCleanupError(t *testing.T) {
	t.Parallel()

	closeErr := errors.New("cleanup failed")
	value := &cleanupFailureValue{err: closeErr}
	engine, err := ferret.New(ferret.WithFunctionsRegistrar(func(namespace coreruntime.Namespace) {
		namespace.Function().A0().Add(
			"MAKE_CLEANUP_FAILURE",
			func(context.Context) (coreruntime.Value, error) {
				return value, nil
			},
		)
	}))
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	t.Cleanup(func() { _ = engine.Close() })

	var universal api.Runtime = engine
	plan, err := universal.Compile(t.Context(), apisource.File{
		Name:    "cleanup.fql",
		Content: "RETURN MAKE_CLEANUP_FAILURE()",
	})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	t.Cleanup(func() { _ = plan.Close() })

	session, err := plan.NewSession(t.Context())
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	output, err := session.Run(t.Context())
	if !errors.Is(err, closeErr) {
		t.Fatalf("run error = %v, want cleanup error", err)
	}

	if output.ContentType != encodingjson.ContentType || string(output.Content) != "1" {
		t.Fatalf("output = %#v, want materialized JSON 1", output)
	}
}

type foreignSessionOptions struct{}

func (foreignSessionOptions) SetParam(string, any) error {
	return nil
}

func (foreignSessionOptions) SetParams(map[string]any) error {
	return nil
}

func (foreignSessionOptions) SetOutputContentType(string) error {
	return nil
}

func TestFerretSessionOptionsRejectForeignRuntimeTargets(t *testing.T) {
	t.Parallel()

	err := ferret.WithSessionParam("value", 1)(foreignSessionOptions{})
	if err == nil || !strings.Contains(err.Error(), "Ferret session option cannot be applied") {
		t.Fatalf("foreign target error = %v", err)
	}
}
