package fs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func mustNewFileSystem(t testing.TB, setters ...Option) FileSystem {
	t.Helper()

	filesystem, err := New(setters...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Cleanup(func() {
		if err := filesystem.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	})

	return filesystem
}

func TestFileSystemOptions(t *testing.T) {
	t.Run("defaults to disabled", func(t *testing.T) {
		filesystem := mustNewFileSystem(t)

		if filesystem != disabledFileSystem {
			t.Fatalf("filesystem = %T, want disabled filesystem", filesystem)
		}
	})

	t.Run("empty root is disabled", func(t *testing.T) {
		filesystem := mustNewFileSystem(t, WithRoot(""))

		if filesystem != disabledFileSystem {
			t.Fatalf("filesystem = %T, want disabled filesystem", filesystem)
		}
	})

	t.Run("valid root", func(t *testing.T) {
		filesystem := mustNewFileSystem(t, WithRoot(t.TempDir()))

		if _, err := filesystem.Stat("."); err != nil {
			t.Fatalf("Stat() error = %v", err)
		}
	})

	t.Run("read only false is writable", func(t *testing.T) {
		root := t.TempDir()
		filesystem := mustNewFileSystem(
			t,
			WithRoot(root),
			WithReadOnly(true),
			WithReadOnly(false),
		)

		if err := filesystem.WriteFile("file.txt", []byte("content"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		content, err := os.ReadFile(filepath.Join(root, "file.txt"))
		if err != nil {
			t.Fatalf("read written file: %v", err)
		}
		if string(content) != "content" {
			t.Fatalf("written content = %q, want %q", content, "content")
		}
	})

	t.Run("read only true rejects writes", func(t *testing.T) {
		filesystem := mustNewFileSystem(
			t,
			WithRoot(t.TempDir()),
			WithReadOnly(false),
			WithReadOnly(true),
		)

		err := filesystem.WriteFile("file.txt", []byte("content"), 0o600)
		if !errors.Is(err, ErrReadOnly) {
			t.Fatalf("WriteFile() error = %v, want %v", err, ErrReadOnly)
		}
	})

	t.Run("later root overrides earlier root", func(t *testing.T) {
		rootA := t.TempDir()
		rootB := t.TempDir()
		filesystem := mustNewFileSystem(
			t,
			WithRoot(rootA),
			WithRoot(rootB),
		)

		if err := filesystem.WriteFile("file.txt", []byte("root-b"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		if _, err := os.Stat(filepath.Join(rootA, "file.txt")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("first root Stat() error = %v, want %v", err, os.ErrNotExist)
		}

		content, err := os.ReadFile(filepath.Join(rootB, "file.txt"))
		if err != nil {
			t.Fatalf("read file from later root: %v", err)
		}
		if string(content) != "root-b" {
			t.Fatalf("later root content = %q, want %q", content, "root-b")
		}
	})
}

func TestNewReturnsOptionError(t *testing.T) {
	want := errors.New("option failed")
	invalid := func(_ *config) error {
		return want
	}

	filesystem, err := New(invalid)
	if !errors.Is(err, want) {
		t.Fatalf("New() error = %v, want %v", err, want)
	}
	if filesystem != nil {
		t.Fatalf("filesystem = %T, want nil", filesystem)
	}
}
