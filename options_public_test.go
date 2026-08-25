package ferret_test

import (
	"context"
	"testing"

	ferret "github.com/MontFerret/ferret/v2"
)

func TestPublicOptionTypesRemainUsable(t *testing.T) {
	t.Parallel()

	engineOptions := []ferret.Option{
		ferret.WithMaxActiveSessions(1),
		ferret.WithMaxIdleVMsPerPlan(0),
		ferret.WithMaxVMsPerPlan(1),
	}

	engine, err := ferret.New(engineOptions...)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	defer func() { _ = engine.Close() }()

	plan, err := engine.Compile(context.Background(), ferret.NewAnonymousSource("RETURN 1"))
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	defer func() { _ = plan.Close() }()

	sessionOptions := []ferret.SessionOption{
		ferret.WithOutputContentType("application/json"),
		ferret.WithParam("request", "external"),
	}

	session, err := plan.NewSession(context.Background(), sessionOptions...)
	if err != nil {
		t.Fatalf("new session: %v", err)
	}

	if err := session.Close(); err != nil {
		t.Fatalf("close session: %v", err)
	}
}
