# Ferret
<p align="center">
	<a href="https://goreportcard.com/report/github.com/MontFerret/ferret">
		<img alt="Go Report Status" src="https://goreportcard.com/badge/github.com/MontFerret/ferret">
	</a>
	<a href="https://github.com/MontFerret/ferret/actions">
		<img alt="Build Status" src="https://github.com/MontFerret/ferret/workflows/build/badge.svg">
	</a>
	<a href="https://codecov.io/gh/MontFerret/ferret">
		<img src="https://codecov.io/gh/MontFerret/ferret/branch/main/graph/badge.svg" />
	</a>
	<a href="https://mastodon.social/@montferret">
		<img alt="Mastodon Follow" src="https://img.shields.io/mastodon/follow/114576925880917699?domain=mastodon.social&style=flat&label=on%20Mastodon">
	</a>
	<a href="https://t.me/montferret_chat">
		<img alt="Telegram Group" src="https://raw.githubusercontent.com/Patrolavia/telegram-badge/master/chat.svg">
	</a>
	<a href="https://github.com/MontFerret/ferret/releases">
		<img alt="Ferret release" src="https://img.shields.io/github/release/MontFerret/ferret.svg">
	</a>
	<a href="https://opensource.org/licenses/Apache-2.0">
		<img alt="Apache-2.0 License" src="http://img.shields.io/badge/license-Apache-brightgreen.svg">
	</a>
</p>

![ferret](https://raw.githubusercontent.com/MontFerret/ferret/main/assets/intro.jpg)

<p align="center">
	<a href="https://ferretlang.org/try" style="margin: 0 15px">
		<span>Try it!</span>
	</a>
	<a href="https://ferretlang.org/docs/introduction" style="margin: 0 15px">
		<span>Docs</span>
	</a>
	<a href="https://github.com/MontFerret/cli" style="margin: 0 15px">
		<span>CLI</span>
	</a>
	<a href="https://github.com/MontFerret/lab" style="margin: 0 15px">
		<span>Test runner</span>
	</a>
	<a href="https://github.com/MontFerret/worker" style="margin: 0 15px">
		<span>Web worker</span>
	</a>
</p>

---

## Explore Ferret v2

Ferret v2 is currently in alpha. You can try the new syntax in the playground and read more about the design behind the new runtime and language capabilities:

- [Try Ferret v2 in the Playground](https://ferretlang.org/try/)
- [On the Road to Ferret v2](https://ferretlang.org/blog/ferret-v2-announcement/)
- [Inside Ferret v2: The New Execution Model](https://ferretlang.org/blog/ferret-v2-execution-model/)
- [Inside Ferret v2: New Language Capabilities](https://ferretlang.org/blog/ferret-v2-new-syntax/)

---

> **Notice:** This branch contains the upcoming **Ferret v2**. For the stable v1 release, please visit [Ferret v1](https://github.com/MontFerret/ferret/tree/v1).

---

## What is it?

Ferret is a declarative runtime for structured data extraction and automation.

It lets you query web pages, browser state, documents, APIs, and host-provided data sources with a dedicated query language, then return the result as structured data.

Instead of writing page-specific glue code for browser control, DOM traversal, waiting, extraction, and transformation, Ferret lets you describe the data you want and run that workflow from the CLI, a worker, or an embedded Go application.

### Features

- Declarative query language for structured data workflows
- Support for static pages, dynamic pages, and browser-driven extraction
- CLI and embeddable Go runtime
- Extensible module, function, and runtime capability system
- Structured results for testing, analytics, AI/ML, and automation pipelines
- Portable execution model with a focused VM

## Getting started

```bash
go get github.com/MontFerret/ferret/v2@latest
```

There are currently two ways to start with Ferret v2:
- Native v2 API - recommended for new projects
- `compat` module - recommended as a first migration step for existing v1 integrations

### New projects

Use the native v2 API built around the following flow:

```
Engine -> compile query -> create session -> run
```

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/MontFerret/ferret/v2/pkg/engine"
)

func main() {
	ctx := context.Background()

	eng, err := engine.New()
	if err != nil {
		log.Fatal(err)
	}
	defer eng.Close()

	plan, err := eng.Compile(`RETURN 1 + 1`)
	if err != nil {
		log.Fatal(err)
	}

	session, err := plan.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	result, err := session.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.Content)
}
```

### Migration from v1

Ferret v2 introduces a new architecture and public API, so embedding it directly is different from v1.

To make migration easier, v2 includes a compat module that provides a v1-style API. Its goal is to make upgrades incremental instead of forcing a full rewrite up front.

For many projects, the easiest migration path will be:
- switch imports from v1 to the compat package
- get the project compiling again
- migrate incrementally to the native v2 API over time

A small helper script for rewriting import paths is planned to simplify this process further.

The compatibility layer is intended as a migration aid, not the long-term preferred API. New projects should use the native v2 packages directly.

### Source-level debugger API

Ferret v2 exposes a local debugger core for embedding tools such as CLIs and
editors. Compile with debugger metadata, create a debug session, then drive it
with breakpoints and stepping commands:

```go
plan, err := eng.CompileDebug(ctx, source.New("query.fql", query))
session, err := plan.NewDebugSession(ctx, ferret.WithSessionParam("input", 21))
defer session.Close()

_, err = session.SetBreakpointAt(
    ferret.DebugSourceLocation{File: "query.fql", Line: 4},
    ferret.DebugBreakpointOptions{},
)
event, err := session.Start(ctx)
event, err = session.Continue(ctx)
locals, err := session.Locals()
value, err := session.Evaluate(ctx, "result")
```

The canonical debugger models and advanced composition API live in
`pkg/debugger`. The top-level package retains `DebugSession`, event/value
types, reasons, breakpoint locations/options, and formatting options as
compatibility aliases. `SetBreakpoint(file, line)` remains a convenience
helper that binds to the next executable location in the file.

`Step` enters calls, `Next` stays at the same or shallower call depth, and
`Out` stops in the caller. Runtime errors pause before state is released so
frames and locals remain inspectable. Evaluation is deliberately conservative
and side-effect-free. It supports literals, locals and parameters, supported
member reads, scalar arithmetic and comparison, boolean logic, and simple
conditional expressions. It does not execute arbitrary Ferret code: function,
host, and module calls; queries and full collection semantics; mutation;
wait/async/event behavior; dispatch; recovery expressions; and opaque
host-value reads are rejected.

Custom runtime values may implement `runtime.DebugInspectable` to provide
optional debugger type-name and display hints. These hints are presentation
only and must be cheap, deterministic, side-effect free, and must not consume
lazy values or perform external work.

The reusable Phase 1 DAP adapter lives in `pkg/dap`. For local protocol testing,
the repository's test CLI can run it over stdin/stdout:

```bash
./bin/ferret --dap
./bin/ferret --dap --dap-trace
./bin/ferret --dap --dap-trace --dap-log-file ./ferret-dap.log
```

This is a test harness, not the production DAP server. A complete editor-facing
server belongs in a downstream integration repository. The adapter speaks
standard DAP framing over stdin/stdout, supports one launch source file and one
thread, and exposes top-frame locals, parameters, safe expression evaluation,
and bounded array/object expansion. Stdout is reserved exclusively for DAP
messages; trace output goes to stderr unless `--dap-log-file` is set.

Phase 1 does not support attach or remote transports, multiple sources or
threads, source maps, frame-specific locals or evaluation, mutation,
conditional breakpoints, hit conditions, logpoints, or a bundled VS Code
extension. Tail calls retain call depth, so `Next` and `Out` cannot always
distinguish tail-call boundaries.

The interactive `ferret debug` command still belongs to the separate CLI
repository.

### Alpha status

Ferret v2 is currently in active development.

Alpha releases are intended for early adopters, experimentation, and feedback. Some APIs and language features may still change before the stable v2 release.
