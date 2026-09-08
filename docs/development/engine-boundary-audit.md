# Engine Boundary Audit

Task 4 of PR #1024 audits the native engine extraction, curated root façade,
and engine-owned internal packages. The audit preserves the root API, FQL
semantics, error identity, and existing lifecycle ownership. No additional
package extraction or abstraction is needed.

## Native exports

All 119 package-level exports in `pkg/engine` correspond to existing root
exports. The groups below account for each one. Native contracts remain public;
the remaining aliases and forwarding helpers stay exported because the root
façade uses them. No package-level native export exists solely to cross a Task 3
internal package boundary. The explicit façade inventory remains in
[`api_contract_test.go`](../../api_contract_test.go); future native exports must
not automatically expand the root API.

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

Engine, Plan, and Session fields remain private. Aliased types retain their
owning packages' method sets; the extraction introduces no extra methods on them.

### Façade vocabulary and helpers: retained

These names preserve existing façade signatures and type identity. Their
semantics remain in the indicated lower-level owner.

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
  `DebugStateError`, `DebugFormatOptions` alias `pkg/debugger`; `DebugLocation`
  and `DebugSourceLocation` alias `pkg/source`.
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
* Diagnostics: `FormatError` is a historical function variable initialized from
  `pkg/diagnostics`. The root declared function forwards to that variable; its
  mutability is retained for this task and deferred below.

## Internal contracts and ownership

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
* `DebugServices` now has private state. `DebugServicesConfig` supplies only
  borrowed hooks, logger, filesystem, network, codecs, and output content type.
  `NewDebugServices` separately accepts the concrete limiter and resource
  manager. It neither acquires resources nor ends construction rollback.
  `BeforeRun`, `AfterRun`, `ExtendContext`, `Materialize`, and `Close` implement
  the existing debugger service contract. Close runs hooks, closes owned host
  resources, then releases and clears its limiter reference. The debug path has
  no VM-return callback; retained VM cleanup belongs to the debugger.

Option validation, module iteration, compiler selection, plan construction,
session rollback/publication, VM-pool borrowing, and object closure remain in
`pkg/engine`: they coordinate native object invariants and ownership transfer.
The remaining private code does not justify another extraction.

## Consumers and dependency direction

| Consumer | Import and reason |
| --- | --- |
| Root façade implementation | `pkg/engine`, for aliases and forwarding |
| Root API tests | root `ferret`; contract tests also compare native identities |
| README embedding example | root `ferret`, modeling downstream usage |
| SDK registration example, authoring tests, external harness test | corrected to root `ferret`, exercising supported authoring/embedding usage |
| `pkg/sdk/sdktest` implementation | retains `pkg/engine`, owning native executions and exposing its engine for lower-level assertions |
| Compatibility adapters | retain `pkg/engine`, adapting native execution |
| Test CLI and security suites | retain `pkg/engine`, exercising native implementation behavior |
| Native engine/component tests | test their owning package directly |

Production imports and test imports were classified separately. No lower-level
production package imports root `ferret`. Compiler, runtime, VM, bytecode, and
stdlib do not depend on the engine. The engine does not import the root façade;
engine internals import neither their parent engine nor the root. Internal
dependencies flow from session to host/resource, and from host to resource.
No generic helper causes an upward dependency, so no abstraction needs moving.

## Deferred to Task 5

* Decide whether the native `FormatError` function variable should become a
  declared function. Reassignment currently changes root formatting through its
  forwarder; removing that behavior needs an explicit API decision.
* Reconsider the advanced façade contracts only as a separate API task:
  raw-bytecode marshal helpers, artifact-loader configuration, VM environment
  options, and debugger vocabulary. Their existing root declarations are
  protected here, not candidates for opportunistic removal.
* Revisit host snapshot fields and combined hook registration/execution only if
  a concrete mutation or ownership problem appears. Their present shape is an
  intentional internal contract, not unfinished accessor work.

Root alias/signature/constant/export tests remain unchanged. Existing native
tests cover rollback, cancellation, filesystem isolation, output, and debugger
closure; focused debugger-service tests cover cleanup order, errors, and sibling
permit preservation. This cleanup changes no documented public behavior and
requires no website documentation update.
