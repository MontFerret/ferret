package ferret_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/module"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	publicExport struct {
		root   any
		owner  any
		name   string
		target string
		kind   token.Token
	}

	exportDeclaration struct {
		target string
		kind   token.Token
		alias  bool
	}
)

func TestPublicExportContracts(t *testing.T) {
	for _, export := range publicExports() {
		t.Run(export.name, func(t *testing.T) {
			if export.kind == token.FUNC {
				root, owner := reflect.TypeOf(export.root), reflect.TypeOf(export.owner)
				if root != owner {
					t.Fatalf("function signature = %v, want %v", root, owner)
				}

				return
			}

			if export.root != export.owner {
				t.Fatalf("root = %v (%T), owner = %v (%T)", export.root, export.root, export.owner, export.owner)
			}
		})
	}
}

func TestPublicExportDeclarations(t *testing.T) {
	actual := packageExportDeclarations(t, ".")
	seen := make(map[string]bool)
	for _, export := range publicExports() {
		if seen[export.name] {
			t.Fatalf("duplicate API contract for %s", export.name)
		}

		seen[export.name] = true

		got, exists := actual[export.name]
		if !exists {
			t.Errorf("missing root export %s", export.name)
		} else {
			if got.kind != export.kind {
				t.Errorf("%s is declared as %s, want %s", export.name, got.kind, export.kind)
			}

			if export.kind == token.TYPE && !got.alias {
				t.Errorf("%s must remain a type alias", export.name)
			}

			if got.target != export.target {
				t.Errorf("%s directly targets %q, want %q", export.name, got.target, export.target)
			}
		}

		delete(actual, export.name)
	}

	for name := range actual {
		t.Errorf("unaudited root export %s", name)
	}
}

func TestNativeExportDeclarations(t *testing.T) {
	// Native exports are independently curated; additions do not imply root additions.
	expected := map[token.Token][]string{
		token.TYPE: {
			"Engine",
			"OptimizationLevel",
			"Option",
			"Plan",
			"PlanOption",
			"Session",
			"SessionOption",
		},
		token.CONST: {
			"OptimizationBasic",
			"OptimizationFull",
			"OptimizationNone",
		},
		token.FUNC: {
			"New",
			"WithAfterCompileHook",
			"WithAfterRunHook",
			"WithBeforeCompileHook",
			"WithBeforeRunHook",
			"WithDebugFormat",
			"WithEncodingCodec",
			"WithEncodingRegistry",
			"WithEngineCloseHook",
			"WithEngineInitHook",
			"WithEnvironmentOptions",
			"WithFSReadOnly",
			"WithFSRoot",
			"WithFunctions",
			"WithFunctionsRegistrar",
			"WithLog",
			"WithLogFields",
			"WithLogLevel",
			"WithMaxActiveSessions",
			"WithMaxIdleVMsPerPlan",
			"WithMaxVMsPerPlan",
			"WithModules",
			"WithNamespace",
			"WithNetwork",
			"WithNetworkOptions",
			"WithOptimizationLevel",
			"WithOutputContentType",
			"WithParam",
			"WithParams",
			"WithPlanCloseHook",
			"WithPlanOptimizationLevel",
			"WithProgramLoader",
			"WithRuntimeParam",
			"WithRuntimeParams",
			"WithSessionCloseHook",
			"WithSessionFSRoot",
			"WithSessionLog",
			"WithSessionLogFields",
			"WithSessionLogLevel",
			"WithSessionParam",
			"WithSessionParams",
			"WithSessionRuntimeParam",
			"WithSessionRuntimeParams",
			"WithStdlib",
			"WithoutStdlib",
		},
	}
	actual := packageExportDeclarations(t, "pkg/engine")
	for kind, names := range expected {
		for _, name := range names {
			got, exists := actual[name]
			if !exists {
				t.Errorf("missing native export %s", name)
			} else if got.kind != kind {
				t.Errorf("%s is declared as %s, want %s", name, got.kind, kind)
			}

			delete(actual, name)
		}
	}

	for name := range actual {
		t.Errorf("unaudited native export %s", name)
	}
}

// Go runs package tests in their source directory, including with -trimpath.
func packageExportDeclarations(t *testing.T, directory string) map[string]exportDeclaration {
	t.Helper()

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}

	actual := make(map[string]exportDeclaration)
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}

		file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(directory, name), nil, 0)
		if err != nil {
			t.Fatal(err)
		}

		imports := make(map[string]string)
		for _, spec := range file.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				t.Fatal(err)
			}

			name := filepath.Base(importPath)
			if spec.Name != nil {
				name = spec.Name.Name
			}

			imports[name] = importPath
		}

		for _, declaration := range file.Decls {
			switch declaration := declaration.(type) {
			case *ast.FuncDecl:
				if declaration.Recv == nil && declaration.Name.IsExported() {
					actual[declaration.Name.Name] = exportDeclaration{
						kind:   token.FUNC,
						target: forwardingTarget(imports, declaration),
					}
				}
			case *ast.GenDecl:
				for _, spec := range declaration.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						if spec.Name.IsExported() {
							actual[spec.Name.Name] = exportDeclaration{
								kind:   token.TYPE,
								alias:  spec.Assign.IsValid(),
								target: importedTarget(imports, spec.Type),
							}
						}
					case *ast.ValueSpec:
						for index, name := range spec.Names {
							if name.IsExported() {
								export := exportDeclaration{kind: declaration.Tok}
								if index < len(spec.Values) {
									export.target = importedTarget(imports, spec.Values[index])
								}

								actual[name.Name] = export
							}
						}
					}
				}
			}
		}
	}

	return actual
}

func forwardingTarget(imports map[string]string, declaration *ast.FuncDecl) string {
	if declaration.Body == nil || len(declaration.Body.List) != 1 {
		return ""
	}

	result, ok := declaration.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(result.Results) != 1 {
		return ""
	}

	call, ok := result.Results[0].(*ast.CallExpr)
	if !ok {
		return ""
	}

	return importedTarget(imports, call.Fun)
}

func importedTarget(imports map[string]string, expression ast.Expr) string {
	selector, ok := expression.(*ast.SelectorExpr)
	if !ok {
		return ""
	}

	qualifier, ok := selector.X.(*ast.Ident)
	if !ok {
		return ""
	}

	importPath, ok := imports[qualifier.Name]
	if !ok {
		return ""
	}

	return strings.TrimPrefix(importPath, "github.com/MontFerret/ferret/v2/") + "." + selector.Sel.Name
}

// This is the curated pre-extraction surface at 1ae9a44f, with FormatError
// intentionally changed from a function variable to a declared function.
// Keep it explicit: future native exports do not automatically belong at root.
func publicExports() []publicExport {
	return []publicExport{
		{name: "AfterCompileHook", kind: token.TYPE, target: "pkg/module.AfterCompileHook", root: reflect.TypeFor[ferret.AfterCompileHook](), owner: reflect.TypeFor[module.AfterCompileHook]()},
		{name: "AfterRunHook", kind: token.TYPE, target: "pkg/module.AfterRunHook", root: reflect.TypeFor[ferret.AfterRunHook](), owner: reflect.TypeFor[module.AfterRunHook]()},
		{name: "BeforeCompileHook", kind: token.TYPE, target: "pkg/module.BeforeCompileHook", root: reflect.TypeFor[ferret.BeforeCompileHook](), owner: reflect.TypeFor[module.BeforeCompileHook]()},
		{name: "BeforeRunHook", kind: token.TYPE, target: "pkg/module.BeforeRunHook", root: reflect.TypeFor[ferret.BeforeRunHook](), owner: reflect.TypeFor[module.BeforeRunHook]()},
		{name: "DebugBreakpoint", kind: token.TYPE, target: "pkg/debugger.Breakpoint", root: reflect.TypeFor[ferret.DebugBreakpoint](), owner: reflect.TypeFor[debugger.Breakpoint]()},
		{name: "DebugBreakpointBindingMode", kind: token.TYPE, target: "pkg/debugger.BreakpointBindingMode", root: reflect.TypeFor[ferret.DebugBreakpointBindingMode](), owner: reflect.TypeFor[debugger.BreakpointBindingMode]()},
		{name: "DebugBreakpointID", kind: token.TYPE, target: "pkg/debugger.BreakpointID", root: reflect.TypeFor[ferret.DebugBreakpointID](), owner: reflect.TypeFor[debugger.BreakpointID]()},
		{name: "DebugBreakpointOptions", kind: token.TYPE, target: "pkg/debugger.BreakpointOptions", root: reflect.TypeFor[ferret.DebugBreakpointOptions](), owner: reflect.TypeFor[debugger.BreakpointOptions]()},
		{name: "DebugEvent", kind: token.TYPE, target: "pkg/debugger.Event", root: reflect.TypeFor[ferret.DebugEvent](), owner: reflect.TypeFor[debugger.Event]()},
		{name: "DebugFormatOptions", kind: token.TYPE, target: "pkg/debugger.FormatOptions", root: reflect.TypeFor[ferret.DebugFormatOptions](), owner: reflect.TypeFor[debugger.FormatOptions]()},
		{name: "DebugFrame", kind: token.TYPE, target: "pkg/debugger.Frame", root: reflect.TypeFor[ferret.DebugFrame](), owner: reflect.TypeFor[debugger.Frame]()},
		{name: "DebugLocation", kind: token.TYPE, target: "pkg/source.Range", root: reflect.TypeFor[ferret.DebugLocation](), owner: reflect.TypeFor[source.Range]()},
		{name: "DebugReason", kind: token.TYPE, target: "pkg/debugger.Reason", root: reflect.TypeFor[ferret.DebugReason](), owner: reflect.TypeFor[debugger.Reason]()},
		{name: "DebugSession", kind: token.TYPE, target: "pkg/debugger.Session", root: reflect.TypeFor[ferret.DebugSession](), owner: reflect.TypeFor[debugger.Session]()},
		{name: "DebugSourceLocation", kind: token.TYPE, target: "pkg/source.Location", root: reflect.TypeFor[ferret.DebugSourceLocation](), owner: reflect.TypeFor[source.Location]()},
		{name: "DebugStateError", kind: token.TYPE, target: "pkg/debugger.StateError", root: reflect.TypeFor[ferret.DebugStateError](), owner: reflect.TypeFor[debugger.StateError]()},
		{name: "DebugValue", kind: token.TYPE, target: "pkg/debugger.Value", root: reflect.TypeFor[ferret.DebugValue](), owner: reflect.TypeFor[debugger.Value]()},
		{name: "DebugValueReference", kind: token.TYPE, target: "pkg/debugger.ValueReference", root: reflect.TypeFor[ferret.DebugValueReference](), owner: reflect.TypeFor[debugger.ValueReference]()},
		{name: "DebugVariable", kind: token.TYPE, target: "pkg/debugger.Variable", root: reflect.TypeFor[ferret.DebugVariable](), owner: reflect.TypeFor[debugger.Variable]()},
		{name: "Engine", kind: token.TYPE, target: "pkg/engine.Engine", root: reflect.TypeFor[ferret.Engine](), owner: reflect.TypeFor[engine.Engine]()},
		{name: "EngineCloseHook", kind: token.TYPE, target: "pkg/module.EngineCloseHook", root: reflect.TypeFor[ferret.EngineCloseHook](), owner: reflect.TypeFor[module.EngineCloseHook]()},
		{name: "EngineInitHook", kind: token.TYPE, target: "pkg/module.EngineInitHook", root: reflect.TypeFor[ferret.EngineInitHook](), owner: reflect.TypeFor[module.EngineInitHook]()},
		{name: "Location", kind: token.TYPE, target: "pkg/source.Location", root: reflect.TypeFor[ferret.Location](), owner: reflect.TypeFor[source.Location]()},
		{name: "LogLevel", kind: token.TYPE, target: "pkg/logging.LogLevel", root: reflect.TypeFor[ferret.LogLevel](), owner: reflect.TypeFor[logging.LogLevel]()},
		{name: "Module", kind: token.TYPE, target: "pkg/module.Module", root: reflect.TypeFor[ferret.Module](), owner: reflect.TypeFor[module.Module]()},
		{name: "OptimizationLevel", kind: token.TYPE, target: "pkg/engine.OptimizationLevel", root: reflect.TypeFor[ferret.OptimizationLevel](), owner: reflect.TypeFor[engine.OptimizationLevel]()},
		{name: "Option", kind: token.TYPE, target: "pkg/engine.Option", root: reflect.TypeFor[ferret.Option](), owner: reflect.TypeFor[engine.Option]()},
		{name: "Output", kind: token.TYPE, target: "pkg/encoding.Output", root: reflect.TypeFor[ferret.Output](), owner: reflect.TypeFor[encoding.Output]()},
		{name: "Params", kind: token.TYPE, target: "pkg/runtime.Params", root: reflect.TypeFor[ferret.Params](), owner: reflect.TypeFor[runtime.Params]()},
		{name: "Plan", kind: token.TYPE, target: "pkg/engine.Plan", root: reflect.TypeFor[ferret.Plan](), owner: reflect.TypeFor[engine.Plan]()},
		{name: "PlanCloseHook", kind: token.TYPE, target: "pkg/module.PlanCloseHook", root: reflect.TypeFor[ferret.PlanCloseHook](), owner: reflect.TypeFor[module.PlanCloseHook]()},
		{name: "PlanOption", kind: token.TYPE, target: "pkg/engine.PlanOption", root: reflect.TypeFor[ferret.PlanOption](), owner: reflect.TypeFor[engine.PlanOption]()},
		{name: "Position", kind: token.TYPE, target: "pkg/source.Position", root: reflect.TypeFor[ferret.Position](), owner: reflect.TypeFor[source.Position]()},
		{name: "ProgramFormat", kind: token.TYPE, target: "pkg/bytecode/artifact.FormatID", root: reflect.TypeFor[ferret.ProgramFormat](), owner: reflect.TypeFor[artifact.FormatID]()},
		{name: "ProgramOption", kind: token.TYPE, target: "pkg/bytecode/artifact.Option", root: reflect.TypeFor[ferret.ProgramOption](), owner: reflect.TypeFor[artifact.Option]()},
		{name: "Range", kind: token.TYPE, target: "pkg/source.Range", root: reflect.TypeFor[ferret.Range](), owner: reflect.TypeFor[source.Range]()},
		{name: "Session", kind: token.TYPE, target: "pkg/engine.Session", root: reflect.TypeFor[ferret.Session](), owner: reflect.TypeFor[engine.Session]()},
		{name: "SessionCloseHook", kind: token.TYPE, target: "pkg/module.SessionCloseHook", root: reflect.TypeFor[ferret.SessionCloseHook](), owner: reflect.TypeFor[module.SessionCloseHook]()},
		{name: "SessionOption", kind: token.TYPE, target: "pkg/engine.SessionOption", root: reflect.TypeFor[ferret.SessionOption](), owner: reflect.TypeFor[engine.SessionOption]()},
		{name: "Source", kind: token.TYPE, target: "pkg/source.Source", root: reflect.TypeFor[ferret.Source](), owner: reflect.TypeFor[source.Source]()},
		{name: "Span", kind: token.TYPE, target: "pkg/source.Span", root: reflect.TypeFor[ferret.Span](), owner: reflect.TypeFor[source.Span]()},
		{name: "Value", kind: token.TYPE, target: "pkg/runtime.Value", root: reflect.TypeFor[ferret.Value](), owner: reflect.TypeFor[runtime.Value]()},

		{name: "DebugBreakpointBindExact", kind: token.CONST, target: "pkg/debugger.BreakpointBindExact", root: ferret.DebugBreakpointBindExact, owner: debugger.BreakpointBindExact},
		{name: "DebugBreakpointBindNextExecutableInFunction", kind: token.CONST, target: "pkg/debugger.BreakpointBindNextExecutableInFunction", root: ferret.DebugBreakpointBindNextExecutableInFunction, owner: debugger.BreakpointBindNextExecutableInFunction},
		{name: "DebugBreakpointBindNextExecutableInSource", kind: token.CONST, target: "pkg/debugger.BreakpointBindNextExecutableInSource", root: ferret.DebugBreakpointBindNextExecutableInSource, owner: debugger.BreakpointBindNextExecutableInSource},
		{name: "DebugReasonBreakpoint", kind: token.CONST, target: "pkg/debugger.ReasonBreakpoint", root: ferret.DebugReasonBreakpoint, owner: debugger.ReasonBreakpoint},
		{name: "DebugReasonCompleted", kind: token.CONST, target: "pkg/debugger.ReasonCompleted", root: ferret.DebugReasonCompleted, owner: debugger.ReasonCompleted},
		{name: "DebugReasonEntry", kind: token.CONST, target: "pkg/debugger.ReasonEntry", root: ferret.DebugReasonEntry, owner: debugger.ReasonEntry},
		{name: "DebugReasonPause", kind: token.CONST, target: "pkg/debugger.ReasonPause", root: ferret.DebugReasonPause, owner: debugger.ReasonPause},
		{name: "DebugReasonRuntimeError", kind: token.CONST, target: "pkg/debugger.ReasonRuntimeError", root: ferret.DebugReasonRuntimeError, owner: debugger.ReasonRuntimeError},
		{name: "DebugReasonStep", kind: token.CONST, target: "pkg/debugger.ReasonStep", root: ferret.DebugReasonStep, owner: debugger.ReasonStep},
		{name: "DebugReasonTerminated", kind: token.CONST, target: "pkg/debugger.ReasonTerminated", root: ferret.DebugReasonTerminated, owner: debugger.ReasonTerminated},
		{name: "LogDebug", kind: token.CONST, target: "pkg/logging.DebugLevel", root: ferret.LogDebug, owner: logging.DebugLevel},
		{name: "LogDisabled", kind: token.CONST, target: "pkg/logging.Disabled", root: ferret.LogDisabled, owner: logging.Disabled},
		{name: "LogError", kind: token.CONST, target: "pkg/logging.ErrorLevel", root: ferret.LogError, owner: logging.ErrorLevel},
		{name: "LogFatal", kind: token.CONST, target: "pkg/logging.FatalLevel", root: ferret.LogFatal, owner: logging.FatalLevel},
		{name: "LogInfo", kind: token.CONST, target: "pkg/logging.InfoLevel", root: ferret.LogInfo, owner: logging.InfoLevel},
		{name: "LogNone", kind: token.CONST, target: "pkg/logging.NoLevel", root: ferret.LogNone, owner: logging.NoLevel},
		{name: "LogPanic", kind: token.CONST, target: "pkg/logging.PanicLevel", root: ferret.LogPanic, owner: logging.PanicLevel},
		{name: "LogTrace", kind: token.CONST, target: "pkg/logging.TraceLevel", root: ferret.LogTrace, owner: logging.TraceLevel},
		{name: "LogWarn", kind: token.CONST, target: "pkg/logging.WarnLevel", root: ferret.LogWarn, owner: logging.WarnLevel},
		{name: "OptimizationBasic", kind: token.CONST, target: "pkg/engine.OptimizationBasic", root: ferret.OptimizationBasic, owner: engine.OptimizationBasic},
		{name: "OptimizationFull", kind: token.CONST, target: "pkg/engine.OptimizationFull", root: ferret.OptimizationFull, owner: engine.OptimizationFull},
		{name: "OptimizationNone", kind: token.CONST, target: "pkg/engine.OptimizationNone", root: ferret.OptimizationNone, owner: engine.OptimizationNone},
		{name: "ProgramFormatJSON", kind: token.CONST, target: "pkg/bytecode/artifact.FormatJSON", root: ferret.ProgramFormatJSON, owner: artifact.FormatJSON},
		{name: "ProgramFormatMsgPack", kind: token.CONST, target: "pkg/bytecode/artifact.FormatMsgPack", root: ferret.ProgramFormatMsgPack, owner: artifact.FormatMsgPack},

		{name: "FormatError", kind: token.FUNC, target: "pkg/diagnostics.Format", root: ferret.FormatError, owner: diagnostics.Format},
		{name: "MarshalProgram", kind: token.FUNC, target: "pkg/bytecode/artifact.Marshal", root: ferret.MarshalProgram, owner: artifact.Marshal},
		{name: "MustParseLogLevel", kind: token.FUNC, target: "pkg/logging.MustParseLogLevel", root: ferret.MustParseLogLevel, owner: logging.MustParseLogLevel},
		{name: "New", kind: token.FUNC, target: "pkg/engine.New", root: ferret.New, owner: engine.New},
		{name: "NewAnonymousSource", kind: token.FUNC, target: "pkg/source.NewAnonymous", root: ferret.NewAnonymousSource, owner: source.NewAnonymous},
		{name: "NewSource", kind: token.FUNC, target: "pkg/source.New", root: ferret.NewSource, owner: source.New},
		{name: "ParseLogLevel", kind: token.FUNC, target: "pkg/logging.ParseLogLevel", root: ferret.ParseLogLevel, owner: logging.ParseLogLevel},
		{name: "UnmarshalProgram", kind: token.FUNC, target: "pkg/bytecode/artifact.Unmarshal", root: ferret.UnmarshalProgram, owner: artifact.Unmarshal},
		{name: "WithAfterCompileHook", kind: token.FUNC, target: "pkg/engine.WithAfterCompileHook", root: ferret.WithAfterCompileHook, owner: engine.WithAfterCompileHook},
		{name: "WithAfterRunHook", kind: token.FUNC, target: "pkg/engine.WithAfterRunHook", root: ferret.WithAfterRunHook, owner: engine.WithAfterRunHook},
		{name: "WithBeforeCompileHook", kind: token.FUNC, target: "pkg/engine.WithBeforeCompileHook", root: ferret.WithBeforeCompileHook, owner: engine.WithBeforeCompileHook},
		{name: "WithBeforeRunHook", kind: token.FUNC, target: "pkg/engine.WithBeforeRunHook", root: ferret.WithBeforeRunHook, owner: engine.WithBeforeRunHook},
		{name: "WithDebugFormat", kind: token.FUNC, target: "pkg/engine.WithDebugFormat", root: ferret.WithDebugFormat, owner: engine.WithDebugFormat},
		{name: "WithEncodingCodec", kind: token.FUNC, target: "pkg/engine.WithEncodingCodec", root: ferret.WithEncodingCodec, owner: engine.WithEncodingCodec},
		{name: "WithEncodingRegistry", kind: token.FUNC, target: "pkg/engine.WithEncodingRegistry", root: ferret.WithEncodingRegistry, owner: engine.WithEncodingRegistry},
		{name: "WithEngineCloseHook", kind: token.FUNC, target: "pkg/engine.WithEngineCloseHook", root: ferret.WithEngineCloseHook, owner: engine.WithEngineCloseHook},
		{name: "WithEngineInitHook", kind: token.FUNC, target: "pkg/engine.WithEngineInitHook", root: ferret.WithEngineInitHook, owner: engine.WithEngineInitHook},
		{name: "WithEnvironmentOptions", kind: token.FUNC, target: "pkg/engine.WithEnvironmentOptions", root: ferret.WithEnvironmentOptions, owner: engine.WithEnvironmentOptions},
		{name: "WithFSReadOnly", kind: token.FUNC, target: "pkg/engine.WithFSReadOnly", root: ferret.WithFSReadOnly, owner: engine.WithFSReadOnly},
		{name: "WithFSRoot", kind: token.FUNC, target: "pkg/engine.WithFSRoot", root: ferret.WithFSRoot, owner: engine.WithFSRoot},
		{name: "WithFunctions", kind: token.FUNC, target: "pkg/engine.WithFunctions", root: ferret.WithFunctions, owner: engine.WithFunctions},
		{name: "WithFunctionsRegistrar", kind: token.FUNC, target: "pkg/engine.WithFunctionsRegistrar", root: ferret.WithFunctionsRegistrar, owner: engine.WithFunctionsRegistrar},
		{name: "WithLog", kind: token.FUNC, target: "pkg/engine.WithLog", root: ferret.WithLog, owner: engine.WithLog},
		{name: "WithLogFields", kind: token.FUNC, target: "pkg/engine.WithLogFields", root: ferret.WithLogFields, owner: engine.WithLogFields},
		{name: "WithLogLevel", kind: token.FUNC, target: "pkg/engine.WithLogLevel", root: ferret.WithLogLevel, owner: engine.WithLogLevel},
		{name: "WithMaxActiveSessions", kind: token.FUNC, target: "pkg/engine.WithMaxActiveSessions", root: ferret.WithMaxActiveSessions, owner: engine.WithMaxActiveSessions},
		{name: "WithMaxIdleVMsPerPlan", kind: token.FUNC, target: "pkg/engine.WithMaxIdleVMsPerPlan", root: ferret.WithMaxIdleVMsPerPlan, owner: engine.WithMaxIdleVMsPerPlan},
		{name: "WithMaxVMsPerPlan", kind: token.FUNC, target: "pkg/engine.WithMaxVMsPerPlan", root: ferret.WithMaxVMsPerPlan, owner: engine.WithMaxVMsPerPlan},
		{name: "WithModules", kind: token.FUNC, target: "pkg/engine.WithModules", root: ferret.WithModules, owner: engine.WithModules},
		{name: "WithNamespace", kind: token.FUNC, target: "pkg/engine.WithNamespace", root: ferret.WithNamespace, owner: engine.WithNamespace},
		{name: "WithNetwork", kind: token.FUNC, target: "pkg/engine.WithNetwork", root: ferret.WithNetwork, owner: engine.WithNetwork},
		{name: "WithNetworkOptions", kind: token.FUNC, target: "pkg/engine.WithNetworkOptions", root: ferret.WithNetworkOptions, owner: engine.WithNetworkOptions},
		{name: "WithOptimizationLevel", kind: token.FUNC, target: "pkg/engine.WithOptimizationLevel", root: ferret.WithOptimizationLevel, owner: engine.WithOptimizationLevel},
		{name: "WithOutputContentType", kind: token.FUNC, target: "pkg/engine.WithOutputContentType", root: ferret.WithOutputContentType, owner: engine.WithOutputContentType},
		{name: "WithParam", kind: token.FUNC, target: "pkg/engine.WithParam", root: ferret.WithParam, owner: engine.WithParam},
		{name: "WithParams", kind: token.FUNC, target: "pkg/engine.WithParams", root: ferret.WithParams, owner: engine.WithParams},
		{name: "WithPlanCloseHook", kind: token.FUNC, target: "pkg/engine.WithPlanCloseHook", root: ferret.WithPlanCloseHook, owner: engine.WithPlanCloseHook},
		{name: "WithPlanOptimizationLevel", kind: token.FUNC, target: "pkg/engine.WithPlanOptimizationLevel", root: ferret.WithPlanOptimizationLevel, owner: engine.WithPlanOptimizationLevel},
		{name: "WithProgramFormat", kind: token.FUNC, target: "pkg/bytecode/artifact.WithFormat", root: ferret.WithProgramFormat, owner: artifact.WithFormat},
		{name: "WithProgramLoader", kind: token.FUNC, target: "pkg/engine.WithProgramLoader", root: ferret.WithProgramLoader, owner: engine.WithProgramLoader},
		{name: "WithRuntimeParam", kind: token.FUNC, target: "pkg/engine.WithRuntimeParam", root: ferret.WithRuntimeParam, owner: engine.WithRuntimeParam},
		{name: "WithRuntimeParams", kind: token.FUNC, target: "pkg/engine.WithRuntimeParams", root: ferret.WithRuntimeParams, owner: engine.WithRuntimeParams},
		{name: "WithSessionCloseHook", kind: token.FUNC, target: "pkg/engine.WithSessionCloseHook", root: ferret.WithSessionCloseHook, owner: engine.WithSessionCloseHook},
		{name: "WithSessionFSRoot", kind: token.FUNC, target: "pkg/engine.WithSessionFSRoot", root: ferret.WithSessionFSRoot, owner: engine.WithSessionFSRoot},
		{name: "WithSessionLog", kind: token.FUNC, target: "pkg/engine.WithSessionLog", root: ferret.WithSessionLog, owner: engine.WithSessionLog},
		{name: "WithSessionLogFields", kind: token.FUNC, target: "pkg/engine.WithSessionLogFields", root: ferret.WithSessionLogFields, owner: engine.WithSessionLogFields},
		{name: "WithSessionLogLevel", kind: token.FUNC, target: "pkg/engine.WithSessionLogLevel", root: ferret.WithSessionLogLevel, owner: engine.WithSessionLogLevel},
		{name: "WithSessionParam", kind: token.FUNC, target: "pkg/engine.WithSessionParam", root: ferret.WithSessionParam, owner: engine.WithSessionParam},
		{name: "WithSessionParams", kind: token.FUNC, target: "pkg/engine.WithSessionParams", root: ferret.WithSessionParams, owner: engine.WithSessionParams},
		{name: "WithSessionRuntimeParam", kind: token.FUNC, target: "pkg/engine.WithSessionRuntimeParam", root: ferret.WithSessionRuntimeParam, owner: engine.WithSessionRuntimeParam},
		{name: "WithSessionRuntimeParams", kind: token.FUNC, target: "pkg/engine.WithSessionRuntimeParams", root: ferret.WithSessionRuntimeParams, owner: engine.WithSessionRuntimeParams},
		{name: "WithStdlib", kind: token.FUNC, target: "pkg/engine.WithStdlib", root: ferret.WithStdlib, owner: engine.WithStdlib},
		{name: "WithoutStdlib", kind: token.FUNC, target: "pkg/engine.WithoutStdlib", root: ferret.WithoutStdlib, owner: engine.WithoutStdlib},
	}
}
