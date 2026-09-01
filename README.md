# Ferret
<p align="center">
	<a href="https://github.com/MontFerret/ferret/actions">
		<img alt="Build Status" src="https://github.com/MontFerret/ferret/workflows/build/badge.svg">
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

Ferret is a declarative-first, expression-oriented embedded language and runtime for data automation.

FQL combines querying, transformation, synchronization, and structured results with host-defined values and capabilities. Applications can embed the runtime and decide exactly which functions, modules, data, and external operations a program can use.

The language keeps a declarative core and adds domain-oriented orchestration plus constrained mutable state for automation that cannot be expressed as a pure data transformation. It is deliberately focused rather than a general-purpose scripting language.

### Features

- Purpose-built declarative language for querying, transforming, synchronizing, and automating structured data
- Embeddable Go runtime with reusable compiled plans and isolated execution sessions
- Capability-based host values for exposing application objects, resources, and external systems directly to FQL
- Unified query model for browsers, APIs, databases, documents, and custom data sources
- Extensible runtime through namespaced functions, modules, hooks, and custom value types
- Event-driven synchronization and dispatch for interacting with asynchronous and stateful resources
- Managed resource lifecycle for files, connections, cursors, streams, and other host resources
- Bytecode VM and portable programs for efficient repeated execution and precompiled artifacts

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

	"github.com/MontFerret/ferret/v2"
)

func main() {
	ctx := context.Background()

	eng, err := ferret.New()
	if err != nil {
		log.Fatal(err)
	}
	defer eng.Close()

	plan, err := eng.Compile(ctx, ferret.NewAnonymousSource(`return 1 + 1`))
	if err != nil {
		log.Fatal(err)
	}
	defer plan.Close()

	session, err := plan.NewSession(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	output, err := session.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output.Content))
}
```

### Migration from v1

Ferret v2 introduces a new architecture and public API, so existing Go applications should migrate in two stages: first to the v2 compatibility API, then incrementally to the native v2 API.

Run the [`ferret migrate`](https://github.com/MontFerret/cli#migrating-embedded-ferret-applications) command from anywhere inside the application's Go module:

```bash
ferret migrate --dry-run # List the files that would change
ferret migrate --print   # Print a unified diff without changing files
ferret migrate           # Apply the migration
```

The command rewrites the documented v1 imports to their v2 compatibility packages, updates `go.mod` and `go.sum` as required by those rewrites, and formats changed Go files. Generated, vendored, and nested-module files are left untouched. Unsupported v1 imports are reported as manual follow-up; if the project vendors dependencies, run `go mod vendor` after applying the migration.

This is only the mechanical compatibility stage. The command does not convert application logic to the native v2 API or migrate drivers and other unsupported v1 packages. Running it again on an already-migrated compatibility project is a no-op and does not upgrade the Ferret v2 dependency merely because the CLI is newer.

After applying the migration, address any reported manual follow-up, build and test the application, and then migrate to the native v2 API over time. The compatibility layer is a migration aid, not the long-term preferred API; new projects should use the native v2 packages directly.

### Alpha status

Ferret v2 is currently in active development.

Alpha releases are intended for early adopters, experimentation, and feedback. Some APIs and language features may still change before the stable v2 release.

## Maintainers

- [Versioned Ferret Core API Reference](docs/maintainers/core-api-reference.md)

## Support Ferret

Ferret is supported by organizations and community members who help fund its continued development. [View all supporters](https://ferretlang.org/sponsor/).
