package uapi

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestPortableParamsMetadata(t *testing.T) {
	for _, tc := range []struct {
		name  string
		query string
		want  []string
	}{
		{"parameters", "RETURN [@second, @first, @second]", []string{"second", "first"}},
		{"no parameters", "RETURN 1", []string{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			portable := newTestRuntime(t)
			plan, err := portable.Compile(t.Context(), api.NewAnonymousSource(tc.query))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = plan.Close() })

			for _, closed := range []bool{false, true} {
				if closed {
					if err := plan.Close(); err != nil {
						t.Fatal(err)
					}
				}

				params, err := plan.Params(t.Context())
				if err != nil || !slices.Equal(params, tc.want) {
					t.Fatalf("Params (closed=%v) = %v, %v, want %v", closed, params, err, tc.want)
				}

				if len(params) > 0 {
					params[0] = "changed"
				}

				again, err := plan.Params(t.Context())
				if err != nil || !slices.Equal(again, tc.want) {
					t.Fatalf("Params after mutation (closed=%v) = %v, %v, want %v", closed, again, err, tc.want)
				}

				checkMetadataContexts(t, func(ctx context.Context) error {
					params, err := plan.Params(ctx)
					if params != nil {
						t.Errorf("rejected context returned parameters: %v", params)
					}

					return err
				})
			}
		})
	}
}

func TestRuntimeVersionMetadata(t *testing.T) {
	for _, constructor := range []struct {
		new  func(*testing.T, api.Version) api.Runtime
		name string
	}{
		{name: "New", new: func(t *testing.T, version api.Version) api.Runtime {
			portable, err := New(version)
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = portable.Close() })

			return portable
		}},
		{name: "Wrap", new: func(t *testing.T, version api.Version) api.Runtime {
			return Wrap(newTestEngine(t), version)
		}},
	} {
		t.Run(constructor.name, func(t *testing.T) {
			for _, tc := range []struct {
				name    string
				version api.Version
			}{
				{name: "opaque version", version: "Core custom build + LOCAL"},
				{name: "empty version"},
			} {
				t.Run(tc.name, func(t *testing.T) {
					portable := constructor.new(t, tc.version)
					for _, closed := range []bool{false, true} {
						if closed {
							if err := portable.Close(); err != nil {
								t.Fatal(err)
							}
						}

						got, err := portable.Version(t.Context())
						if err != nil || got != tc.version {
							t.Fatalf("Version (closed=%v) = %q, %v, want %q", closed, got, err, tc.version)
						}

						checkMetadataContexts(t, func(ctx context.Context) error {
							version, err := portable.Version(ctx)
							if version != "" {
								t.Errorf("rejected context returned version: %q", version)
							}

							return err
						})
					}
				})
			}
		})
	}
}

func checkMetadataContexts(t *testing.T, call func(context.Context) error) {
	t.Helper()

	canceled, cancel := context.WithCancel(t.Context())
	cancel()

	expired, stop := context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
	t.Cleanup(stop)

	for _, tc := range []struct {
		ctx  context.Context
		want error
	}{
		{nil, runtime.ErrInvalidArgument},
		{canceled, context.Canceled},
		{expired, context.DeadlineExceeded},
	} {
		if err := call(tc.ctx); !errors.Is(err, tc.want) {
			t.Fatalf("metadata context error = %v, want %v", err, tc.want)
		}
	}
}
