package strings

import (
	"context"
	"math/rand"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// randSrc is a global variable because of this issue
// https://github.com/golang/go/issues/8926
var randSrc = rand.NewSource(time.Now().UnixNano())

// RANDOM_TOKEN generates a pseudo-random token string with the specified length. The algorithm for token generation should be treated as opaque.
// @param {Int} len - The desired string length for the token. It must be greater than 0 and at most 65536.
// @return {String} - A generated token consisting of lowercase letters, uppercase letters and numbers.
func RandomToken(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.EmptyString, err
	}

	size, err := runtime.CastInt(args[0])

	if err != nil {
		return runtime.EmptyString, err
	}

	b := make([]byte, size)

	for i, cache, remain := size-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return runtime.NewString(string(b)), nil
}
