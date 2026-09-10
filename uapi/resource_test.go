package uapi

import (
	"context"
	"io"
	"path/filepath"
	"testing"
	"testing/synctest"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestSessionFilesystemOwnershipRemainsNative(t *testing.T) {
	for _, debug := range []bool{false, true} {
		for _, owned := range []bool{false, true} {
			name := map[bool]string{false: "ordinary", true: "debug"}[debug] + "/" + map[bool]string{false: "borrowed", true: "owned"}[owned]
			t.Run(name, func(t *testing.T) {
				var captured fs.FileSystem
				r := newTestRuntime(t, engine.WithFSRoot(t.TempDir()), engine.WithFunctionsRegistrar(func(ns runtime.Namespace) {
					ns.Function().A0().Add("CAPTURE_FS", func(ctx context.Context) (runtime.Value, error) {
						var err error
						captured, err = fs.FileSystemFrom(ctx)

						return runtime.True, err
					})
				}))
				p, err := r.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN CAPTURE_FS()"))
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = p.Close() })
				var options []api.SessionOption
				if owned {
					options = []api.SessionOption{api.WithFSRoot(t.TempDir())}
				}

				var closer io.Closer
				if debug {
					s, err := p.NewDebugSession(t.Context(), options...)
					if err != nil {
						t.Fatal(err)
					}

					closer = s
					if _, err := s.Start(t.Context()); err != nil {
						t.Fatal(err)
					}

					if _, err := s.Continue(t.Context()); err != nil {
						t.Fatal(err)
					}
				} else {
					s, err := p.NewSession(t.Context(), options...)
					if err != nil {
						t.Fatal(err)
					}

					closer = s
					if _, err := s.Run(t.Context()); err != nil {
						t.Fatal(err)
					}
				}

				if captured == nil {
					t.Fatal("Native filesystem did not reach execution")
				}

				for range 2 {
					if err := closer.Close(); err != nil {
						t.Fatal(err)
					}
				}

				if _, err := captured.Stat("."); (err != nil) != owned {
					t.Fatalf("filesystem err=%v, want closed=%t", err, owned)
				}
			})
		}
	}
}

func TestFailedSessionConstructionReleasesCapacity(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				r := newTestRuntime(t, engine.WithMaxActiveSessions(1))
				p, err := r.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN 1"))
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = p.Close() })
				option := api.WithFSRoot(filepath.Join(t.TempDir(), "missing"))
				var failed io.Closer
				if debug {
					failed, err = p.NewDebugSession(t.Context(), option)
				} else {
					failed, err = p.NewSession(t.Context(), option)
				}

				if failed != nil || err == nil {
					t.Fatalf("invalid root published session=%v err=%v", failed, err)
				}

				// This would durably deadlock in the bubble if rollback leaked the
				// sole Native permit.
				s, err := p.NewSession(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				if err := s.Close(); err != nil {
					t.Fatal(err)
				}
			})
		})
	}
}
