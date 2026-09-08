package ferret_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/engine"
)

type publicExport struct {
	root   any
	native any
	name   string
	kind   token.Token
}

func TestPublicExportContracts(t *testing.T) {
	for _, export := range publicExports() {
		t.Run(export.name, func(t *testing.T) {
			if export.kind == token.FUNC {
				root, native := reflect.TypeOf(export.root), reflect.TypeOf(export.native)
				if root != native {
					t.Fatalf("function signature = %v, want %v", root, native)
				}

				return
			}

			if export.root != export.native {
				t.Fatalf("root = %v (%T), native = %v (%T)", export.root, export.root, export.native, export.native)
			}
		})
	}
}

func TestPublicExportDeclarations(t *testing.T) {
	expected := make(map[string]token.Token)
	for _, export := range publicExports() {
		if _, exists := expected[export.name]; exists {
			t.Fatalf("duplicate API contract for %s", export.name)
		}

		expected[export.name] = export.kind
	}

	// Go runs package tests in their source directory, including with -trimpath.
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}

	actual := make(map[string]token.Token)
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}

		file, err := parser.ParseFile(token.NewFileSet(), name, nil, 0)
		if err != nil {
			t.Fatal(err)
		}

		for _, declaration := range file.Decls {
			switch declaration := declaration.(type) {
			case *ast.FuncDecl:
				if declaration.Recv == nil && declaration.Name.IsExported() {
					actual[declaration.Name.Name] = token.FUNC
				}
			case *ast.GenDecl:
				for _, spec := range declaration.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						if spec.Name.IsExported() {
							actual[spec.Name.Name] = token.TYPE

							if !spec.Assign.IsValid() {
								t.Errorf("%s must remain a type alias", spec.Name.Name)
							}
						}
					case *ast.ValueSpec:
						for _, name := range spec.Names {
							if name.IsExported() {
								actual[name.Name] = declaration.Tok
							}
						}
					}
				}
			}
		}
	}

	for name, want := range expected {
		got, exists := actual[name]
		if !exists {
			t.Errorf("missing root export %s", name)
		} else if got != want {
			t.Errorf("%s is declared as %s, want %s", name, got, want)
		}

		delete(actual, name)
	}

	for name := range actual {
		t.Errorf("unaudited root export %s", name)
	}
}

// This is the curated pre-extraction surface at 1ae9a44f, with FormatError
// intentionally changed from a function variable to a declared function.
// Keep it explicit: future native exports do not automatically belong at root.
func publicExports() []publicExport {
	return []publicExport{
		{name: "AfterCompileHook", kind: token.TYPE, root: reflect.TypeFor[ferret.AfterCompileHook](), native: reflect.TypeFor[engine.AfterCompileHook]()},
		{name: "AfterRunHook", kind: token.TYPE, root: reflect.TypeFor[ferret.AfterRunHook](), native: reflect.TypeFor[engine.AfterRunHook]()},
		{name: "BeforeCompileHook", kind: token.TYPE, root: reflect.TypeFor[ferret.BeforeCompileHook](), native: reflect.TypeFor[engine.BeforeCompileHook]()},
		{name: "BeforeRunHook", kind: token.TYPE, root: reflect.TypeFor[ferret.BeforeRunHook](), native: reflect.TypeFor[engine.BeforeRunHook]()},
		{name: "DebugBreakpoint", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugBreakpoint](), native: reflect.TypeFor[engine.DebugBreakpoint]()},
		{name: "DebugBreakpointBindingMode", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugBreakpointBindingMode](), native: reflect.TypeFor[engine.DebugBreakpointBindingMode]()},
		{name: "DebugBreakpointID", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugBreakpointID](), native: reflect.TypeFor[engine.DebugBreakpointID]()},
		{name: "DebugBreakpointOptions", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugBreakpointOptions](), native: reflect.TypeFor[engine.DebugBreakpointOptions]()},
		{name: "DebugEvent", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugEvent](), native: reflect.TypeFor[engine.DebugEvent]()},
		{name: "DebugFormatOptions", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugFormatOptions](), native: reflect.TypeFor[engine.DebugFormatOptions]()},
		{name: "DebugFrame", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugFrame](), native: reflect.TypeFor[engine.DebugFrame]()},
		{name: "DebugLocation", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugLocation](), native: reflect.TypeFor[engine.DebugLocation]()},
		{name: "DebugReason", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugReason](), native: reflect.TypeFor[engine.DebugReason]()},
		{name: "DebugSession", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugSession](), native: reflect.TypeFor[engine.DebugSession]()},
		{name: "DebugSourceLocation", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugSourceLocation](), native: reflect.TypeFor[engine.DebugSourceLocation]()},
		{name: "DebugStateError", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugStateError](), native: reflect.TypeFor[engine.DebugStateError]()},
		{name: "DebugValue", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugValue](), native: reflect.TypeFor[engine.DebugValue]()},
		{name: "DebugValueReference", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugValueReference](), native: reflect.TypeFor[engine.DebugValueReference]()},
		{name: "DebugVariable", kind: token.TYPE, root: reflect.TypeFor[ferret.DebugVariable](), native: reflect.TypeFor[engine.DebugVariable]()},
		{name: "Engine", kind: token.TYPE, root: reflect.TypeFor[ferret.Engine](), native: reflect.TypeFor[engine.Engine]()},
		{name: "EngineCloseHook", kind: token.TYPE, root: reflect.TypeFor[ferret.EngineCloseHook](), native: reflect.TypeFor[engine.EngineCloseHook]()},
		{name: "EngineInitHook", kind: token.TYPE, root: reflect.TypeFor[ferret.EngineInitHook](), native: reflect.TypeFor[engine.EngineInitHook]()},
		{name: "Location", kind: token.TYPE, root: reflect.TypeFor[ferret.Location](), native: reflect.TypeFor[engine.Location]()},
		{name: "LogLevel", kind: token.TYPE, root: reflect.TypeFor[ferret.LogLevel](), native: reflect.TypeFor[engine.LogLevel]()},
		{name: "Module", kind: token.TYPE, root: reflect.TypeFor[ferret.Module](), native: reflect.TypeFor[engine.Module]()},
		{name: "OptimizationLevel", kind: token.TYPE, root: reflect.TypeFor[ferret.OptimizationLevel](), native: reflect.TypeFor[engine.OptimizationLevel]()},
		{name: "Option", kind: token.TYPE, root: reflect.TypeFor[ferret.Option](), native: reflect.TypeFor[engine.Option]()},
		{name: "Output", kind: token.TYPE, root: reflect.TypeFor[ferret.Output](), native: reflect.TypeFor[engine.Output]()},
		{name: "Params", kind: token.TYPE, root: reflect.TypeFor[ferret.Params](), native: reflect.TypeFor[engine.Params]()},
		{name: "Plan", kind: token.TYPE, root: reflect.TypeFor[ferret.Plan](), native: reflect.TypeFor[engine.Plan]()},
		{name: "PlanCloseHook", kind: token.TYPE, root: reflect.TypeFor[ferret.PlanCloseHook](), native: reflect.TypeFor[engine.PlanCloseHook]()},
		{name: "PlanOption", kind: token.TYPE, root: reflect.TypeFor[ferret.PlanOption](), native: reflect.TypeFor[engine.PlanOption]()},
		{name: "Position", kind: token.TYPE, root: reflect.TypeFor[ferret.Position](), native: reflect.TypeFor[engine.Position]()},
		{name: "ProgramFormat", kind: token.TYPE, root: reflect.TypeFor[ferret.ProgramFormat](), native: reflect.TypeFor[engine.ProgramFormat]()},
		{name: "ProgramOption", kind: token.TYPE, root: reflect.TypeFor[ferret.ProgramOption](), native: reflect.TypeFor[engine.ProgramOption]()},
		{name: "Range", kind: token.TYPE, root: reflect.TypeFor[ferret.Range](), native: reflect.TypeFor[engine.Range]()},
		{name: "Session", kind: token.TYPE, root: reflect.TypeFor[ferret.Session](), native: reflect.TypeFor[engine.Session]()},
		{name: "SessionCloseHook", kind: token.TYPE, root: reflect.TypeFor[ferret.SessionCloseHook](), native: reflect.TypeFor[engine.SessionCloseHook]()},
		{name: "SessionOption", kind: token.TYPE, root: reflect.TypeFor[ferret.SessionOption](), native: reflect.TypeFor[engine.SessionOption]()},
		{name: "Source", kind: token.TYPE, root: reflect.TypeFor[ferret.Source](), native: reflect.TypeFor[engine.Source]()},
		{name: "Span", kind: token.TYPE, root: reflect.TypeFor[ferret.Span](), native: reflect.TypeFor[engine.Span]()},
		{name: "Value", kind: token.TYPE, root: reflect.TypeFor[ferret.Value](), native: reflect.TypeFor[engine.Value]()},

		{name: "DebugBreakpointBindExact", kind: token.CONST, root: ferret.DebugBreakpointBindExact, native: engine.DebugBreakpointBindExact},
		{name: "DebugBreakpointBindNextExecutableInFunction", kind: token.CONST, root: ferret.DebugBreakpointBindNextExecutableInFunction, native: engine.DebugBreakpointBindNextExecutableInFunction},
		{name: "DebugBreakpointBindNextExecutableInSource", kind: token.CONST, root: ferret.DebugBreakpointBindNextExecutableInSource, native: engine.DebugBreakpointBindNextExecutableInSource},
		{name: "DebugReasonBreakpoint", kind: token.CONST, root: ferret.DebugReasonBreakpoint, native: engine.DebugReasonBreakpoint},
		{name: "DebugReasonCompleted", kind: token.CONST, root: ferret.DebugReasonCompleted, native: engine.DebugReasonCompleted},
		{name: "DebugReasonEntry", kind: token.CONST, root: ferret.DebugReasonEntry, native: engine.DebugReasonEntry},
		{name: "DebugReasonPause", kind: token.CONST, root: ferret.DebugReasonPause, native: engine.DebugReasonPause},
		{name: "DebugReasonRuntimeError", kind: token.CONST, root: ferret.DebugReasonRuntimeError, native: engine.DebugReasonRuntimeError},
		{name: "DebugReasonStep", kind: token.CONST, root: ferret.DebugReasonStep, native: engine.DebugReasonStep},
		{name: "DebugReasonTerminated", kind: token.CONST, root: ferret.DebugReasonTerminated, native: engine.DebugReasonTerminated},
		{name: "LogDebug", kind: token.CONST, root: ferret.LogDebug, native: engine.LogDebug},
		{name: "LogDisabled", kind: token.CONST, root: ferret.LogDisabled, native: engine.LogDisabled},
		{name: "LogError", kind: token.CONST, root: ferret.LogError, native: engine.LogError},
		{name: "LogFatal", kind: token.CONST, root: ferret.LogFatal, native: engine.LogFatal},
		{name: "LogInfo", kind: token.CONST, root: ferret.LogInfo, native: engine.LogInfo},
		{name: "LogNone", kind: token.CONST, root: ferret.LogNone, native: engine.LogNone},
		{name: "LogPanic", kind: token.CONST, root: ferret.LogPanic, native: engine.LogPanic},
		{name: "LogTrace", kind: token.CONST, root: ferret.LogTrace, native: engine.LogTrace},
		{name: "LogWarn", kind: token.CONST, root: ferret.LogWarn, native: engine.LogWarn},
		{name: "OptimizationBasic", kind: token.CONST, root: ferret.OptimizationBasic, native: engine.OptimizationBasic},
		{name: "OptimizationFull", kind: token.CONST, root: ferret.OptimizationFull, native: engine.OptimizationFull},
		{name: "OptimizationNone", kind: token.CONST, root: ferret.OptimizationNone, native: engine.OptimizationNone},
		{name: "ProgramFormatJSON", kind: token.CONST, root: ferret.ProgramFormatJSON, native: engine.ProgramFormatJSON},
		{name: "ProgramFormatMsgPack", kind: token.CONST, root: ferret.ProgramFormatMsgPack, native: engine.ProgramFormatMsgPack},

		{name: "FormatError", kind: token.FUNC, root: ferret.FormatError, native: diagnostics.Format},
		{name: "MarshalProgram", kind: token.FUNC, root: ferret.MarshalProgram, native: engine.MarshalProgram},
		{name: "MustParseLogLevel", kind: token.FUNC, root: ferret.MustParseLogLevel, native: engine.MustParseLogLevel},
		{name: "New", kind: token.FUNC, root: ferret.New, native: engine.New},
		{name: "NewAnonymousSource", kind: token.FUNC, root: ferret.NewAnonymousSource, native: engine.NewAnonymousSource},
		{name: "NewSource", kind: token.FUNC, root: ferret.NewSource, native: engine.NewSource},
		{name: "ParseLogLevel", kind: token.FUNC, root: ferret.ParseLogLevel, native: engine.ParseLogLevel},
		{name: "UnmarshalProgram", kind: token.FUNC, root: ferret.UnmarshalProgram, native: engine.UnmarshalProgram},
		{name: "WithAfterCompileHook", kind: token.FUNC, root: ferret.WithAfterCompileHook, native: engine.WithAfterCompileHook},
		{name: "WithAfterRunHook", kind: token.FUNC, root: ferret.WithAfterRunHook, native: engine.WithAfterRunHook},
		{name: "WithBeforeCompileHook", kind: token.FUNC, root: ferret.WithBeforeCompileHook, native: engine.WithBeforeCompileHook},
		{name: "WithBeforeRunHook", kind: token.FUNC, root: ferret.WithBeforeRunHook, native: engine.WithBeforeRunHook},
		{name: "WithDebugFormat", kind: token.FUNC, root: ferret.WithDebugFormat, native: engine.WithDebugFormat},
		{name: "WithEncodingCodec", kind: token.FUNC, root: ferret.WithEncodingCodec, native: engine.WithEncodingCodec},
		{name: "WithEncodingRegistry", kind: token.FUNC, root: ferret.WithEncodingRegistry, native: engine.WithEncodingRegistry},
		{name: "WithEngineCloseHook", kind: token.FUNC, root: ferret.WithEngineCloseHook, native: engine.WithEngineCloseHook},
		{name: "WithEngineInitHook", kind: token.FUNC, root: ferret.WithEngineInitHook, native: engine.WithEngineInitHook},
		{name: "WithEnvironmentOptions", kind: token.FUNC, root: ferret.WithEnvironmentOptions, native: engine.WithEnvironmentOptions},
		{name: "WithFSReadOnly", kind: token.FUNC, root: ferret.WithFSReadOnly, native: engine.WithFSReadOnly},
		{name: "WithFSRoot", kind: token.FUNC, root: ferret.WithFSRoot, native: engine.WithFSRoot},
		{name: "WithFunctions", kind: token.FUNC, root: ferret.WithFunctions, native: engine.WithFunctions},
		{name: "WithFunctionsRegistrar", kind: token.FUNC, root: ferret.WithFunctionsRegistrar, native: engine.WithFunctionsRegistrar},
		{name: "WithLog", kind: token.FUNC, root: ferret.WithLog, native: engine.WithLog},
		{name: "WithLogFields", kind: token.FUNC, root: ferret.WithLogFields, native: engine.WithLogFields},
		{name: "WithLogLevel", kind: token.FUNC, root: ferret.WithLogLevel, native: engine.WithLogLevel},
		{name: "WithMaxActiveSessions", kind: token.FUNC, root: ferret.WithMaxActiveSessions, native: engine.WithMaxActiveSessions},
		{name: "WithMaxIdleVMsPerPlan", kind: token.FUNC, root: ferret.WithMaxIdleVMsPerPlan, native: engine.WithMaxIdleVMsPerPlan},
		{name: "WithMaxVMsPerPlan", kind: token.FUNC, root: ferret.WithMaxVMsPerPlan, native: engine.WithMaxVMsPerPlan},
		{name: "WithModules", kind: token.FUNC, root: ferret.WithModules, native: engine.WithModules},
		{name: "WithNamespace", kind: token.FUNC, root: ferret.WithNamespace, native: engine.WithNamespace},
		{name: "WithNetwork", kind: token.FUNC, root: ferret.WithNetwork, native: engine.WithNetwork},
		{name: "WithNetworkOptions", kind: token.FUNC, root: ferret.WithNetworkOptions, native: engine.WithNetworkOptions},
		{name: "WithOptimizationLevel", kind: token.FUNC, root: ferret.WithOptimizationLevel, native: engine.WithOptimizationLevel},
		{name: "WithOutputContentType", kind: token.FUNC, root: ferret.WithOutputContentType, native: engine.WithOutputContentType},
		{name: "WithParam", kind: token.FUNC, root: ferret.WithParam, native: engine.WithParam},
		{name: "WithParams", kind: token.FUNC, root: ferret.WithParams, native: engine.WithParams},
		{name: "WithPlanCloseHook", kind: token.FUNC, root: ferret.WithPlanCloseHook, native: engine.WithPlanCloseHook},
		{name: "WithPlanOptimizationLevel", kind: token.FUNC, root: ferret.WithPlanOptimizationLevel, native: engine.WithPlanOptimizationLevel},
		{name: "WithProgramFormat", kind: token.FUNC, root: ferret.WithProgramFormat, native: engine.WithProgramFormat},
		{name: "WithProgramLoader", kind: token.FUNC, root: ferret.WithProgramLoader, native: engine.WithProgramLoader},
		{name: "WithRuntimeParam", kind: token.FUNC, root: ferret.WithRuntimeParam, native: engine.WithRuntimeParam},
		{name: "WithRuntimeParams", kind: token.FUNC, root: ferret.WithRuntimeParams, native: engine.WithRuntimeParams},
		{name: "WithSessionCloseHook", kind: token.FUNC, root: ferret.WithSessionCloseHook, native: engine.WithSessionCloseHook},
		{name: "WithSessionFSRoot", kind: token.FUNC, root: ferret.WithSessionFSRoot, native: engine.WithSessionFSRoot},
		{name: "WithSessionLog", kind: token.FUNC, root: ferret.WithSessionLog, native: engine.WithSessionLog},
		{name: "WithSessionLogFields", kind: token.FUNC, root: ferret.WithSessionLogFields, native: engine.WithSessionLogFields},
		{name: "WithSessionLogLevel", kind: token.FUNC, root: ferret.WithSessionLogLevel, native: engine.WithSessionLogLevel},
		{name: "WithSessionParam", kind: token.FUNC, root: ferret.WithSessionParam, native: engine.WithSessionParam},
		{name: "WithSessionParams", kind: token.FUNC, root: ferret.WithSessionParams, native: engine.WithSessionParams},
		{name: "WithSessionRuntimeParam", kind: token.FUNC, root: ferret.WithSessionRuntimeParam, native: engine.WithSessionRuntimeParam},
		{name: "WithSessionRuntimeParams", kind: token.FUNC, root: ferret.WithSessionRuntimeParams, native: engine.WithSessionRuntimeParams},
		{name: "WithStdlib", kind: token.FUNC, root: ferret.WithStdlib, native: engine.WithStdlib},
		{name: "WithoutStdlib", kind: token.FUNC, root: ferret.WithoutStdlib, native: engine.WithoutStdlib},
	}
}
