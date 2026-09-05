package crypto_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/crypto"
)

func TestDigestVectors(t *testing.T) {
	for _, tc := range []struct {
		name string
		fn   func(context.Context, runtime.Value) (runtime.Value, error)
		want string
	}{
		{name: "md5", fn: crypto.MD5, want: "3858f62230ac3c915f300c664312c63f"},
		{name: "sha1", fn: crypto.SHA1, want: "8843d7f92416211de9ebb963ff4ce28125932878"},
		{name: "sha256", fn: crypto.SHA256, want: "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2"},
		{name: "sha512", fn: crypto.SHA512, want: "0a50261ebd1a390fed2bf326f2673c145582a6342d523204973d0219337f81616a8069b012587cf5635f6925f1b56c360230c19b273500ee013e030601bf2425"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, value := range []runtime.Value{runtime.String("foobar"), runtime.Binary("foobar")} {
				got, err := tc.fn(context.Background(), value)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(runtime.String(tc.want), got) {
					t.Fatalf("got %#v, want %#v", got, runtime.String(tc.want))
				}
			}

			for _, bad := range []runtime.Value{runtime.Int(1), runtime.True, runtime.None, runtime.NewArray(0)} {
				_, err := tc.fn(context.Background(), bad)
				if !errors.Is(err, runtime.ErrInvalidArgument) {
					t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
				}
			}
		})
	}
}
