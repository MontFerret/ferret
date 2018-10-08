package env

import "context"

type (
	ctxKey struct{}

	Environment struct {
		CDPAddress   string
		ProxyAddress string
		UserAgent    string
	}
)

const RandomUserAgent = "*"

func WithContext(ctx context.Context, e Environment) context.Context {
	return context.WithValue(ctx, ctxKey{}, e)
}

func FromContext(ctx context.Context) Environment {
	res := ctx.Value(ctxKey{})

	val, ok := res.(Environment)

	if !ok {
		return Environment{}
	}

	return val
}
