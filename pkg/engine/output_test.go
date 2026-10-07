package engine

import (
	"context"
	stdjson "encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type aliasCodec struct {
	base        ferretencoding.Codec
	pre         ferretencoding.PreEncoderHook
	contentType string
}

func (c aliasCodec) ContentType() string {
	return c.contentType
}

func (c aliasCodec) Encode(ctx context.Context, dst io.Writer, value runtime.Value) error {
	return c.base.Encode(ctx, dst, value)
}

func (c aliasCodec) EncodeWith() ferretencoding.EncoderConfigurer {
	return c.base.EncodeWith().PreHook(c.pre)
}

func (c aliasCodec) Decode(ctx context.Context, src io.Reader) (runtime.Value, error) {
	return c.base.Decode(ctx, src)
}

func (c aliasCodec) DecodeWith() ferretencoding.DecoderConfigurer {
	return c.base.DecodeWith()
}

func TestSessionRunReturnsDefaultJSONOutput(t *testing.T) {
	eng := mustNewEngine(t)
	plan := mustCompilePlan(t, eng, "RETURN 1")
	session := mustNewSession(t, plan)

	out, err := collectSession(session, context.Background())
	if err != nil {
		t.Fatalf("expected session run to succeed, got %v", err)
	}

	if out.Metadata.ContentType != encodingjson.ContentType {
		t.Fatalf("unexpected content type: got %q, want %q", out.Metadata.ContentType, encodingjson.ContentType)
	}

	if got := strings.TrimSpace(string(out.Data)); got != "1" {
		t.Fatalf("unexpected output payload: got %q", got)
	}
}

func TestSessionRunUsesRequestedOutputCodec(t *testing.T) {
	const customContentType = "application/x-test-json"

	eng := mustNewEngine(t, WithEncodingCodec(customContentType, aliasCodec{
		contentType: customContentType,
		base:        encodingjson.Default,
	}))
	plan := mustCompilePlan(t, eng, "RETURN 1")
	session := mustNewSession(t, plan, WithOutputContentType(customContentType))

	out, err := collectSession(session, context.Background())
	if err != nil {
		t.Fatalf("expected session run to succeed, got %v", err)
	}

	if out.Metadata.ContentType != customContentType {
		t.Fatalf("unexpected content type: got %q, want %q", out.Metadata.ContentType, customContentType)
	}

	if got := strings.TrimSpace(string(out.Data)); got != "1" {
		t.Fatalf("unexpected output payload: got %q", got)
	}
}

func TestSessionRunClosesResultWhenRequestedCodecIsMissing(t *testing.T) {
	root := newTrackingJSONCloser("root-missing-codec", "1")

	eng := mustNewEngine(t, WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
			return root, nil
		})
	}))
	plan := mustCompilePlan(t, eng, "RETURN MAKE()")
	session := mustNewSession(t, plan, WithOutputContentType("application/x-missing"))

	out, err := collectSession(session, context.Background())
	if out != nil {
		t.Fatal("expected output to be nil when codec resolution fails")
	}

	if err == nil {
		t.Fatal("expected session run to fail when codec is missing")
	}

	if !errors.Is(err, ferretencoding.ErrCodecNotFound) {
		t.Fatalf("expected codec-not-found error, got %v", err)
	}

	if got := root.closed; got != 1 {
		t.Fatalf("expected missing-codec path to close encountered result resources once, got %d closes", got)
	}
}

func TestSessionRunReturnsOutputWhenResultCleanupFails(t *testing.T) {
	closeErr := errors.New("close boom")
	closer := newTrackingJSONCloser("cleanup-failure", "1")
	closer.closeErr = closeErr

	eng := mustNewEngine(t, WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
			return closer, nil
		})
	}))
	plan := mustCompilePlan(t, eng, "RETURN MAKE()")
	session := mustNewSession(t, plan)

	out, err := collectSession(session, context.Background())
	if out == nil {
		t.Fatal("expected output to be returned when cleanup fails after materialization")
	}

	if !errors.Is(err, closeErr) {
		t.Fatalf("expected cleanup failure to be returned, got %v", err)
	}

	if got := strings.TrimSpace(string(out.Data)); got != "1" {
		t.Fatalf("unexpected output payload: got %q", got)
	}
}

func TestSessionRunClosesNestedLiveValuesDiscoveredDuringEncoding(t *testing.T) {
	nested := newTrackingJSONCloser("nested-live", `"ok"`)

	eng := mustNewEngine(t, WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("MAKE", func(context.Context) (runtime.Value, error) {
			return runtime.NewArrayWith(nested), nil
		})
	}))
	plan := mustCompilePlan(t, eng, "RETURN MAKE()")
	session := mustNewSession(t, plan)

	out, err := collectSession(session, context.Background())
	if err != nil {
		t.Fatalf("expected session run to succeed, got %v", err)
	}

	var decoded []string
	if err := stdjson.Unmarshal(out.Data, &decoded); err != nil {
		t.Fatalf("failed to decode output: %v", err)
	}

	if len(decoded) != 1 || decoded[0] != "ok" {
		t.Fatalf("unexpected output payload: %v", decoded)
	}

	if got := nested.closed; got != 1 {
		t.Fatalf("expected nested live value to be closed after materialization, got %d closes", got)
	}
}
