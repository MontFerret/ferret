package fs

import "context"

type fsContextKey struct{}

var fsCtxKey = fsContextKey{}

// WithFileSystem adds FileSystem to context.
func WithFileSystem(ctx context.Context, registry FileSystem) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, fsCtxKey, registry)
}

// FileSystemFrom gets FileSystem from context.
func FileSystemFrom(ctx context.Context) (FileSystem, error) {
	val := ctx.Value(fsCtxKey)
	if val == nil {
		return nil, ErrNotFound
	}

	fs, ok := val.(FileSystem)

	if !ok {
		return nil, ErrNotFound
	}

	return fs, nil
}

// ReaderFrom gets Reader from context.
func ReaderFrom(ctx context.Context) (Reader, error) {
	return FileSystemFrom(ctx)
}

// DirectoriesFrom gets Directories from context.
func DirectoriesFrom(ctx context.Context) (Directories, error) {
	return FileSystemFrom(ctx)
}

// WriterFrom gets Writer from context.
func WriterFrom(ctx context.Context) (Writer, error) {
	return FileSystemFrom(ctx)
}

// RemoverFrom gets Remover from context.
func RemoverFrom(ctx context.Context) (Remover, error) {
	return FileSystemFrom(ctx)
}
