package ferret

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
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
