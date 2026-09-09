package engine

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"

	ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestSessionHostResourcesAreIsolated(t *testing.T) {
	for _, debug := range []bool{false, true} {
		for _, owned := range []bool{false, true} {
			name := map[bool]string{false: "ordinary", true: "debug"}[debug] + "/" + map[bool]string{false: "borrowed", true: "owned"}[owned]
			t.Run(name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					client := &recordingHTTPClient{}
					var capturedFS ferretfs.FileSystem
					engine := mustNewEngine(t, withCapturedSessionFileSystem(&capturedFS), WithFSRoot(t.TempDir()), WithNetworkOptions(ferretnet.WithHTTPClient(client)), WithMaxActiveSessions(2))
					t.Cleanup(func() { _ = engine.Close() })
					plan, err := engine.CompileDebug(t.Context(), source.NewAnonymous("RETURN CAPTURE_FS()"))
					if err != nil {
						t.Fatal(err)
					}

					t.Cleanup(func() { _ = plan.Close() })
					sibling := mustNewSession(t, plan, WithSessionFSRoot(t.TempDir()))
					t.Cleanup(func() { _ = sibling.Close() })
					runNativeLifecycleSession(t, sibling)
					siblingFS := capturedFS
					var setters []SessionOption

					if owned {
						setters = []SessionOption{WithSessionFSRoot(t.TempDir())}
					}

					session, err := createNativeLifecycleSession(plan, t.Context(), debug, setters...)
					if err != nil {
						t.Fatal(err)
					}

					t.Cleanup(func() { _ = session.Close() })
					runNativeLifecycleSession(t, session)
					filesystem := capturedFS

					for range 2 {
						if err := session.Close(); err != nil {
							t.Fatal(err)
						}
					}

					if _, err := filesystem.Stat("."); (err != nil) != owned {
						t.Fatalf("session filesystem stat = %v, want closed=%t", err, owned)
					}

					if _, err := engine.host.FileSystem.Stat("."); err != nil {
						t.Fatalf("session closed the Engine filesystem: %v", err)
					}

					if _, err := siblingFS.Stat("."); err != nil {
						t.Fatalf("session closed its sibling's filesystem: %v", err)
					}

					if client.idleCloseCount() != 0 {
						t.Fatal("session closed the Engine network")
					}

					replacement := mustNewSession(t, plan)
					t.Cleanup(func() { _ = replacement.Close() })
					assertSessionAdmissionBlocked(t, plan)
					if err := replacement.Close(); err != nil {
						t.Fatal(err)
					}

					if err := sibling.Close(); err != nil {
						t.Fatal(err)
					}

					if err := plan.Close(); err != nil {
						t.Fatal(err)
					}

					if err := engine.Close(); err != nil {
						t.Fatal(err)
					}

					if client.idleCloseCount() != 1 {
						t.Fatal("Engine did not retain network ownership")
					}
				})
			})
		}
	}
}

// assertSessionAdmissionBlocked must run inside a synctest bubble. It observes
// admission through Plan rather than exposing the limiter's channel to consumers.
func assertSessionAdmissionBlocked(t *testing.T, plan *Plan) {
	t.Helper()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var child *Session
	var err error
	finished := false

	go func() {
		child, err = plan.NewSession(ctx)
		finished = true
	}()

	synctest.Wait()
	if finished {
		if child != nil {
			_ = child.Close()
		}

		t.Errorf("admission bypassed occupied capacity: session=%v error=%v", child, err)

		return
	}

	cancel()
	synctest.Wait()
	if !finished || child != nil || !errors.Is(err, context.Canceled) {
		t.Errorf("canceled admission: finished=%t session=%v error=%v", finished, child, err)
	}
}
