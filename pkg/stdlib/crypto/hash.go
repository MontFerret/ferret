package crypto

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// md5 computes the MD5 digest of raw bytes. Strings supply their UTF-8 bytes. This legacy digest is unsuitable for security-sensitive collision resistance.
// @param value {String | Binary} The input data.
// @return {String} The lowercase hexadecimal digest.
func MD5(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg, 0, runtime.TypeString, runtime.TypeBinary); err != nil {
		return runtime.None, err
	}

	digest := md5.Sum(runtime.ToBinary(arg))

	return runtime.String(hex.EncodeToString(digest[:])), nil
}

// sha1 computes the SHA1 digest of raw bytes. Strings supply their UTF-8 bytes. This legacy digest is unsuitable for security-sensitive collision resistance.
// @param value {String | Binary} The input data.
// @return {String} The lowercase hexadecimal digest.
func SHA1(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg, 0, runtime.TypeString, runtime.TypeBinary); err != nil {
		return runtime.None, err
	}

	digest := sha1.Sum(runtime.ToBinary(arg))

	return runtime.String(hex.EncodeToString(digest[:])), nil
}

// sha256 computes the SHA256 digest of raw bytes. Strings supply their UTF-8 bytes.
// @param value {String | Binary} The input data.
// @return {String} The lowercase hexadecimal digest.
func SHA256(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg, 0, runtime.TypeString, runtime.TypeBinary); err != nil {
		return runtime.None, err
	}

	digest := sha256.Sum256(runtime.ToBinary(arg))

	return runtime.String(hex.EncodeToString(digest[:])), nil
}

// sha512 computes the SHA512 digest of raw bytes. Strings supply their UTF-8 bytes.
// @param value {String | Binary} The input data.
// @return {String} The lowercase hexadecimal digest.
func SHA512(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg, 0, runtime.TypeString, runtime.TypeBinary); err != nil {
		return runtime.None, err
	}

	digest := sha512.Sum512(runtime.ToBinary(arg))

	return runtime.String(hex.EncodeToString(digest[:])), nil
}
