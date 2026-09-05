package crypto

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RegisterLib registers digest and secure token operations in the crypto namespace.
func RegisterLib(ns runtime.Namespace) {
	ns.Namespace("crypto").Function().A1().Add("md5", MD5).Add("sha1", SHA1).
		Add("sha256", SHA256).Add("sha512", SHA512).Add("random_token", RandomToken)
}
