package ferret_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/module"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

type hostAccessorCodec struct {
	encoding.Codec
	contentType string
}

func (c *hostAccessorCodec) ContentType() string {
	return c.contentType
}

func newHostAccessorTestEngine(t *testing.T, opts ...ferret.Option) *ferret.Engine {
	t.Helper()

	eng, err := ferret.New(opts...)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := eng.Close(); err != nil {
			t.Errorf("close engine: %v", err)
		}
	})

	return eng
}

func TestPublicEngineHostAccessorsDefaults(t *testing.T) {
	eng := newHostAccessorTestEngine(t)

	if got := ferret.LogLevel(eng.Logger().GetLevel()); got != ferret.LogError {
		t.Errorf("default log level = %v, want %v", got, ferret.LogError)
	}
	codecs := eng.Encoding()
	for _, contentType := range []string{encodingjson.ContentType, encodingmsgpack.ContentType} {
		if _, err := codecs.Codec(contentType); err != nil {
			t.Errorf("default codec for %s: %v", contentType, err)
		}
	}
	if params := eng.Params(); params == nil || len(params) != 0 {
		t.Errorf("expected an allocated empty parameter snapshot, got %#v", params)
	}
	if functions := eng.Functions(); functions == nil || !functions.Has("CONCAT") {
		t.Error("expected the default function registry to contain CONCAT")
	}
}

func TestPublicEngineHostAccessorsLogger(t *testing.T) {
	var logs bytes.Buffer
	eng := newHostAccessorTestEngine(t,
		ferret.WithLog(&logs),
		ferret.WithLogLevel(ferret.LogWarn),
		ferret.WithLogFields(map[string]any{"service": "host-accessors"}),
	)
	logger := eng.Logger()
	if got := ferret.LogLevel(logger.GetLevel()); got != ferret.LogWarn {
		t.Fatalf("configured log level = %v, want %v", got, ferret.LogWarn)
	}
	logger.Info().Msg("suppressed")
	logger.Warn().Msg("visible")

	var event map[string]any
	if err := json.Unmarshal(logs.Bytes(), &event); err != nil {
		t.Fatalf("expected one warning in the configured writer, got %q: %v", logs.String(), err)
	}
	if event["level"] != "warn" || event["message"] != "visible" || event["service"] != "host-accessors" {
		t.Fatalf("unexpected log event: %v", event)
	}
}

func TestPublicEngineHostAccessorsEncodingSnapshot(t *testing.T) {
	const contentType = "application/x-host-accessor-json"
	const snapshotContentType = "application/x-host-accessor-snapshot"
	codec := &hostAccessorCodec{Codec: encodingjson.Default, contentType: contentType}
	mod := sdk.NewModule("host-accessor-codec", func(boot module.Bootstrap) error {
		return boot.Host().Encoding().Register(codec)
	})
	eng := newHostAccessorTestEngine(t, ferret.WithModules(mod))

	snapshot := eng.Encoding()
	if got, err := snapshot.Codec(contentType); err != nil || got != codec {
		t.Fatalf("snapshot should share the module codec, got %T, %v", got, err)
	}
	replacement := &hostAccessorCodec{Codec: encodingmsgpack.Default, contentType: contentType}
	if err := snapshot.Register(replacement); err != nil {
		t.Fatal(err)
	}
	if err := snapshot.Register(&hostAccessorCodec{Codec: encodingjson.Default, contentType: snapshotContentType}); err != nil {
		t.Fatal(err)
	}
	if got, err := snapshot.Codec(contentType); err != nil || got != replacement {
		t.Fatalf("expected the replacement in the edited snapshot, got %T, %v", got, err)
	}

	fresh := eng.Encoding()
	if got, err := fresh.Codec(contentType); err != nil || got != codec {
		t.Fatalf("expected the engine codec in a fresh snapshot, got %T, %v", got, err)
	}
	if _, err := fresh.Codec(snapshotContentType); !errors.Is(err, encoding.ErrCodecNotFound) {
		t.Fatalf("expected the added registration to stay in the edited snapshot, got %v", err)
	}

	ctx := t.Context()
	src := ferret.NewAnonymousSource("RETURN 42")
	out, err := collectEngine(eng, ctx, src, ferret.WithOutputContentType(contentType))
	if err != nil || out == nil || out.Metadata.ContentType != contentType || string(out.Data) != "42" {
		t.Fatalf("expected query output from the original codec, got %#v, %v", out, err)
	}
	if _, err := collectEngine(eng, ctx, src, ferret.WithOutputContentType(snapshotContentType)); !errors.Is(err, encoding.ErrCodecNotFound) {
		t.Fatalf("expected execution to reject the snapshot-only content type, got %v", err)
	}
}

func TestPublicEngineHostAccessorsParamsSnapshot(t *testing.T) {
	object := runtime.NewObject()
	mod := sdk.NewModule("host-accessor-params", func(boot module.Bootstrap) error {
		return boot.Host().Params().Set("registered", 22)
	})
	eng := newHostAccessorTestEngine(t,
		ferret.WithParam("configured", 20),
		ferret.WithRuntimeParam("mutable", object),
		ferret.WithModules(mod),
	)

	snapshot := eng.Params()
	if snapshot["configured"] != runtime.Int(20) || snapshot["registered"] != runtime.Int(22) || snapshot["mutable"] != object {
		t.Fatalf("unexpected finalized parameter bindings: %v", snapshot)
	}
	snapshot["configured"] = runtime.Int(1)
	delete(snapshot, "registered")
	snapshot["added"] = runtime.Int(2)

	fresh := eng.Params()
	if fresh["configured"] != runtime.Int(20) || fresh["registered"] != runtime.Int(22) || fresh["mutable"] != object || fresh.Has("added") {
		t.Fatalf("snapshot edits changed the engine defaults: %v", fresh)
	}
	out, err := collectEngine(eng, t.Context(), ferret.NewAnonymousSource("RETURN @configured + @registered"))
	if err != nil || out == nil || string(out.Data) != "42" {
		t.Fatalf("expected the original parameter defaults during execution, got %#v, %v", out, err)
	}
}

func TestPublicEngineHostAccessorsFunctions(t *testing.T) {
	mod := sdk.NewModule("host-accessor-functions", func(boot module.Bootstrap) error {
		boot.Host().Library().Namespace("Accessors").Function().A0().Add("ModuleValue", func(context.Context) (runtime.Value, error) {
			return runtime.Int(22), nil
		})
		return nil
	})
	eng := newHostAccessorTestEngine(t, ferret.WithModules(mod), ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("CustomValue", func(context.Context) (runtime.Value, error) {
			return runtime.Int(20), nil
		})
	}))

	functions := eng.Functions()
	names := slices.Clone(functions.List())
	size, hash := functions.Size(), functions.Hash()
	for _, name := range []string{"concat", "CustomValue", "Accessors::ModuleValue"} {
		if !slices.Contains(names, name) || !functions.Has(strings.ToUpper(name)) {
			t.Fatalf("expected to list and resolve %q", name)
		}
	}
	if size < len(names) || hash == 0 {
		t.Fatalf("unexpected registry metadata: size %d for %d names, hash %d", size, len(names), hash)
	}

	functions.List()[0] = "changed"
	*functions = runtime.Functions{}

	fresh := eng.Functions()
	if !slices.Equal(fresh.List(), names) || fresh.Size() != size || fresh.Hash() != hash {
		t.Fatal("editing the returned registry or name list changed engine introspection")
	}
	out, err := collectEngine(eng, t.Context(), ferret.NewAnonymousSource("RETURN CustomValue() + Accessors::ModuleValue()"))
	if err != nil || out == nil || string(out.Data) != "42" {
		t.Fatalf("expected execution to use the finalized functions, got %#v, %v", out, err)
	}
}
