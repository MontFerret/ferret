package fs

import (
	"context"
	"errors"
	"testing"
)

func TestWithFileSystemRoundTrip(t *testing.T) {
	filesystem := disabledFileSystem
	ctx := WithFileSystem(context.Background(), filesystem)

	resolved, err := FileSystemFrom(ctx)
	if err != nil {
		t.Fatalf("filesystem from context failed: %v", err)
	}

	if resolved != filesystem {
		t.Fatalf("expected same filesystem instance from context")
	}
}

func TestContextResolvers(t *testing.T) {
	filesystem := disabledFileSystem
	ctx := WithFileSystem(context.Background(), filesystem)

	if resolved, err := FileSystemFrom(ctx); err != nil {
		t.Fatalf("filesystem from context failed: %v", err)
	} else if resolved != filesystem {
		t.Fatalf("expected same filesystem from FileSystemFrom")
	}

	if resolved, err := ReaderFrom(ctx); err != nil {
		t.Fatalf("reader from context failed: %v", err)
	} else if resolved != filesystem {
		t.Fatalf("expected same filesystem from ReaderFrom")
	}

	if resolved, err := WriterFrom(ctx); err != nil {
		t.Fatalf("writer from context failed: %v", err)
	} else if resolved != filesystem {
		t.Fatalf("expected same filesystem from WriterFrom")
	}

	if resolved, err := DirectoriesFrom(ctx); err != nil {
		t.Fatalf("directories from context failed: %v", err)
	} else if resolved != filesystem {
		t.Fatalf("expected same filesystem from DirectoriesFrom")
	}

	if resolved, err := RemoverFrom(ctx); err != nil {
		t.Fatalf("remover from context failed: %v", err)
	} else if resolved != filesystem {
		t.Fatalf("expected same filesystem from RemoverFrom")
	}
}

func TestFileSystemFromError(t *testing.T) {
	if _, err := FileSystemFrom(context.Background()); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for background context, got %v", err)
	}

	if _, err := FileSystemFrom(nil); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound for nil context, got %v", err)
	}
}
