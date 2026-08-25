package sdktest

import (
	"context"
	"testing"

	apisource "github.com/MontFerret/api/source"
	ferret "github.com/MontFerret/ferret/v2"
)

// Harness owns an Engine configured for a module test.
type Harness struct {
	engine *ferret.Engine
}

// New creates a test engine and registers its cleanup with t.
func New(t testing.TB, options ...ferret.Option) *Harness {
	t.Helper()

	engine, err := ferret.New(options...)
	if err != nil {
		t.Fatalf("sdktest: create engine: %v", err)
	}

	harness := &Harness{engine: engine}
	t.Cleanup(func() {
		if err := engine.Close(); err != nil {
			t.Errorf("sdktest: close engine: %v", err)
		}
	})

	return harness
}

// Engine returns the harness engine for lower-level assertions.
func (h *Harness) Engine() *ferret.Engine {
	if h == nil {
		return nil
	}

	return h.engine
}

// Run compiles and executes query with a fresh session.
func (h *Harness) Run(
	ctx context.Context,
	query string,
	options ...ferret.SessionOption,
) (ferret.Output, error) {
	return h.engine.Run(ctx, apisource.File{
		Name:    "anonymous",
		Content: query,
	}, options...)
}
