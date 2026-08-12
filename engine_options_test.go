package ferret

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func mustNewOptionsForTest(t *testing.T, setters ...Option) *options {
	t.Helper()

	opts, err := newOptions(setters)
	if err != nil {
		t.Fatalf("failed to create options: %v", err)
	}

	return opts
}

func TestNewOptionsIgnoresNilTopLevelOption(t *testing.T) {
	t.Parallel()

	opts := mustNewOptionsForTest(t, nil, WithParam("param1", "value1"))

	value, ok := opts.params.Get("param1")
	if !ok {
		t.Fatal("expected param from later option to be applied")
	}

	if value != runtime.NewString("value1") {
		t.Fatalf("expected param1 to remain value1, got: %v", value)
	}
}

func TestNewOptionsIgnoresEmptyParamsOptions(t *testing.T) {
	t.Parallel()

	opts := mustNewOptionsForTest(
		t,
		WithParam("param1", "value1"),
		WithParams(nil),
		WithParams(map[string]any{}),
	)

	if len(opts.params) != 1 {
		t.Fatalf("expected params to remain unchanged, got %d entries", len(opts.params))
	}

	value, ok := opts.params.Get("param1")
	if !ok {
		t.Fatal("expected param1 to remain configured")
	}

	if value != runtime.NewString("value1") {
		t.Fatalf("expected param1 to remain value1, got: %v", value)
	}
}

func TestNewOptionsIgnoresEmptyRuntimeParamsOptions(t *testing.T) {
	t.Parallel()

	opts := mustNewOptionsForTest(
		t,
		WithRuntimeParam("param1", runtime.NewString("value1")),
		WithRuntimeParams(nil),
		WithRuntimeParams(runtime.Params{}),
	)

	if len(opts.params) != 1 {
		t.Fatalf("expected params to remain unchanged, got %d entries", len(opts.params))
	}

	value, ok := opts.params.Get("param1")
	if !ok {
		t.Fatal("expected param1 to remain configured")
	}

	if value != runtime.NewString("value1") {
		t.Fatalf("expected param1 to remain value1, got: %v", value)
	}
}

func TestNewOptionsIgnoresEmptyLogFields(t *testing.T) {
	t.Parallel()

	opts := mustNewOptionsForTest(
		t,
		WithLogFields(map[string]any{"component": "engine"}),
		WithLogFields(nil),
		WithLogFields(map[string]any{}),
	)

	if len(opts.logger) != 1 {
		t.Fatalf("expected logger options to remain unchanged, got %d entries", len(opts.logger))
	}
}

func TestNewOptionsAcceptsEmptyModulesOption(t *testing.T) {
	t.Parallel()

	opts := mustNewOptionsForTest(t, WithModules())

	if len(opts.modules) != 0 {
		t.Fatalf("expected no modules to be registered, got %d", len(opts.modules))
	}
}

func TestNewOptionsRejectsNilModule(t *testing.T) {
	t.Parallel()

	_, err := newOptions([]Option{WithModules(nil)})
	if err == nil {
		t.Fatal("expected nil module to fail")
	}

	if !strings.Contains(err.Error(), "module cannot be nil") {
		t.Fatalf("expected nil module validation error, got: %v", err)
	}
}

func TestNewOptionsAcceptsEmptyCompilerOptions(t *testing.T) {
	t.Parallel()

	opts := mustNewOptionsForTest(
		t,
		WithCompilerOptions(compiler.WithOptimizationLevel(compiler.O0)),
		WithCompilerOptions(),
	)

	if len(opts.compiler) != 1 {
		t.Fatalf("expected compiler options to remain unchanged, got %d entries", len(opts.compiler))
	}
}

func TestNewOptionsRejectsNilCompilerOption(t *testing.T) {
	t.Parallel()

	_, err := newOptions([]Option{WithCompilerOptions(nil)})
	if err == nil {
		t.Fatal("expected nil compiler option to fail")
	}

	if !strings.Contains(err.Error(), "compiler option cannot be nil") {
		t.Fatalf("expected nil compiler option validation error, got: %v", err)
	}
}

func TestNewOptionsTrimsFSRoot(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	opts := mustNewOptionsForTest(t, WithFSRoot("  "+root+"\n"))

	if opts.fsRoot != root {
		t.Fatalf("expected fs root to be trimmed to %q, got %q", root, opts.fsRoot)
	}
}

func TestNewOptionsRejectsBlankFSRoot(t *testing.T) {
	t.Parallel()

	_, err := newOptions([]Option{WithFSRoot(" \t\n ")})
	if err == nil {
		t.Fatal("expected blank fs root to fail")
	}

	if !strings.Contains(err.Error(), "fs root cannot be empty") {
		t.Fatalf("expected blank fs root validation error, got: %v", err)
	}
}

func TestWithNetworkOptionsWithoutSettersIsNoOp(t *testing.T) {
	t.Parallel()

	network := mustNewTestNetwork(t)
	opts := mustNewOptionsForTest(t, WithNetwork(network), WithNetworkOptions())

	if opts.network != network {
		t.Fatalf("expected injected network to remain configured, got %T", opts.network)
	}

	if !opts.hostNetwork {
		t.Fatal("expected injected network to remain caller-owned")
	}
}

func TestWithNetworkOptionsReturnsNetworkConstructionError(t *testing.T) {
	t.Parallel()

	_, err := newOptions([]Option{WithNetworkOptions(
		ferretnet.WithHTTPPolicies(ferrethttp.WithMaxResponseSize(-1)),
	)})
	if err == nil {
		t.Fatal("expected invalid network options to fail")
	}

	if !errors.Is(err, ferrethttp.ErrInvalidPolicyConfiguration) {
		t.Fatalf("expected invalid policy configuration error, got: %v", err)
	}

	if !strings.Contains(err.Error(), "create network: http client:") {
		t.Fatalf("expected network construction context, got: %v", err)
	}
}

func TestWithStdlibSafeRegistersSelectedGroups(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithStdlib(stdlib.Safe()))
	defer func() { _ = eng.Close() }()

	if !eng.host.functions.Has("CONCAT") {
		t.Fatal("expected safe stdlib to register non-IO functions")
	}

	for _, name := range []string{"IO::FS::READ", "IO::NET::HTTP::GET"} {
		if eng.host.functions.Has(name) {
			t.Fatalf("expected safe stdlib to exclude %s", name)
		}
	}
}

func TestWithFunctionsRegistrarCanonicalizesQualifiedHostFunctionNames(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t,
		WithoutStdlib(),
		WithFunctionsRegistrar(func(ns runtime.Namespace) {
			ns.Namespace("Tools").Namespace("Risk").Function().A0().Add("Calculate_Risk", func(context.Context) (runtime.Value, error) {
				return runtime.NewString("ok"), nil
			})
		}),
	)
	defer func() { _ = eng.Close() }()

	if got := eng.host.functions.List(); !slices.Equal(got, []string{"tools::risk::calculate_risk"}) {
		t.Fatalf("expected canonical host metadata, got %v", got)
	}

	for _, query := range []string{
		"return tools::risk::calculate_risk()",
		"return TOOLS::RISK::CALCULATE_RISK()",
		"return Tools::Risk::Calculate_Risk()",
	} {
		output, err := eng.Run(t.Context(), source.NewAnonymous(query))
		if err != nil {
			t.Fatalf("run %q: %v", query, err)
		}

		if got := string(output.Content); got != `"ok"` {
			t.Fatalf("run %q = %s, want %q", query, got, `"ok"`)
		}
	}
}

func TestWithFunctionsRegistrarRejectsCaseOnlyDuplicate(t *testing.T) {
	t.Parallel()

	eng, err := New(
		WithoutStdlib(),
		WithFunctionsRegistrar(func(ns runtime.Namespace) {
			fn := func(context.Context) (runtime.Value, error) { return runtime.None, nil }
			ns.Function().A0().Add("foo", fn).Add("FOO", fn)
		}),
	)
	if eng != nil {
		_ = eng.Close()
	}

	if err == nil {
		t.Fatal("expected case-only duplicate registration to fail")
	}
}

func TestUnknownHostFunctionDiagnosticPreservesSourceSpelling(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithoutStdlib())
	defer func() { _ = eng.Close() }()

	const name = "TOOLS::MiSsInG"
	_, err := eng.Run(t.Context(), source.NewAnonymous("return "+name+"()"))
	if err == nil {
		t.Fatal("expected unknown host function to fail")
	}

	var runtimeErr *vm.RuntimeError
	if !errors.As(err, &runtimeErr) {
		t.Fatalf("expected RuntimeError, got %T: %v", err, err)
	}

	want := "function '" + name + "' is not registered"
	if got := runtimeErr.Spans[0].Label; got != want {
		t.Fatalf("diagnostic label = %q, want %q", got, want)
	}
}

func TestResolvedHostFunctionDiagnosticUsesCanonicalQualifiedName(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t,
		WithoutStdlib(),
		WithFunctionsRegistrar(func(ns runtime.Namespace) {
			ns.Namespace("DB").Namespace("POSTGRES").Function().A1().Add("QUERY", func(context.Context, runtime.Value) (runtime.Value, error) {
				return runtime.None, nil
			})
		}),
	)
	defer func() { _ = eng.Close() }()

	_, err := eng.Run(t.Context(), source.NewAnonymous("return Db::Postgres::Query()"))
	if err == nil {
		t.Fatal("expected invalid host function arity to fail")
	}

	var runtimeErr *vm.RuntimeError
	if !errors.As(err, &runtimeErr) {
		t.Fatalf("expected RuntimeError, got %T: %v", err, err)
	}

	if got, want := runtimeErr.Note, "db::postgres::query expects 1 argument, but got 0"; got != want {
		t.Fatalf("diagnostic note = %q, want %q", got, want)
	}
}

func TestWithStdlibEmptyMatchesWithoutStdlib(t *testing.T) {
	t.Parallel()

	eng := mustNewEngine(t, WithStdlib(stdlib.Empty()))
	defer func() { _ = eng.Close() }()

	if eng.host.functions.Size() != 0 {
		t.Fatalf("expected empty stdlib to register no functions, got %d", eng.host.functions.Size())
	}
}

func TestStdlibOptionsUseLastSelection(t *testing.T) {
	t.Parallel()

	withoutThenFull := mustNewEngine(t, WithoutStdlib(), WithStdlib(stdlib.Full()))
	defer func() { _ = withoutThenFull.Close() }()

	if !withoutThenFull.host.functions.Has("CONCAT") {
		t.Fatal("expected WithStdlib after WithoutStdlib to restore full stdlib")
	}

	fullThenWithout := mustNewEngine(t, WithStdlib(stdlib.Full()), WithoutStdlib())
	defer func() { _ = fullThenWithout.Close() }()

	if fullThenWithout.host.functions.Size() != 0 {
		t.Fatalf("expected WithoutStdlib after WithStdlib to disable stdlib, got %d functions", fullThenWithout.host.functions.Size())
	}
}

func TestWithStdlibRejectsInvalidGroup(t *testing.T) {
	t.Parallel()

	_, err := New(WithStdlib(stdlib.Only(stdlib.Group("unknown"))))
	if err == nil {
		t.Fatal("expected invalid stdlib group to fail engine creation")
	}

	if !strings.Contains(err.Error(), "stdlib: invalid stdlib group(s): unknown") {
		t.Fatalf("expected invalid stdlib group error, got: %v", err)
	}
}
