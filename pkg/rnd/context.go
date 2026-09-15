package rnd

import "context"

// WithContext transports an already-owned source to execution-time code.
// The caller supplies a non-nil context and retains ownership of src.
func WithContext(ctx context.Context, src *Source) context.Context {
	return context.WithValue(ctx, contextKey{}, src)
}

// FromContext returns the transported source. Missing and nil sources return
// false; lookup never constructs a source or advances a sequence.
func FromContext(ctx context.Context) (*Source, bool) {
	if ctx == nil {
		return nil, false
	}

	src, ok := ctx.Value(contextKey{}).(*Source)

	return src, ok && src != nil
}

type contextKey struct{}
