package crypto

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const tokenAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// random_token generates a cryptographically secure token with uniformly sampled lowercase letters, uppercase letters, and digits.
// @param length {Int} Token length, from 1 through 65536 inclusive.
// @return {String} The generated token. Entropy-source errors are propagated.
func RandomToken(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	length, err := runtime.CastArg[runtime.Int](arg, 0)
	if err != nil {
		return runtime.None, err
	}

	if length < 1 || length > 65536 {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "length must be between 1 and 65536"), 0)
	}

	return randomToken(ctx, rand.Reader, int(length))
}

func randomToken(ctx context.Context, reader io.Reader, length int) (runtime.Value, error) {
	out := make([]byte, length)
	var buffer [256]byte
	for written := 0; written < length; {
		if err := ctx.Err(); err != nil {
			return runtime.None, err
		}

		batch := buffer[:min(len(buffer), length-written)]
		if _, err := io.ReadFull(reader, batch); err != nil {
			return runtime.None, err
		}

		for _, value := range batch {
			// Reject the incomplete final alphabet interval to avoid modulo bias.
			if int(value) >= 256-256%len(tokenAlphabet) {
				continue
			}

			out[written] = tokenAlphabet[int(value)%len(tokenAlphabet)]
			written++
		}
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	return runtime.String(out), nil
}
