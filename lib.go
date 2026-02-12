package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Engine = engine.Engine
type Option = engine.Option

type Plan = engine.Plan

type Session = engine.Session
type SessionOption = engine.SessionOption
type Environment = vm.Environment

type Result = engine.Result

var New = engine.New
var NewEnvironment = vm.NewEnvironment
var ToValue = runtime.Parse

// Engine helpers
var IsScalar = engine.IsScalar
var ForEach = engine.ForEach
var Collect = engine.Collect
var One = engine.One
var JSONStream = engine.JSONStream

// Env options
var WithParams = vm.WithParams
var WithParam = vm.WithParam
var WithFunctions = vm.WithFunctions
var WithFunction = vm.WithFunction
var WithFunctionSetter = vm.WithFunctionsBuilder
var WithLog = vm.WithLog
var WithLogLevel = vm.WithLogLevel
var WithLogFields = vm.WithLogFields

// Runtime helpers
var MustParseLogLevel = runtime.MustParseLogLevel

// Diagnostic
var FormatError = diagnostics.Format

type Formattable = diagnostics.Formattable
type FormattableError = diagnostics.FormattableError
