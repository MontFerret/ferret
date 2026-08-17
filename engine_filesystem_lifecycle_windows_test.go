//go:build windows

package ferret

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/module"
)

func TestEngineCloseReleasesRootDirectoryOnWindows(t *testing.T) {
	root := filepath.Join(t.TempDir(), "workspace")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}

	engine, err := New(WithFSRoot(root))
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}

	if err := engine.Close(); err != nil {
		t.Fatalf("close engine: %v", err)
	}

	if err := os.Remove(root); err != nil {
		t.Fatalf("remove root after engine close: %v", err)
	}
}

func TestNewReleasesRootDirectoryOnRegistrationFailureOnWindows(t *testing.T) {
	root := filepath.Join(t.TempDir(), "workspace")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatalf("create root: %v", err)
	}

	registerErr := errors.New("register failed")
	mod := testModule{
		registerFn: func(module.Bootstrap) error {
			return registerErr
		},
	}

	_, err := New(WithFSRoot(root), WithModules(mod))
	if !errors.Is(err, registerErr) {
		t.Fatalf("expected registration error, got %v", err)
	}

	if err := os.Remove(root); err != nil {
		t.Fatalf("remove root after construction failure: %v", err)
	}
}
