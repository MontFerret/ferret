package strings

import (
	"context"
	"math/rand"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
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
func RandomToken(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	err = core.ValidateType(args[0], types.Int)

	if err != nil {
		return values.EmptyString, err
	}

	size := args[0].(values.Int)
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

	return values.NewString(string(b)), nil
}
