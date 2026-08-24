package ferret

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	gooptions "github.com/ziflex/go-options"

	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
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

	return &opts
}

func TestNewOptionsPreservesDefaults(t *testing.T) {
	t.Parallel()

	opts := mustNewOptionsForTest(t)

	if opts.library == nil {
		t.Fatal("expected default runtime library")
	}

	if opts.encoding == nil {
		t.Fatal("expected default encoding registry")
	}

	if opts.programLoader == nil {
		t.Fatal("expected default program loader")
	}

	if opts.maxActiveSessions != defaultMaxActiveSessions {
		t.Fatalf("max active sessions = %d, want %d", opts.maxActiveSessions, defaultMaxActiveSessions)
	}

	if opts.maxIdleVMsPerPlan != defaultVMPoolSize {
		t.Fatalf("max idle VMs per plan = %d, want %d", opts.maxIdleVMsPerPlan, defaultVMPoolSize)
	}

	if opts.maxVMsPerPlan != defaultMaxVMsPerPlan {
		t.Fatalf("max VMs per plan = %d, want %d", opts.maxVMsPerPlan, defaultMaxVMsPerPlan)
	}
}

func TestEngineSimpleOptionsApplyValidValues(t *testing.T) {
	t.Parallel()

	registry := encoding.NewRegistry()
	loader := artifact.NewDefaultLoader()
	root := t.TempDir()
	opts := mustNewOptionsForTest(
		t,
		WithEncodingRegistry(registry),
		WithProgramLoader(loader),
		WithMaxActiveSessions(3),
		WithMaxIdleVMsPerPlan(4),
		WithMaxVMsPerPlan(5),
		WithFSRoot(" \t"+root+"\n"),
		WithFSReadOnly(),
	)

	if opts.encoding != registry {
		t.Fatal("expected custom encoding registry")
	}

	if opts.programLoader != loader {
		t.Fatal("expected custom program loader")
	}

	if opts.maxActiveSessions != 3 || opts.maxIdleVMsPerPlan != 4 || opts.maxVMsPerPlan != 5 {
		t.Fatalf(
			"unexpected limits: active=%d idle=%d total=%d",
			opts.maxActiveSessions,
			opts.maxIdleVMsPerPlan,
			opts.maxVMsPerPlan,
		)
	}

	if opts.fsRoot != root {
		t.Fatalf("fs root = %q, want %q", opts.fsRoot, root)
	}

	if !opts.fsReadOnly {
		t.Fatal("expected file system to be read-only")
	}
}

func TestEngineLimitsAcceptZero(t *testing.T) {
	t.Parallel()

	opts := mustNewOptionsForTest(
		t,
		WithMaxActiveSessions(0),
		WithMaxIdleVMsPerPlan(0),
		WithMaxVMsPerPlan(0),
	)

	if opts.maxActiveSessions != 0 || opts.maxIdleVMsPerPlan != 0 || opts.maxVMsPerPlan != 0 {
		t.Fatalf(
			"expected zero limits, got active=%d idle=%d total=%d",
			opts.maxActiveSessions,
			opts.maxIdleVMsPerPlan,
			opts.maxVMsPerPlan,
		)
	}
}

func TestEngineSimpleOptionsReturnStructuredValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		option Option
		name   string
		field  string
		value  string
		reason string
	}{
		{
			name:   "encoding registry",
			field:  "encoding registry",
			value:  "<nil>",
			reason: "cannot be nil",
			option: WithEncodingRegistry(nil),
		},
		{
			name:   "program loader",
			field:  "program loader",
			value:  "<nil>",
			reason: "cannot be nil",
			option: WithProgramLoader(nil),
		},
		{
			name:   "max active sessions",
			field:  "max active sessions",
			value:  "-1",
			reason: "must be non-negative",
			option: WithMaxActiveSessions(-1),
		},
		{
			name:   "max idle VMs per plan",
			field:  "max idle VMs per plan",
			value:  "-1",
			reason: "must be non-negative",
			option: WithMaxIdleVMsPerPlan(-1),
		},
		{
			name:   "max VMs per plan",
			field:  "max VMs per plan",
			value:  "-1",
			reason: "must be non-negative",
			option: WithMaxVMsPerPlan(-1),
		},
		{
			name:   "fs root",
			field:  "fs root",
			value:  `" \t\n "`,
			reason: "must not be blank",
			option: WithFSRoot(" \t\n "),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := newOptions([]Option{tt.option})
			if err == nil {
				t.Fatal("expected validation error")
			}

			var validationErr gooptions.ValidationError
			if !errors.As(err, &validationErr) {
				t.Fatalf("expected options.ValidationError, got %T: %v", err, err)
			}

			if validationErr.Field != tt.field {
				t.Fatalf("validation field = %q, want %q", validationErr.Field, tt.field)
			}

			if validationErr.Value != tt.value {
				t.Fatalf("validation value = %q, want %q", validationErr.Value, tt.value)
			}

			if validationErr.Reason == nil || validationErr.Reason.Error() != tt.reason {
				t.Fatalf("validation reason = %v, want %q", validationErr.Reason, tt.reason)
			}
		})
	}
}

func TestInvalidEngineBuilderOptionDoesNotMutateConfig(t *testing.T) {
	t.Parallel()

	config := defaultOptions()
	registry := config.encoding
	config.maxActiveSessions = 7
	config.fsRoot = "existing"

	for _, option := range []Option{
		WithEncodingRegistry(nil),
		WithMaxActiveSessions(-1),
		WithFSRoot(" \t "),
	} {
		if err := option(&config); err == nil {
			t.Fatal("expected invalid option to fail")
		}
	}

	if config.encoding != registry {
		t.Fatal("invalid registry option mutated the config")
	}

	if config.maxActiveSessions != 7 {
		t.Fatalf("invalid limit mutated the config to %d", config.maxActiveSessions)
	}

	if config.fsRoot != "existing" {
		t.Fatalf("invalid fs root mutated the config to %q", config.fsRoot)
	}
}

func TestNewOptionsAppliesAllOptionsAndJoinsFailuresInOrder(t *testing.T) {
	t.Parallel()

	firstErr := errors.New("first option failed")
	secondErr := errors.New("second option failed")
	var calls []string

	_, err := newOptions([]Option{
		func(*options) error {
			calls = append(calls, "first")

			return firstErr
		},
		nil,
		func(*options) error {
			calls = append(calls, "middle")

			return nil
		},
		func(*options) error {
			calls = append(calls, "second")

			return secondErr
		},
	})
	if err == nil {
		t.Fatal("expected joined option failures")
	}

	if !slices.Equal(calls, []string{"first", "middle", "second"}) {
		t.Fatalf("option calls = %v", calls)
	}

	if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
		t.Fatalf("expected both option failures, got %v", err)
	}

	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected joined error, got %T", err)
	}

	failures := joined.Unwrap()
	if len(failures) != 2 || failures[0] != firstErr || failures[1] != secondErr {
		t.Fatalf("joined failures = %v, want [%v %v]", failures, firstErr, secondErr)
	}
}

func TestJoinedEngineOptionFailuresRemainInspectable(t *testing.T) {
	t.Parallel()

	_, err := newOptions([]Option{
		WithMaxActiveSessions(-1),
		WithParams(map[string]any{"unsupported": make(chan int)}),
	})
	if err == nil {
		t.Fatal("expected joined option failures")
	}

	if !errors.Is(err, runtime.ErrInvalidType) {
		t.Fatalf("expected runtime.ErrInvalidType, got %v", err)
	}

	var validationErr gooptions.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected options.ValidationError, got %T: %v", err, err)
	}

	if validationErr.Field != "max active sessions" {
		t.Fatalf("validation field = %q", validationErr.Field)
	}
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

	var validationErr gooptions.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected blank fs root validation error, got %T: %v", err, err)
	}

	if validationErr.Field != "fs root" {
		t.Fatalf("validation field = %q, want %q", validationErr.Field, "fs root")
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

	var validationErr gooptions.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected options.ValidationError, got %T: %v", err, err)
	}
	if validationErr.Field != "max response size" || validationErr.Value != "-1" {
		t.Fatalf("unexpected validation error: %+v", validationErr)
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

func TestWithFunctionsRegistrarPreservesQualifiedHostFunctionNames(t *testing.T) {
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

	if got := eng.host.functions.List(); !slices.Equal(got, []string{"Tools::Risk::Calculate_Risk"}) {
		t.Fatalf("expected declared host metadata, got %v", got)
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

func TestResolvedHostFunctionDiagnosticUsesRegisteredQualifiedName(t *testing.T) {
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

	if got, want := runtimeErr.Note, "DB::POSTGRES::QUERY expects 1 argument, but got 0"; got != want {
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
