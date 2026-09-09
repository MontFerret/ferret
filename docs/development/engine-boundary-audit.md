# Engine Boundary Audit

Tasks 4 and 5 of PR #1024 audit the native engine extraction, curated root
façade, and engine-owned internal packages, then remove redundant engine API
relays. The root API, FQL semantics, error identity, and lifecycle ownership
remain unchanged. Native callers use lower-level owners for shared vocabulary.
Subsequent construction and session extractions add `internal/bootstrap` and
`internal/session.Execution` with the ownership boundaries described below;
the native and root export inventories remain unchanged.

## Native exports

`pkg/engine` retains 55 package-level exports: seven types, three constants,
and 45 functions. Each represents native engine configuration, behavior, or
semantic vocabulary. Task 5 removes 35 aliases, 21 constants, and seven forwarding
functions; the earlier diagnostics cleanup already removed `engine.FormatError`.
The root retains all 119 exports. Independent root and native inventories in
[`api_contract_test.go`](../../api_contract_test.go) prevent relay reintroduction
and accidental façade growth. Future native exports must not automatically
expand the root API.

### Native engine contracts: retained

* Construction and lifecycle: `Engine`, `Plan`, `Session`, `New`.
* Configuration types: `Option`, `PlanOption`, `SessionOption`.
* Optimization: `OptimizationLevel`, `OptimizationNone`, `OptimizationBasic`,
  `OptimizationFull`, `WithOptimizationLevel`, `WithPlanOptimizationLevel`.
* Capacity: `WithMaxActiveSessions`, `WithMaxIdleVMsPerPlan`, `WithMaxVMsPerPlan`.
* Module and function registration: `WithModules`, `WithStdlib`, `WithoutStdlib`,
  `WithFunctions`, `WithFunctionsRegistrar`, `WithNamespace`.
* Encoding: `WithEncodingCodec`, `WithEncodingRegistry`, `WithOutputContentType`.
* Filesystem and network configuration: `WithFSRoot`, `WithFSReadOnly`,
  `WithSessionFSRoot`, `WithNetwork`, `WithNetworkOptions`.
* Engine parameters: `WithParam`, `WithParams`, `WithRuntimeParam`,
  `WithRuntimeParams`.
* Session parameters: `WithSessionParam`, `WithSessionParams`,
  `WithSessionRuntimeParam`, `WithSessionRuntimeParams`.
* Logging configuration: `WithLog`, `WithLogLevel`, `WithLogFields`,
  `WithSessionLog`, `WithSessionLogLevel`, `WithSessionLogFields`.
* Hook registration: `WithEngineInitHook`, `WithEngineCloseHook`,
  `WithBeforeCompileHook`, `WithAfterCompileHook`, `WithPlanCloseHook`,
  `WithBeforeRunHook`, `WithAfterRunHook`, `WithSessionCloseHook`.
* Advanced integration: `WithProgramLoader`, `WithEnvironmentOptions`,
  `WithDebugFormat`. These configure native artifact loading, VM environments,
  and debugger output without moving those subsystem contracts into the engine.

The 13 methods declared on native public types are intentional API:

* `Engine`: `Compile`, `CompileDebug`, `Load`, `Run`, `Close`.
* `Plan`: `Params`, `NewSession`, `NewDebugSession`, `Marshal`, `Close`.
* `Session`: `Run`, `Close`.
* `OptimizationLevel`: `String`.

Engine, Plan, and Session fields remain private. `Option`, `PlanOption`, and
`SessionOption` remain aliases instantiated with engine-owned configuration
types. `OptimizationLevel` remains a distinct engine type with validated
conversion to compiler levels, rather than a compiler alias.

### Shared vocabulary and helpers: direct owner references

The names below remain at root and target their lower-level owners directly.
Their engine equivalents are removed. Root types remain aliases, root constants
retain their types and values, and convenience functions remain declared
functions rather than assignable function variables.

* Source: `Source`, `Position`, `Span`, `Range`, `Location`, `NewSource`,
  `NewAnonymousSource` forward to `pkg/source`.
* Runtime and output: `Value`, `Params` alias `pkg/runtime`; `Output` aliases
  `pkg/encoding`.
* Module contracts: `Module`, `EngineInitHook`, `EngineCloseHook`,
  `BeforeCompileHook`, `AfterCompileHook`, `PlanCloseHook`, `BeforeRunHook`,
  `AfterRunHook`, `SessionCloseHook` alias `pkg/module`.
* Logging: `LogLevel`, `LogTrace`, `LogDebug`, `LogInfo`, `LogWarn`, `LogError`,
  `LogFatal`, `LogPanic`, `LogNone`, `LogDisabled`, `ParseLogLevel`,
  `MustParseLogLevel` forward to `pkg/logging`.
* Debugger types: `DebugSession`, `DebugReason`, `DebugValue`, `DebugVariable`,
  `DebugFrame`, `DebugBreakpoint`, `DebugBreakpointID`, `DebugBreakpointOptions`,
  `DebugBreakpointBindingMode`, `DebugValueReference`, `DebugEvent`,
  `DebugStateError`, `DebugFormatOptions` alias `pkg/debugger`, with
  `DebugSession = debugger.Session`; `DebugLocation = source.Range` and
  `DebugSourceLocation = source.Location` target `pkg/source`.
* Debugger constants: `DebugReasonEntry`, `DebugReasonBreakpoint`,
  `DebugReasonStep`, `DebugReasonPause`, `DebugReasonRuntimeError`,
  `DebugReasonCompleted`, `DebugReasonTerminated`,
  `DebugBreakpointBindNextExecutableInSource`, `DebugBreakpointBindExact`,
  `DebugBreakpointBindNextExecutableInFunction` preserve debugger values.
* Program serialization: `ProgramFormat`, `ProgramOption`, `ProgramFormatJSON`,
  `ProgramFormatMsgPack`, `WithProgramFormat`, `MarshalProgram`,
  `UnmarshalProgram` forward to `pkg/bytecode/artifact`. The standalone marshal
  helpers are existing low-level compatibility APIs; ordinary embedding uses
  `Plan.Marshal` and `Engine.Load`.
* Diagnostics: the declared root `FormatError` function already forwards
  directly to `diagnostics.Format` and remains unchanged in Task 5.

The source constructors call `source.New` and `source.NewAnonymous`; logging
parsers call `logging.ParseLogLevel` and `logging.MustParseLogLevel`. Program
helpers call `artifact.WithFormat`, `artifact.Marshal`, and `artifact.Unmarshal`.
Native `Plan.Marshal` also calls `artifact.Marshal` directly, retaining its
closed-plan guard.

Native callers must replace removed `engine.*` names with the owners above;
there are no compatibility aliases. For example, compilation takes
`source.Source`, execution returns `*encoding.Output`, `Plan.NewDebugSession`
returns `*debugger.Session`, and serialization accepts `artifact.Option`.
Option constructors remain in `pkg/engine`, accepting `module` hook types,
`runtime` values/parameters, `logging.LogLevel`, and `debugger.FormatOptions`
directly. Root callers keep their existing names and assignability.

## Internal contracts and ownership

* Engine construction uses concrete bootstrap `Config`, `Result`, and `Build`.
  Config supplies applied host inputs, hooks, resources, modules, compiler
  optimization, and session capacity. Build constructs normal/debug compilers,
  registers modules, finalizes the host, snapshots and initializes hooks, and
  creates the limiter. Result transfers these dependencies and the same resource
  manager to Engine on success. Bootstrap does not import either embedding
  package or construct the Engine object.
* Host construction uses concrete `Config`, `NewBootstrap`, and `Bootstrap`.
  `Bootstrap.Host` and `Bootstrap.Hooks` implement the existing module contract;
  `Bootstrap.Build` returns the finalized host. The underlying `hostContext`
  service methods implement `module.HostContext`. Its package-local `Build`
  method is narrowed to `build` because no external contract requires it.
* `Host` intentionally exposes the functions, parameters, codecs, logger,
  filesystem, network, and filesystem policy needed by engine orchestration.
  Parameter and codec snapshots isolate bootstrap mutations; service references
  remain borrowed, with lifecycle metadata owned by the resource manager.
  Fields are shared by convention, not structurally immutable. Accessors would
  expose the same references without improving this ownership contract.
* `Hooks`, `NewHooks`, `EngineHooks`, `PlanHooks`, and `SessionHooks` intentionally
  combine module registration and native execution. `Engine`, `Plan`, and
  `Session` expose module registrars; `Clone` isolates callback slices before
  execution. The concrete groups retain their registration and `Run*` methods
  for lifecycle-specific consumers. Their exported group pointers stay intact;
  separate registrar/runner types would add plumbing without a current need.
* `resource.Name`, `FileSystem`, `Network`, `Manager`, and `NewManager` identify
  owned or borrowed service slots. `Own`, `Borrow`, and `Close` retain private
  lifecycle metadata, replacement handling, and cleanup callbacks. These
  callbacks capture the owned resource; no lookup or dependency injection is
  exposed. This API is already narrow and stays engine-owned.
* Session admission retains `Limiter`, `NewLimiter`, `Acquire`, and `Release`.
  `PermitRelease` and `NewPermitRelease` remain the ordinary session mechanism
  for once-only permit release and returning a borrowed VM to its plan pool.
  `Materialize` remains shared native encoding/resource adoption behavior.
* Session preparation uses concrete `Config`, `NewConfig`, `NewExecution`, and
  `NewDebugSession`. Config contains applied settings and parameter conversion;
  a private native wrapper preserves the `SessionOption` target and native option
  validation/application. A private acquisition owner shares logger, permit,
  filesystem, and environment preparation and retains rollback until ownership
  transfers to the completed execution or debug session. Constructors receive
  concrete host, hook, limiter, pool/program, and Plan-close dependencies.
* `Execution` privately owns the ordinary VM, environment, host-resource manager,
  permit release, logger, filesystem/network references, and encoding settings.
  `Run` invokes the VM with its context services; `MaterializeAndClose` returns
  output and separate encoding/cleanup errors; `Close` runs close hooks, closes
  resources, returns the permit and VM, and retains the once-only cleanup result.
  Native `Session` retains the execution pointer, run hooks, and closed flag;
  it owns run admission, hook pairing, and output/error decisions. No operational
  state accessors or preparation-result bag cross the component boundary.
* `DebugServices` now has private state. `DebugServicesConfig` supplies only
  borrowed hooks, logger, filesystem, network, codecs, and output content type.
  `NewDebugServices` separately accepts the concrete limiter and resource
  manager. It neither acquires resources nor ends construction rollback.
  `BeforeRun`, `AfterRun`, `ExtendContext`, `Materialize`, and `Close` implement
  the existing debugger service contract. Close runs hooks, closes owned host
  resources, then releases and clears its limiter reference. The debug path has
  no VM-return callback; retained VM cleanup belongs to the debugger.

Engine defaults, option validation/application, stdlib registration, loader
configuration, public optimization-level translation, and final Engine assembly
remain in `pkg/engine`. Configuration failures retain their original cleanup
paths; bootstrap owns subsequent construction rollback. Native shutdown and
bootstrap rollback separately orchestrate hooks and the resource manager,
preserving error aggregation without a cleanup callback dependency.
`engine.go` retains the public Engine API and pre-bootstrap optimization-level
failure cleanup. Private compilation and plan construction live in
`engine_compile.go`; `engine_close.go` owns native shutdown. Engine configuration
and defaults live in `engine_config.go`, with public options grouped by core
program/limit settings, runtime registration, I/O, and lifecycle hooks.

Per-plan compiler selection, plan construction, session admission, and the final
session cancellation check remain in `pkg/engine`. Plan-specific session
orchestration lives in `plan_session.go`; public session options retain their
native configuration bridge. Successful construction
transfers ownership before that check; cancellation closes the completed native
or debugger session, while parent closure does not revoke its ownership.
Partial construction rollback belongs to the internal acquisition owner and
runs without session close hooks. Retained debug execution stays debugger-owned;
debug services do not own an ordinary Execution or borrow from the Plan's pool.

## Consumers and dependency direction

| Consumer | Import and reason |
| --- | --- |
| Root façade implementation | `pkg/engine` for engine semantics; lower-level owners for shared vocabulary and helpers |
| Root API tests | root `ferret`; contract tests compare owner identities, direct targets, and separate export inventories |
| README embedding example | root `ferret`, modeling downstream usage |
| SDK registration example, authoring tests, external harness test | corrected to root `ferret`, exercising supported authoring/embedding usage |
| `pkg/sdk/sdktest` implementation | retains `pkg/engine`, owning native executions and exposing its engine for lower-level assertions; output uses `pkg/encoding` directly |
| Compatibility adapters | retain `pkg/engine`, adapting native execution |
| Test CLI | root `ferret`, exercising the embedding façade |
| Security suites | retain `pkg/engine`, exercising native implementation behavior |
| Native engine/component tests | test their owning package directly |

Production imports and test imports were classified separately. No lower-level
production package imports root `ferret`. Compiler, runtime, VM, bytecode, and
stdlib do not depend on the engine. The engine does not import the root façade;
engine internals import neither their parent engine nor the root. Internal
dependencies flow from session to host/resource, and from host to resource.
No generic helper causes an upward dependency, so no abstraction needs moving.

## Deliberately unchanged

* Raw-bytecode root marshal helpers remain convenience APIs; artifact-loader
  configuration, VM environment options, and debugger formatting options remain
  native engine configuration. Removing vocabulary relays does not remove or
  redesign these contracts.
* Revisit host snapshot fields and combined hook registration/execution only if
  a concrete mutation or ownership problem appears. Their present shape is an
  intentional internal contract, not unfinished accessor work.

Root tests preserve the full alias/signature/constant/export inventory and check
direct ownership using resolved import paths. Native tests and benchmarks use
owner-qualified names while preserving rollback, cancellation, filesystem
isolation, output, debugger closure, cleanup order, errors, and sibling permit
coverage. Only native import spellings change; execution semantics do not.
Existing website examples use the unchanged root API and require no update.
