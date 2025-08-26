package ferret

import (
	"github.com/MontFerret/ferret/pkg/exec"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type Engine = exec.Engine
type Option = exec.Option

type Plan = exec.Plan
type PlanOption = exec.PlanOption

type Session = exec.Session
type SessionOption = exec.SessionOption
type Environment = vm.Environment

type Result = exec.Result

var New = exec.New
var NewEnvironment = vm.NewEnvironment
var ToValue = runtime.Parse

// Engine helpers
var IsScalar = exec.IsScalar
var ForEach = exec.ForEach
var Collect = exec.Collect
var One = exec.One
var JSONStream = exec.JSONStream

// Env options
var WithParams = vm.WithParams
var WithParam = vm.WithParam
var WithFunctions = vm.WithFunctions
var WithFunction = vm.WithFunction
var WithFunctionSetter = vm.WithFunctionSetter
var WithLog = vm.WithLog
var WithLogLevel = vm.WithLogLevel
var WithLogFields = vm.WithLogFields

// Runtime helpers
var MustParseLogLevel = runtime.MustParseLogLevel
