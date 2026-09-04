package ferret_test

import (
	"context"
	"testing"

	ferret "github.com/MontFerret/ferret/v2"
)

func TestPublicOptionTypesRemainUsable(t *testing.T) {
	t.Parallel()

	levels := []ferret.LogLevel{
		ferret.LogTrace,
		ferret.LogDebug,
		ferret.LogInfo,
		ferret.LogWarn,
		ferret.LogError,
		ferret.LogFatal,
		ferret.LogPanic,
		ferret.LogNone,
		ferret.LogDisabled,
	}
	if len(levels) != 9 {
		t.Fatalf("log levels = %d, want 9", len(levels))
	}

	level, err := ferret.ParseLogLevel("info")
	if err != nil {
		t.Fatalf("parse log level: %v", err)
	}
	if level != ferret.LogInfo {
		t.Fatalf("parsed log level = %v, want %v", level, ferret.LogInfo)
	}
	if level := ferret.MustParseLogLevel("debug"); level != ferret.LogDebug {
		t.Fatalf("must-parse log level = %v, want %v", level, ferret.LogDebug)
	}

	var mod ferret.Module
	var moduleOption ferret.Option = ferret.WithModules(mod)
	if moduleOption == nil {
		t.Fatal("module option is nil")
	}

	var engineInitHook ferret.EngineInitHook = func() error {
		return nil
	}
	var engineCloseHook ferret.EngineCloseHook = func() error {
		return nil
	}
	var beforeCompileHook ferret.BeforeCompileHook = func(context.Context) error {
		return nil
	}
	var afterCompileHook ferret.AfterCompileHook = func(context.Context, error) error {
		return nil
	}
	var planCloseHook ferret.PlanCloseHook = func() error {
		return nil
	}
	var beforeRunHook ferret.BeforeRunHook = func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	}
	var afterRunHook ferret.AfterRunHook = func(context.Context, error) error {
		return nil
	}
	var sessionCloseHook ferret.SessionCloseHook = func() error {
		return nil
	}

	params := ferret.Params{}
	if err := params.Set("engineParam", 1); err != nil {
		t.Fatalf("set engine param: %v", err)
	}
	var value ferret.Value = params.MustGet("engineParam")

	engineOptions := []ferret.Option{
		ferret.WithOptimizationLevel(ferret.OptimizationBasic),
		ferret.WithMaxActiveSessions(1),
		ferret.WithMaxIdleVMsPerPlan(0),
		ferret.WithMaxVMsPerPlan(1),
		ferret.WithLogLevel(ferret.LogInfo),
		ferret.WithRuntimeParams(params),
		ferret.WithRuntimeParam("engineValue", value),
		ferret.WithEngineInitHook(engineInitHook),
		ferret.WithEngineCloseHook(engineCloseHook),
		ferret.WithBeforeCompileHook(beforeCompileHook),
		ferret.WithAfterCompileHook(afterCompileHook),
		ferret.WithPlanCloseHook(planCloseHook),
		ferret.WithBeforeRunHook(beforeRunHook),
		ferret.WithAfterRunHook(afterRunHook),
		ferret.WithSessionCloseHook(sessionCloseHook),
	}

	engine, err := ferret.New(engineOptions...)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	t.Cleanup(func() {
		if err := engine.Close(); err != nil {
			t.Errorf("close engine: %v", err)
		}
	})

	plan, err := engine.Compile(
		context.Background(),
		ferret.NewAnonymousSource("RETURN @engineParam + @engineValue + @sessionParam + @sessionValue"),
	)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	t.Cleanup(func() {
		if err := plan.Close(); err != nil {
			t.Errorf("close plan: %v", err)
		}
	})

	sessionParams := ferret.Params{}
	if err := sessionParams.Set("sessionParam", 3); err != nil {
		t.Fatalf("set session param: %v", err)
	}
	sessionValue := sessionParams.MustGet("sessionParam")

	sessionOptions := []ferret.SessionOption{
		ferret.WithOutputContentType("application/json"),
		ferret.WithSessionFSRoot(t.TempDir()),
		ferret.WithSessionLogLevel(ferret.LogDebug),
		ferret.WithSessionRuntimeParams(sessionParams),
		ferret.WithSessionRuntimeParam("sessionValue", sessionValue),
	}

	session, err := plan.NewSession(context.Background(), sessionOptions...)
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	t.Cleanup(func() {
		if err := session.Close(); err != nil {
			t.Errorf("close session: %v", err)
		}
	})

	output, err := session.Run(context.Background())
	if err != nil {
		t.Fatalf("run session: %v", err)
	}
	if got := string(output.Content); got != "8" {
		t.Fatalf("session output = %q, want %q", got, "8")
	}

	if err := session.Close(); err != nil {
		t.Fatalf("close session: %v", err)
	}

	var format ferret.ProgramFormat = ferret.ProgramFormatJSON
	var programOption ferret.ProgramOption = ferret.WithProgramFormat(format)

	data, err := plan.Marshal(programOption)
	if err != nil {
		t.Fatalf("marshal plan: %v", err)
	}

	program, err := ferret.UnmarshalProgram(data)
	if err != nil {
		t.Fatalf("unmarshal program: %v", err)
	}

	var msgpackOption ferret.ProgramOption = ferret.WithProgramFormat(ferret.ProgramFormatMsgPack)
	data, err = ferret.MarshalProgram(program, msgpackOption)
	if err != nil {
		t.Fatalf("marshal program: %v", err)
	}

	loadedPlan, err := engine.Load(data)
	if err != nil {
		t.Fatalf("load plan: %v", err)
	}
	t.Cleanup(func() {
		if err := loadedPlan.Close(); err != nil {
			t.Errorf("close loaded plan: %v", err)
		}
	})

	loadedSession, err := loadedPlan.NewSession(context.Background(), sessionOptions...)
	if err != nil {
		t.Fatalf("new loaded session: %v", err)
	}
	t.Cleanup(func() {
		if err := loadedSession.Close(); err != nil {
			t.Errorf("close loaded session: %v", err)
		}
	})

	output, err = loadedSession.Run(context.Background())
	if err != nil {
		t.Fatalf("run loaded session: %v", err)
	}
	if got := string(output.Content); got != "8" {
		t.Fatalf("loaded session output = %q, want %q", got, "8")
	}
}
