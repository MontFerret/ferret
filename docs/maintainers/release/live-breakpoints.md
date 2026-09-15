# Live breakpoint replacement release follow-up

Native Ferret and its Universal adapter implement live, atomic source-wide
replacement. This does not establish support in an already released daemon or
IDE integration. ferretd changes are deferred to a separate clean task to preserve
unrelated work in that repository.

## Upstream dependency alignment

The API contract change is required: `api v1.0.0-alpha.16` does not declare
`debugger.Session.ReplaceBreakpoints` or `BreakpointRequest`.

1. Validate the coordinated API and Ferret changes together using a temporary
   Go workspace containing both repositories and Ferret's two tool modules.
2. Release the API contract, then update Ferret's root and `tools/apiref` module
   files to that published version. Do not invent a release version or commit a
   local dependency override.
3. Validate the root and tool modules independently with `GOWORK=off`, rerun
   affected tests, and release Ferret.
4. Only then begin the separate clean ferretd change below.

Until step 2, the staged Ferret change requires the coordinated API checkout;
the retained published module pins cannot compile the new portable type alias.

The low-level `vm.DebugExecution.Resume` integration now accepts a synchronous
`func(int) bool` breakpoint predicate instead of a fixed PC map. Native adapters
and test implementations must update that parameter. Existing source-level
single-breakpoint methods remain available.

## Separate ferretd task

Pin the released API/Ferret versions supporting live replacement and update the
current Universal bridge or official `uapi` composition, as appropriate to that
clean checkout. Do not remove the running restriction against the old dependency.

- Replace `internal/debug/session_breakpoints.go`'s delete/add loop with one
  upstream `ReplaceBreakpoints` call. Remove its shadow breakpoint set and
  breakpoint-specific inspectability restriction. Keep inspection restrictions.
- Forward the request immediately, without automatic pause, deferred mutation,
  or an IDE pending queue. Keep native synchronization authoritative.
- Preserve session request-arrival order. The inspected DAP read/dispatch loop
  processes requests sequentially; do not introduce parallel replacements that
  allow older requests to finish after newer ones and overwrite their state.
- Use native session breakpoint IDs directly in DAP responses and hit events.
  Remove native-to-DAP remapping and cleanup that discards removed identities:
  a previously committed stop may still report an ID removed by replacement,
  including when its event is delivered before the replacement response.
- Preserve launched-source canonicalization, quiet unverified results for other
  files, request ordering, resolved coordinates, and verification messages. Add
  a useful message for valid locations with no executable binding.
- Cover launch/configurationDone/run followed by live add, remove, replace,
  clear, rapid updates, relocation, and stop/termination races. Include independent
  sessions and assert that the former running-state rejection is gone.
- Run debug-domain and DAP tests with race detection, the full suite, and the
  clean checkout's lint/build gates. Update DAP and development documentation.

Release the daemon after validation. A later Editorium integration updates its
daemon pin. This prerequisite includes no Editorium or JetBrains implementation.

## CLI compatibility

Preserve the existing stopped add/list/delete/step workflow and native incremental
breakpoint tests. The inspected CLI checkout pins Ferret alpha.53 and still uses
older debugger names such as `Step`, `Next`, and `Out`; validating it against a
new native release requires its separate API migration. Do not fold that
migration or a new interactive live-command UX into this prerequisite.
