package crypto

import (
	"bytes"
	"context"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestRandomTokenSampling(t *testing.T) {
	// Eight rejected bytes followed by every accepted byte: each alphabet
	// character must occur exactly four times, without depending on chance.
	input := []byte{248, 249, 250, 251, 252, 253, 254, 255}
	for value := 0; value < 248; value++ {
		input = append(input, byte(value))
	}

	got, err := randomToken(context.Background(), bytes.NewReader(input), 248)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(runtime.String(strings.Repeat(tokenAlphabet, 4)), got) {
		t.Fatalf("got %#v, want %#v", got, runtime.String(strings.Repeat(tokenAlphabet, 4)))
	}

	_, err = randomToken(context.Background(), bytes.NewReader([]byte{0}), 2)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("error = %v, want %v", err, io.ErrUnexpectedEOF)
	}

	failure := errors.New("entropy failed")
	_, err = randomToken(context.Background(), iotestReader{err: failure}, 1)
	if !errors.Is(err, failure) {
		t.Fatalf("error = %v, want %v", err, failure)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = randomToken(ctx, cancelingTokenReader{cancel: cancel}, 4)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want cancellation during rejection sampling", err)
	}
}

func TestRandomTokenValidationAndConcurrency(t *testing.T) {
	for _, bad := range []runtime.Value{runtime.Int(-1), runtime.Int(0), runtime.Int(65537), runtime.Int(1<<63 - 1), runtime.True, runtime.String("32"), runtime.Float(32), runtime.None} {
		_, err := RandomToken(context.Background(), bad)
		if !errors.Is(err, runtime.ErrInvalidArgument) {
			t.Fatalf("error = %v, want %v", err, runtime.ErrInvalidArgument)
		}
	}

	for _, size := range []runtime.Int{1, 32, 65536} {
		got, err := RandomToken(context.Background(), size)
		if err != nil {
			t.Fatal(err)
		}

		if len(got.String()) != int(size) {
			t.Fatalf("length = %d, want %d", len(got.String()), int(size))
		}

		if strings.Trim(got.String(), tokenAlphabet) != "" {
			t.Fatalf("unexpected characters: %q", strings.Trim(got.String(), tokenAlphabet))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := RandomToken(ctx, runtime.Int(32))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want %v", err, context.Canceled)
	}

	for i := 0; i < 32; i++ {
		t.Run("concurrent", func(t *testing.T) {
			t.Parallel()

			for j := 0; j < 32; j++ {
				got, err := RandomToken(context.Background(), runtime.Int(128))
				if err != nil {
					t.Fatal(err)
				}

				if len(got.String()) != 128 {
					t.Fatalf("length = %d, want %d", len(got.String()), 128)
				}

				if strings.Trim(got.String(), tokenAlphabet) != "" {
					t.Fatalf("unexpected characters: %q", strings.Trim(got.String(), tokenAlphabet))
				}
			}
		})
	}
}
