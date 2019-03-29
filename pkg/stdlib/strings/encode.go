package strings

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"net/url"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// EncodeURIComponent returns the encoded String of uri.
// @param (String) - Uri to encode.
// @returns String - Encoded string.
func EncodeURIComponent(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	str := url.QueryEscape(args[0].String())

	return values.NewString(str), nil
}

// Md5 calculates the MD5 checksum for text and return it in a hexadecimal string representation.
// @param text (String) - The string to do calculations against to.
// @return (String) - MD5 checksum as hex string.
func Md5(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	res := md5.Sum([]byte(text))

	return values.NewString(string(res[:])), nil
}

// Sha1 calculates the SHA1 checksum for text and returns it in a hexadecimal string representation.
// @param text (String) - The string to do calculations against to.
// @return (String) - Sha1 checksum as hex string.
func Sha1(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	res := sha1.Sum([]byte(text))

	return values.NewString(string(res[:])), nil
}

// Sha512 calculates the SHA512 checksum for text and returns it in a hexadecimal string representation.
// @param text (String) - The string to do calculations against to.
// @return (String) - SHA512 checksum as hex string.
func Sha512(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	res := sha512.Sum512([]byte(text))

	return values.NewString(string(res[:])), nil
}

// ToBase64 returns the base64 representation of value.
// @param value (string) - The string to encode.
// @returns toBase64String (String) - A base64 representation of the string.
func ToBase64(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	value := args[0].String()
	out := base64.StdEncoding.EncodeToString([]byte(value))

	return values.NewString(out), nil
}
