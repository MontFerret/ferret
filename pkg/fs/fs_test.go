package fs

import (
	"context"
	"errors"
	stdfs "io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestRootFSLstatDoesNotFollowFinalSymlink(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "target.txt"), []byte("target"), 0o600); err != nil {
		t.Fatalf("write target: %v", err)
	}
	if err := os.Symlink("target.txt", filepath.Join(root, "link.txt")); err != nil {
		t.Skipf("symlinks are unavailable: %v", err)
	}

	filesystem, err := New(WithRoot(root))
	if err != nil {
		t.Fatalf("create filesystem: %v", err)
	}

	metadata, err := ReaderFrom(WithFileSystem(context.Background(), filesystem))
	if err != nil {
		t.Fatalf("resolve link metadata: %v", err)
	}

	info, err := metadata.Lstat("link.txt")
	if err != nil {
		t.Fatalf("lstat link: %v", err)
	}
	if info.Mode()&stdfs.ModeSymlink == 0 {
		t.Fatalf("expected symlink metadata, got mode %v", info.Mode())
	}
}

func TestDisabledFSLstatPreservesRootDenial(t *testing.T) {
	filesystem, err := New()
	if err != nil {
		t.Fatalf("create disabled filesystem: %v", err)
	}
	metadata, err := ReaderFrom(WithFileSystem(context.Background(), filesystem))
	if err != nil {
		t.Fatalf("resolve disabled link metadata: %v", err)
	}

	if _, err := metadata.Lstat("file.txt"); !errors.Is(err, ErrRootNotConfigured) {
		t.Fatalf("expected ErrRootNotConfigured, got %v", err)
	}
}
