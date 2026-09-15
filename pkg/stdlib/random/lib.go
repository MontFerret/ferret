package random

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RegisterLib registers non-cryptographic random value generation. These
// functions are unsuitable for passwords, authentication tokens, or secrets.
// @namespace random
func RegisterLib(ns runtime.Namespace) {
	ns = ns.Namespace("random")
	ns.Function().A0().Add("float", float0).Add("bool", Bool)
	ns.Function().A2().Add("float", float2).Add("int", Int)
}

func source(ctx context.Context) (*rnd.Source, error) {
	src, ok := rnd.FromContext(ctx)
	if !ok {
		return nil, runtime.Error(runtime.ErrUnexpected, "random source missing from execution context")
	}

	return src, nil
}
