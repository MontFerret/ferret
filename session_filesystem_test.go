package ferret

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestSessionFSRootOverridesEngineFSRoot(t *testing.T) {
	t.Parallel()

	engineRoot := t.TempDir()
	sessionRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(engineRoot, "value.txt"), []byte("engine"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sessionRoot, "value.txt"), []byte("session"), 0o600); err != nil {
		t.Fatal(err)
	}

	engine := mustNewEngine(t, WithFSRoot(engineRoot))
	defer func() { _ = engine.Close() }()
	plan := mustCompilePlan(t, engine, `RETURN TO_STRING(IO::FS::READ("value.txt"))`)
	defer func() { _ = plan.Close() }()

	inherited := mustNewSession(t, plan)
	inheritedOutput, err := inherited.Run(context.Background())
	if err != nil {
		t.Fatalf("run inherited filesystem session: %v", err)
	}
	if err := inherited.Close(); err != nil {
		t.Fatalf("close inherited filesystem session: %v", err)
	}
	if got := string(inheritedOutput.Content); got != `"engine"` {
		t.Fatalf("inherited filesystem output = %s, want %q", got, "engine")
	}

	overridden := mustNewSession(t, plan, WithSessionFSRoot(sessionRoot))
	overriddenOutput, err := overridden.Run(context.Background())
	if err != nil {
		t.Fatalf("run overridden filesystem session: %v", err)
	}
	if err := overridden.Close(); err != nil {
		t.Fatalf("close overridden filesystem session: %v", err)
	}
	if got := string(overriddenOutput.Content); got != `"session"` {
		t.Fatalf("overridden filesystem output = %s, want %q", got, "session")
	}
}

func TestSessionFSRootWritesOutsideEngineFSRoot(t *testing.T) {
	t.Parallel()

	engineRoot := t.TempDir()
	sessionRoot := t.TempDir()
	engine := mustNewEngine(t, WithFSRoot(engineRoot))
	defer func() { _ = engine.Close() }()
	plan := mustCompilePlan(t, engine, `RETURN IO::FS::WRITE("created.txt", TO_BINARY("session"))`)
	defer func() { _ = plan.Close() }()

	session := mustNewSession(t, plan, WithSessionFSRoot(sessionRoot))
	if _, err := session.Run(context.Background()); err != nil {
		t.Fatalf("run session with overridden filesystem: %v", err)
	}
	if err := session.Close(); err != nil {
		t.Fatalf("close session with overridden filesystem: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(sessionRoot, "created.txt"))
	if err != nil {
		t.Fatalf("read session-root output: %v", err)
	}
	if string(content) != "session" {
		t.Fatalf("session-root output = %q, want %q", content, "session")
	}
	if _, err := os.Stat(filepath.Join(engineRoot, "created.txt")); !os.IsNotExist(err) {
		t.Fatalf("session write escaped to the engine root: %v", err)
	}
}

func TestSessionFSRootPreservesEngineReadOnlyPolicy(t *testing.T) {
	t.Parallel()

	engine := mustNewEngine(t, WithFSRoot(t.TempDir()), WithFSReadOnly())
	defer func() { _ = engine.Close() }()
	plan := mustCompilePlan(t, engine, `RETURN IO::FS::WRITE("created.txt", TO_BINARY("value"))`)
	defer func() { _ = plan.Close() }()

	sessionRoot := t.TempDir()
	session := mustNewSession(t, plan, WithSessionFSRoot(sessionRoot))
	defer func() { _ = session.Close() }()

	if _, err := session.Run(context.Background()); !errors.Is(err, ferretfs.ErrReadOnly) {
		t.Fatalf("run error = %v, want %v", err, ferretfs.ErrReadOnly)
	}
	if _, err := os.Stat(filepath.Join(sessionRoot, "created.txt")); !os.IsNotExist(err) {
		t.Fatalf("read-only session created a file: %v", err)
	}
}

func TestConcurrentSessionsUseIndependentFSRoots(t *testing.T) {
	t.Parallel()

	engine := mustNewEngine(t, WithFSRoot(t.TempDir()))
	defer func() { _ = engine.Close() }()
	plan := mustCompilePlan(t, engine, `RETURN TO_STRING(IO::FS::READ("value.txt"))`)
	defer func() { _ = plan.Close() }()

	roots := []string{t.TempDir(), t.TempDir()}
	for i, value := range []string{"first", "second"} {
		if err := os.WriteFile(filepath.Join(roots[i], "value.txt"), []byte(value), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	sessions := []*Session{
		mustNewSession(t, plan, WithSessionFSRoot(roots[0])),
		mustNewSession(t, plan, WithSessionFSRoot(roots[1])),
	}
	defer func() {
		for _, session := range sessions {
			_ = session.Close()
		}
	}()

	outputs := make([]string, len(sessions))
	errs := make([]error, len(sessions))
	var wait sync.WaitGroup
	for i, session := range sessions {
		wait.Add(1)
		go func(index int, current *Session) {
			defer wait.Done()
			output, err := current.Run(context.Background())
			errs[index] = err
			if output != nil {
				outputs[index] = string(output.Content)
			}
		}(i, session)
	}
	wait.Wait()

	for i, want := range []string{`"first"`, `"second"`} {
		if errs[i] != nil {
			t.Fatalf("session %d error = %v", i, errs[i])
		}
		if outputs[i] != want {
			t.Fatalf("session %d output = %s, want %s", i, outputs[i], want)
		}
	}
}

func TestSessionClosesOwnedFSRootWithoutClosingBorrowedRoot(t *testing.T) {
	t.Parallel()

	engineRoot := t.TempDir()
	engine := mustNewEngine(t, WithFSRoot(engineRoot))
	defer func() { _ = engine.Close() }()
	plan := mustCompilePlan(t, engine, coverageValidQuery)
	defer func() { _ = plan.Close() }()

	borrowed := mustNewSession(t, plan)
	if err := borrowed.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := engine.host.fs.Stat("."); err != nil {
		t.Fatalf("closing a default session closed the engine filesystem: %v", err)
	}

	owned := mustNewSession(t, plan, WithSessionFSRoot(t.TempDir()))
	ownedFS := owned.fs
	if err := owned.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := ownedFS.Stat("."); err == nil {
		t.Fatal("closing a session did not close its owned filesystem")
	}
}

func TestSessionCloseJoinsOwnedFileSystemErrorExactlyOnce(t *testing.T) {
	t.Parallel()

	hookErr := errors.New("session hook close failed")
	fileSystemErr := errors.New("session filesystem close failed")
	engine := mustNewEngine(t, WithSessionCloseHook(func() error { return hookErr }))
	defer func() { _ = engine.Close() }()
	plan := mustCompilePlan(t, engine, coverageValidQuery)
	defer func() { _ = plan.Close() }()
	session := mustNewSession(t, plan, WithSessionFSRoot(t.TempDir()))

	filesystem := &countingCloseFileSystem{FileSystem: session.fs, closeErr: fileSystemErr}
	session.fs = filesystem
	firstErr := session.Close()
	secondErr := session.Close()

	if !errors.Is(firstErr, hookErr) || !errors.Is(firstErr, fileSystemErr) {
		t.Fatalf("first close error = %v, want hook and filesystem failures", firstErr)
	}
	if firstErr != secondErr {
		t.Fatalf("repeated close returned different errors: %v and %v", firstErr, secondErr)
	}
	if calls := filesystem.closeCalls.Load(); calls != 1 {
		t.Fatalf("filesystem close calls = %d, want 1", calls)
	}
}

func TestDebugSessionUsesAndClosesOwnedFSRoot(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	fileSystemErr := errors.New("debug session filesystem close failed")
	if err := os.WriteFile(filepath.Join(root, "value.txt"), []byte("debug"), 0o600); err != nil {
		t.Fatal(err)
	}

	engine := mustNewEngine(t, WithFSRoot(t.TempDir()))
	defer func() { _ = engine.Close() }()
	plan, err := engine.CompileDebug(
		context.Background(),
		source.NewAnonymous(`RETURN TO_STRING(IO::FS::READ("value.txt"))`),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = plan.Close() }()

	var filesystem *countingCloseFileSystem
	session, err := newPlanSession(
		plan,
		context.Background(),
		[]SessionOption{WithSessionFSRoot(root)},
		planSessionSetup{requiresDebugInfo: true},
		func(dependencies planSessionDependencies) (*DebugSession, error) {
			filesystem = &countingCloseFileSystem{
				FileSystem: dependencies.filesystem,
				closeErr:   fileSystemErr,
			}
			dependencies.filesystem = filesystem

			return buildDebugSession(dependencies)
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := session.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	event, err := session.Continue(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if event.Output == nil || string(event.Output.Content) != `"debug"` {
		t.Fatalf("debug output = %#v, want %q", event.Output, "debug")
	}
	firstErr := session.Close()
	secondErr := session.Close()
	if !errors.Is(firstErr, fileSystemErr) {
		t.Fatalf("first debug close error = %v, want %v", firstErr, fileSystemErr)
	}
	if firstErr != secondErr {
		t.Fatalf("repeated debug close returned different errors: %v and %v", firstErr, secondErr)
	}
	if calls := filesystem.closeCalls.Load(); calls != 1 {
		t.Fatalf("debug filesystem close calls = %d, want 1", calls)
	}
}

func TestNewPlanSessionClosesOwnedFSRootOnBuildFailure(t *testing.T) {
	t.Parallel()

	engine := mustNewEngine(t)
	defer func() { _ = engine.Close() }()
	plan := mustCompilePlan(t, engine, coverageValidQuery)
	defer func() { _ = plan.Close() }()

	buildErr := errors.New("session build failed")
	var filesystem ferretfs.FileSystem
	_, err := newPlanSession(
		plan,
		context.Background(),
		[]SessionOption{WithSessionFSRoot(t.TempDir())},
		planSessionSetup{},
		func(dependencies planSessionDependencies) (struct{}, error) {
			filesystem = dependencies.filesystem

			return struct{}{}, buildErr
		},
	)
	if !errors.Is(err, buildErr) {
		t.Fatalf("session construction error = %v, want %v", err, buildErr)
	}
	if _, err := filesystem.Stat("."); err == nil {
		t.Fatal("session construction failure did not close its owned filesystem")
	}
}

func TestSessionFSRootWaitsForPermitBeforeOpening(t *testing.T) {
	t.Parallel()

	engine := mustNewEngine(t, WithMaxActiveSessions(1))
	defer func() { _ = engine.Close() }()
	plan := mustCompilePlan(t, engine, coverageValidQuery)
	defer func() { _ = plan.Close() }()

	active := mustNewSession(t, plan)
	defer func() { _ = active.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	missingRoot := filepath.Join(t.TempDir(), "missing")
	session, err := plan.NewSession(ctx, WithSessionFSRoot(missingRoot))
	if session != nil {
		_ = session.Close()
		t.Fatal("canceled session creation returned a session")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("session creation error = %v, want %v", err, context.Canceled)
	}
}

func TestSessionFSRootCreationFailureReleasesPermit(t *testing.T) {
	t.Parallel()

	engine := mustNewEngine(t, WithMaxActiveSessions(1))
	defer func() { _ = engine.Close() }()
	plan := mustCompilePlan(t, engine, coverageValidQuery)
	defer func() { _ = plan.Close() }()

	missingRoot := filepath.Join(t.TempDir(), "missing")
	if session, err := plan.NewSession(context.Background(), WithSessionFSRoot(missingRoot)); err == nil {
		_ = session.Close()
		t.Fatal("session creation with a missing filesystem root unexpectedly succeeded")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	session, err := plan.NewSession(ctx)
	if err != nil {
		t.Fatalf("session creation after filesystem failure: %v", err)
	}
	if err := session.Close(); err != nil {
		t.Fatalf("close session after filesystem failure: %v", err)
	}
}

func TestSessionFSRootRejectsUnusablePaths(t *testing.T) {
	t.Parallel()

	engine := mustNewEngine(t)
	defer func() { _ = engine.Close() }()
	plan := mustCompilePlan(t, engine, coverageValidQuery)
	defer func() { _ = plan.Close() }()

	file := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(file, []byte("value"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, root := range []string{filepath.Join(t.TempDir(), "missing"), file} {
		if session, err := plan.NewSession(context.Background(), WithSessionFSRoot(root)); err == nil {
			_ = session.Close()
			t.Fatalf("NewSession(%q) unexpectedly succeeded", root)
		}
	}
}
