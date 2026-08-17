package ferret

import (
	"errors"
	"strings"
	"testing"

	ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/module"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
)

type failingCloseFileSystem struct {
	ferretfs.FileSystem
	closeErr error
}

func (f *failingCloseFileSystem) Close() error {
	return errors.Join(f.closeErr, f.FileSystem.Close())
}

func TestEngineCloseClosesRootFileSystem(t *testing.T) {
	t.Parallel()

	engine, err := New(WithFSRoot(t.TempDir()))
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	filesystem := engine.host.fs

	if err := engine.Close(); err != nil {
		t.Fatalf("close engine: %v", err)
	}

	if _, err := filesystem.Stat("."); err == nil {
		t.Fatal("expected engine filesystem to be closed")
	}
}

func TestNewClosesRootFileSystemOnRegistrationFailure(t *testing.T) {
	t.Parallel()

	registerErr := errors.New("register failed")
	var filesystem ferretfs.FileSystem
	mod := testModule{
		registerFn: func(boot module.Bootstrap) error {
			filesystem = boot.Host().FileSystem()

			return registerErr
		},
	}

	_, err := New(WithFSRoot(t.TempDir()), WithModules(mod))
	if !errors.Is(err, registerErr) {
		t.Fatalf("expected registration error, got %v", err)
	}

	if _, err := filesystem.Stat("."); err == nil {
		t.Fatal("expected construction failure to close the filesystem")
	}
}

func TestNewClosesRootFileSystemOnHostBuildFailure(t *testing.T) {
	t.Parallel()

	var filesystem ferretfs.FileSystem
	mod := testModule{
		registerFn: func(boot module.Bootstrap) error {
			filesystem = boot.Host().FileSystem()
			boot.Host().Library().Function().A0().Add("FILESYSTEM_DUPLICATE_FN", testFn0)
			boot.Host().Library().Function().A0().Add("FILESYSTEM_DUPLICATE_FN", testFn0)

			return nil
		},
	}

	_, err := New(WithFSRoot(t.TempDir()), WithModules(mod))
	if err == nil {
		t.Fatal("expected host build failure")
	}

	if _, err := filesystem.Stat("."); err == nil {
		t.Fatal("expected host build failure to close the filesystem")
	}
}

func TestNewClosesRootFileSystemOnInitFailure(t *testing.T) {
	t.Parallel()

	initErr := errors.New("init failed")
	var filesystem ferretfs.FileSystem
	mod := testModule{
		registerFn: func(boot module.Bootstrap) error {
			filesystem = boot.Host().FileSystem()

			return nil
		},
	}

	_, err := New(
		WithFSRoot(t.TempDir()),
		WithModules(mod),
		WithEngineInitHook(func() error {
			return initErr
		}),
	)
	if !errors.Is(err, initErr) {
		t.Fatalf("expected init error, got %v", err)
	}

	if _, err := filesystem.Stat("."); err == nil {
		t.Fatal("expected init failure to close the filesystem")
	}
}

func TestNewJoinsConstructionHookAndFileSystemCloseErrors(t *testing.T) {
	t.Parallel()

	registerErr := errors.New("register failed")
	hookErr := errors.New("hook close failed")
	filesystemErr := errors.New("filesystem close failed")
	client := &recordingHTTPClient{}
	mod := testModule{
		registerFn: func(boot module.Bootstrap) error {
			internal, ok := boot.(*bootstrap)
			if !ok {
				t.Fatalf("expected internal bootstrap, got %T", boot)
			}

			internal.host.fs = &failingCloseFileSystem{
				FileSystem: internal.host.fs,
				closeErr:   filesystemErr,
			}
			internal.host.network = mustNewTestNetwork(t, ferretnet.WithHTTPClient(client))
			boot.Hooks().Engine().OnClose(func() error {
				return hookErr
			})

			return registerErr
		},
	}

	_, err := New(WithFSRoot(t.TempDir()), WithModules(mod))
	if !errors.Is(err, registerErr) {
		t.Fatalf("expected registration error, got %v", err)
	}

	if !errors.Is(err, hookErr) {
		t.Fatalf("expected hook close error, got %v", err)
	}

	if !errors.Is(err, filesystemErr) {
		t.Fatalf("expected filesystem close error, got %v", err)
	}

	if !strings.Contains(err.Error(), "close hooks") {
		t.Fatalf("expected close hooks label, got %v", err)
	}

	if !strings.Contains(err.Error(), "close filesystem") {
		t.Fatalf("expected close filesystem label, got %v", err)
	}

	if got := client.idleCloseCount(); got != 1 {
		t.Fatalf("expected network cleanup after close errors, got %d calls", got)
	}
}
