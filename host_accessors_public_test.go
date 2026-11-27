package ferret_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	encodingmsgpack "github.com/MontFerret/ferret/v2/pkg/encoding/msgpack"
	"github.com/MontFerret/ferret/v2/pkg/module"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/sdk"
)

var _ interface {
	Encoding() encoding.CodecResolver
	Params() runtime.ParamsReader
	Functions() runtime.FunctionRegistry
} = (*ferret.Engine)(nil)

var hostAccessorViewSink [3]any

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
	if params := eng.Params(); params == nil || params.Size() != 0 || len(params.Names()) != 0 {
		t.Errorf("expected an empty parameter view, got %#v", params)
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

	first, second := logger.With().Logger(), logger.With().Logger()
	first.UpdateContext(func(ctx zerolog.Context) zerolog.Context { return ctx.Str("request", "first") })
	second.UpdateContext(func(ctx zerolog.Context) zerolog.Context { return ctx.Str("request", "other") })
	first.Warn().Msg("visible")

	var event map[string]any
	if err := json.Unmarshal(logs.Bytes(), &event); err != nil {
		t.Fatalf("expected one warning in the configured writer, got %q: %v", logs.String(), err)
	}
	if event["level"] != "warn" || event["message"] != "visible" || event["service"] != "host-accessors" || event["request"] != "first" {
		t.Fatalf("unexpected log event: %v", event)
	}
}

func TestPublicEngineHostAccessorsEncodingView(t *testing.T) {
	const contentType = "application/x-host-accessor-json"
	const bootstrapContentType = "application/x-host-accessor-bootstrap"
	codec := &hostAccessorCodec{Codec: encodingjson.Default, contentType: contentType}
	var registrar encoding.CodecRegistrar
	mod := sdk.NewModule("host-accessor-codec", func(boot module.Bootstrap) error {
		registrar = boot.Host().Encoding()
		return registrar.Register(codec)
	})
	eng := newHostAccessorTestEngine(t, ferret.WithModules(mod))

	resolver := eng.Encoding()
	if got, err := resolver.Codec(contentType); err != nil || got != codec {
		t.Fatalf("expected the module codec instance, got %T, %v", got, err)
	}
	replacement := &hostAccessorCodec{Codec: encodingmsgpack.Default, contentType: contentType}
	if err := registrar.Register(replacement); err != nil {
		t.Fatal(err)
	}
	if err := registrar.Register(&hostAccessorCodec{Codec: encodingjson.Default, contentType: bootstrapContentType}); err != nil {
		t.Fatal(err)
	}

	if got, err := resolver.Codec(contentType); err != nil || got != codec {
		t.Fatalf("expected the original codec after editing bootstrap registrations, got %T, %v", got, err)
	}
	if got, err := resolver.Encoder(contentType); err != nil || got != codec {
		t.Fatalf("expected the original encoder, got %T, %v", got, err)
	}
	if got, err := resolver.Decoder(contentType); err != nil || got != codec {
		t.Fatalf("expected the original decoder, got %T, %v", got, err)
	}
	if _, err := resolver.Codec(bootstrapContentType); !errors.Is(err, encoding.ErrCodecNotFound) {
		t.Fatalf("expected bootstrap-only registration to be absent, got %v", err)
	}

	ctx := t.Context()
	src := ferret.NewAnonymousSource("RETURN 42")
	out, err := collectEngine(eng, ctx, src, ferret.WithOutputContentType(contentType))
	if err != nil || out == nil || out.Metadata.ContentType != contentType || string(out.Data) != "42" {
		t.Fatalf("expected query output from the finalized codec, got %#v, %v", out, err)
	}
	if _, err := collectEngine(eng, ctx, src, ferret.WithOutputContentType(bootstrapContentType)); !errors.Is(err, encoding.ErrCodecNotFound) {
		t.Fatalf("expected execution to reject the bootstrap-only content type, got %v", err)
	}
}

func TestPublicEngineExecutionCodecRegistrationsAreReadOnly(t *testing.T) {
	const contentType = "application/x-host-accessor-json"
	const addedType = "application/x-host-accessor-added"
	original := &hostAccessorCodec{Codec: encodingjson.Default, contentType: contentType}
	replacement := &hostAccessorCodec{Codec: encodingmsgpack.Default, contentType: contentType}
	added := &hostAccessorCodec{Codec: encodingjson.Default, contentType: addedType}
	register := ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
		ns.Function().A0().Add("CHECK_CODEC_REGISTRATIONS", func(ctx context.Context) (runtime.Value, error) {
			snapshot, err := encoding.RegistryFrom(ctx)
			if err != nil {
				return runtime.None, err
			}
			if err := snapshot.Register(replacement); err != nil {
				return runtime.None, err
			}
			if err := snapshot.Register(added); err != nil {
				return runtime.None, err
			}
			if _, err := encoding.CodecFrom(ctx, addedType); !errors.Is(err, encoding.ErrCodecNotFound) {
				return runtime.None, errors.New("context codec lookup observed registration in the snapshot")
			}
			if got, err := encoding.EncoderFrom(ctx, contentType); err != nil || got != original {
				return runtime.None, errors.New("context encoder lookup lost the original codec")
			}
			if got, err := encoding.DecoderFrom(ctx, contentType); err != nil || got != original {
				return runtime.None, errors.New("context decoder lookup lost the original codec")
			}

			return runtime.Int(42), nil
		})
	})
	src := ferret.NewAnonymousSource("RETURN CHECK_CODEC_REGISTRATIONS()")
	outputType := ferret.WithOutputContentType(contentType)

	for _, mode := range []string{"ordinary", "debugger"} {
		t.Run(mode, func(t *testing.T) {
			eng := newHostAccessorTestEngine(t, ferret.WithEncodingRegistry(encoding.NewRegistry(original)), register)
			view := eng.Encoding()
			ctx := t.Context()
			var output *encoding.Content

			if mode == "ordinary" {
				var err error
				output, err = collectEngine(eng, ctx, src, outputType)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				plan, err := eng.CompileDebug(ctx, src)
				if err != nil {
					t.Fatal(err)
				}
				defer plan.Close()

				session, err := plan.NewDebugSession(ctx, outputType)
				if err != nil {
					t.Fatal(err)
				}
				defer session.Close()

				if _, err := session.Start(ctx); err != nil {
					t.Fatal(err)
				}
				event, err := session.Continue(ctx)
				if err != nil || event == nil {
					t.Fatalf("expected a debugger event, got %#v, %v", event, err)
				}
				if event.Reason != debugger.ReasonCompleted {
					t.Fatalf("debugger stopped with %s: %v", event.Reason, event.Error)
				}
				output = event.Output
			}

			if output == nil || output.Metadata.ContentType != contentType || string(output.Data) != "42" {
				t.Fatalf("expected output from the finalized codec, got %#v", output)
			}
			if got, err := view.Codec(contentType); err != nil || got != original {
				t.Fatalf("held view lost the original codec, got %T, %v", got, err)
			}
			if _, err := view.Codec(addedType); !errors.Is(err, encoding.ErrCodecNotFound) {
				t.Fatalf("held view observed a registration from the snapshot: %v", err)
			}
			output, err := collectEngine(eng, ctx, ferret.NewAnonymousSource("RETURN 7"), outputType)
			if err != nil || output == nil || output.Metadata.ContentType != contentType || string(output.Data) != "7" {
				t.Fatalf("expected later output from the finalized codec, got %#v, %v", output, err)
			}
		})
	}
}

func TestPublicEngineHostAccessorsParamsView(t *testing.T) {
	object := runtime.NewObject()
	var bootstrapParams runtime.Params
	mod := sdk.NewModule("host-accessor-params", func(boot module.Bootstrap) error {
		bootstrapParams = boot.Host().Params()
		return bootstrapParams.Set("registered", 22)
	})
	eng := newHostAccessorTestEngine(t,
		ferret.WithRuntimeParam("mutable", object),
		ferret.WithParam("configured", 20),
		ferret.WithModules(mod),
	)

	params := eng.Params()
	if value, found := params.Get("configured"); !found || value != runtime.Int(20) {
		t.Fatalf("configured parameter = %v, %t", value, found)
	}
	if got := params.MustGet("registered"); got != runtime.Int(22) {
		t.Fatalf("module parameter = %v, want 22", got)
	}
	if params.GetOr("mutable", runtime.None) != object || params.GetOr("missing", object) != object {
		t.Fatal("expected existing and fallback values to retain their identity")
	}
	if !params.Has("registered") || params.Has("missing") {
		t.Fatal("expected Has to report finalized parameter membership")
	}

	wantNames := []string{"configured", "mutable", "registered"}
	if got := params.Size(); got != len(wantNames) {
		t.Fatalf("parameter count = %d, want %d", got, len(wantNames))
	}
	names := params.Names()
	if !slices.Equal(names, wantNames) {
		t.Fatalf("parameter names = %v, want %v", names, wantNames)
	}
	names[0] = "changed"
	bootstrapParams["configured"] = runtime.Int(1)
	bootstrapParams.Delete("registered")
	bootstrapParams["added"] = runtime.Int(2)

	for range 4 {
		if got := params.Names(); !slices.Equal(got, wantNames) {
			t.Fatalf("parameter names = %v, want %v", got, wantNames)
		}
	}
	if params.MustGet("configured") != runtime.Int(20) || params.MustGet("registered") != runtime.Int(22) || params.Has("added") {
		t.Fatal("editing bootstrap bindings changed the finalized parameters")
	}
	out, err := collectEngine(eng, t.Context(), ferret.NewAnonymousSource("RETURN @configured + @registered"))
	if err != nil || out == nil || string(out.Data) != "42" {
		t.Fatalf("expected the finalized parameter defaults during execution, got %#v, %v", out, err)
	}
}

func TestPublicEngineHostAccessorsFunctions(t *testing.T) {
	var library runtime.Library
	mod := sdk.NewModule("host-accessor-functions", func(boot module.Bootstrap) error {
		library = boot.Host().Library()
		library.Namespace("Accessors").Function().A0().Add("ModuleValue", func(context.Context) (runtime.Value, error) {
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
	for _, check := range []struct {
		collection string
		found      bool
	}{
		{"A0", functions.A0().Has("Accessors::ModuleValue")},
		{"A1", functions.A1().Has("lower")},
		{"A2", functions.A2().Has("contains")},
		{"A3", functions.A3().Has("substring")},
		{"A4", functions.A4().Has("find_first")},
		{"Var", functions.Var().Has("concat")},
	} {
		if !check.found {
			t.Errorf("expected a registered function in %s", check.collection)
		}
	}

	functions.List()[0] = "changed"
	moduleFunctions := functions.A0().GetAll()
	delete(moduleFunctions, "Accessors::ModuleValue")
	library.Namespace("Accessors").Function().A0().Remove("ModuleValue")
	library.Function().A0().Add("BootstrapOnly", func(context.Context) (runtime.Value, error) {
		return runtime.None, nil
	})

	if !slices.Equal(functions.List(), names) || functions.Size() != size || functions.Hash() != hash {
		t.Fatal("editing returned names or bootstrap registrations changed engine introspection")
	}
	if !functions.A0().Has("Accessors::ModuleValue") || functions.Has("BootstrapOnly") {
		t.Fatal("expected the finalized definitions to survive changes to bootstrap or returned collections")
	}
	out, err := collectEngine(eng, t.Context(), ferret.NewAnonymousSource("RETURN CustomValue() + Accessors::ModuleValue()"))
	if err != nil || out == nil || string(out.Data) != "42" {
		t.Fatalf("expected execution to use the finalized functions, got %#v, %v", out, err)
	}
}

func TestPublicEngineHostAccessorsHideMutableTypes(t *testing.T) {
	eng := newHostAccessorTestEngine(t)
	codecs, params, functions := eng.Encoding(), eng.Params(), eng.Functions()

	if _, ok := codecs.(*encoding.Registry); ok {
		t.Error("codec view exposes the mutable registry")
	}
	if _, ok := codecs.(encoding.CodecRegistrar); ok {
		t.Error("codec view exposes registration")
	}
	if _, ok := params.(runtime.Params); ok {
		t.Error("parameter view exposes the mutable map")
	}
	if _, ok := params.(interface{ Set(string, any) error }); ok {
		t.Error("parameter view exposes binding changes")
	}
	if _, ok := functions.(*runtime.Functions); ok {
		t.Error("function view exposes the replaceable registry")
	}
	if _, ok := functions.(runtime.Library); ok {
		t.Error("function view exposes the registration library")
	}
	if _, ok := any(functions).(runtime.FunctionDefs); ok {
		t.Error("function view exposes function definitions for registration")
	}
	for _, check := range []struct {
		view any
		api  reflect.Type
	}{
		{codecs, reflect.TypeFor[encoding.CodecResolver]()},
		{params, reflect.TypeFor[runtime.ParamsReader]()},
		{functions, reflect.TypeFor[runtime.FunctionRegistry]()},
	} {
		if got := reflect.TypeOf(check.view).NumMethod(); got != check.api.NumMethod() {
			t.Errorf("%v view exposes %d methods, want %d", check.api, got, check.api.NumMethod())
		}
	}
}

func TestPublicEngineHostAccessorsReuseViewsWithoutAllocating(t *testing.T) {
	eng := newHostAccessorTestEngine(t, ferret.WithParam("configured", 20))
	views := [3]any{eng.Encoding(), eng.Params(), eng.Functions()}
	defer func() { hostAccessorViewSink = [3]any{} }()

	reused := true
	allocations := testing.AllocsPerRun(100, func() {
		hostAccessorViewSink = [3]any{eng.Encoding(), eng.Params(), eng.Functions()}
		if hostAccessorViewSink != views {
			reused = false
		}
	})
	if !reused {
		t.Error("expected each accessor to reuse its initialized view")
	}
	if allocations != 0 {
		t.Errorf("allocations per set of accessor calls = %g, want zero", allocations)
	}
}

func TestPublicEngineHostAccessorsConcurrentReads(t *testing.T) {
	eng := newHostAccessorTestEngine(t, ferret.WithParam("configured", 20))

	for worker := range 8 {
		t.Run(strconv.Itoa(worker), func(t *testing.T) {
			t.Parallel()

			for range 20 {
				if _, err := eng.Encoding().Codec(encodingjson.ContentType); err != nil {
					t.Fatal(err)
				}
				params := eng.Params()
				if params.MustGet("configured") != runtime.Int(20) || !slices.Equal(params.Names(), []string{"configured"}) {
					t.Fatal("unexpected parameter view")
				}
				functions := eng.Functions()
				if !functions.Var().Has("concat") || !slices.Contains(functions.List(), "concat") {
					t.Fatal("unexpected function view")
				}
			}
		})
	}
}
