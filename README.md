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
		<img alt="Mastodon Follow" src="https://img.shields.io/mastodon/follow/114576925880917699?domain=mastodon.social&style=flat">
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
	<a href="https://www.montferret.dev/try" style="margin: 0 15px">
		<span>Try it!</span>
	</a>
	<a href="https://www.montferret.dev/docs/introduction" style="margin: 0 15px">
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

> **📢 Notice:** This branch contains the upcoming **Ferret v2**. For the stable v1 release, please visit [Ferret v1](https://github.com/MontFerret/ferret/tree/v1).

---

## What is it?
Ferret is a declarative system for working with web data - extracting it, querying it, and turning it into structured results for testing, analytics, machine learning, and other workflows. 
It allows users to focus on the data they need while abstracting away the complexity of browser automation, page interaction, and underlying execution details.

### Features

- Declarative query language
- Works with static and dynamic web pages
- Embeddable in Go applications
- Extensible runtime and function system
- Portable and fast

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

### Alpha status

Ferret v2 is currently in active development.

Alpha releases are intended for early adopters, experimentation, and feedback. Some APIs and language features may still change before the stable v2 release.
