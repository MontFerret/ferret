# Go Code Style

This guide defines the shared conventions for Ferret's handwritten Go code.
Preserve the rules and their exceptions when organizing or changing code.

## Go type and file structure

These rules are mandatory unless the change explicitly requires otherwise.

* Prefer grouped `type ( ... )` declarations for package-level types.
* Types declared in the same file should normally be placed in a single grouped
  `type` declaration rather than written as independent `type` declarations.
* This applies equally to structs, interfaces, aliases, named primitive types,
  and method-bearing types.
* Do not split types into independent declarations merely because one or more of
  them have methods.
* Keep related types together when they belong to the same narrow responsibility
  and their proximity makes the implementation easier to understand.
* A file may contain multiple related behavioral types when they form one
  cohesive concern.
* Split types into separate files based on responsibility and ownership, not
  simply because multiple types have methods.
* When a file contains only one package-level type, a standalone declaration is
  acceptable; do not create an artificial group containing a single type.
* When adding a package-level type to a file that already contains type
  declarations, incorporate it into the existing type group when it belongs to
  the same concern.
* Avoid scattering a cohesive family of small types across multiple files.
* Do not create `helpers.go`, `utils.go`, or similarly generic files as dumping
  grounds. Organize files around predictable responsibilities.

Preferred:

```go
type (
	PassResult struct {
		Metadata map[string]any
		Modified bool
	}

	PassContext struct {
		Program *bytecode.Program
	}

	Pass interface {
		Run(*PassContext) (*PassResult, error)
	}
)
```

Avoid independent declarations when the types belong to the same concern:

```go
type PassResult struct {
	Metadata map[string]any
	Modified bool
}

type PassContext struct {
	Program *bytecode.Program
}

type Pass interface {
	Run(*PassContext) (*PassResult, error)
}
```

The grouped declaration expresses that these types form one related family.

## Function and method ownership

These rules are mandatory unless the change explicitly requires otherwise.

* Organize files around cohesive responsibilities rather than individual types.
* A file may contain multiple related types and their methods when they
  participate in the same narrow concern.
* Keep methods close to the types they belong to.
* A file containing methods must not also contain unrelated package-level
  functions.
* Package-level functions may coexist with methods in a type-centered file only
  when they are constructors for types owned by that file.
* Constructors include conventional `New...` functions and other explicit
  construction functions whose primary responsibility is creating or
  initializing one of the file's types.
* If package-level behavior is not a constructor and has no natural receiver,
  place it in a separate responsibility-focused file.
* If behavior conceptually belongs to a type's state, invariants, lifecycle, or
  resources, implement it as a method rather than a package-level function.
* Do not keep a regular function beside methods merely because that function is
  used only by those methods.
* Split a file when it contains distinct responsibilities, not merely because it
  contains multiple behavioral or method-bearing types.
* Do not split cohesive behavior across files merely to enforce one type or one
  method-bearing type per file.

Preferred:

```go
type (
	Registry struct {
		entries map[string]*Entry
	}

	Entry struct {
		value Value
	}
)

func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string]*Entry),
	}
}

func (r *Registry) Add(name string, value Value) {
	// ...
}

func (r *Registry) Get(name string) (*Entry, bool) {
	// ...
}
```

Avoid mixing regular package functions with methods:

```go
func (r *Registry) Add(name string, value Value) {
	// ...
}

func normalizeName(name string) string {
	// ...
}

func (r *Registry) Get(name string) (*Entry, bool) {
	// ...
}
```

If `normalizeName` is intrinsic to `Registry`, make it a method. If it is
genuinely package-level behavior, move it to an appropriately named
responsibility-focused file.

## Comment conventions

* Do not comment every function or method mechanically.
* Exported functions and methods should normally have doc comments, especially
  in embedding-facing and extension-facing packages.
* Comment unexported code only when it carries non-obvious semantics,
  invariants, side effects, ownership, cleanup, recovery, or protocol behavior.
* Explain why the code exists, what must remain true, or how it must be used.
* Do not merely restate the symbol name or signature.
* Prefer semantic and lifecycle comments in compiler, VM, runtime, encoding,
  diagnostics, and debugger internals.
* Avoid comment wallpaper.

Preferred:

```go
// Close releases resources associated with the result.
// It is safe to call multiple times. Once closed, the result must not be reused.
func (r *Result) Close() error
```

Avoid comments such as `// Close closes the result.`

## Go control-flow spacing

These rules are mandatory for handwritten Go code. Blank lines separate logical
units and make control transfer visible.

### Declaration ordering

Within a file, place exported/public declarations before unexported/private declarations.

Do not interleave private helpers between public methods or functions. Keep the public API grouped toward the top of the file and implementation helpers below it.

Prefer:

`Public1 → Public2 → Public3 → Private1 → Private2`

over:

`Public1 → Public2 → Private1 → Public3 → Private2`

### Producer and immediate check

A declaration, assignment, call, assertion, lookup, or parse operation stays
adjacent to an `if` that immediately validates or consumes its result:

```go
res, err := doSome()
if err != nil {
	return err
}

named, ok := value.(*types.Named)
if !ok {
	return ErrUnsupported
}
```

Do not insert a blank line between the producer and its immediate check.

If this producer/check unit follows separate logic, add a blank line before it:

```go
prepareState()

value := lookup(name)
if value == nil {
	return ErrNotFound
}
```

No leading blank line is needed when the producer starts the enclosing block.

### Independent control flow

Separate independent `if` blocks with a blank line:

```go
if foo != nil {
	useFoo(foo)
}

if bar != nil {
	useBar(bar)
}
```

After a completed control-flow block, add a blank line before a separate
statement or logical unit.

### Return and break

`return` and `break` begin a new logical group when another statement precedes
them in the same block:

```go
result := buildResult()

return result
```

```go
if is(curr, "WAITFOR") {
	found = true

	break
}
```

No blank line is required when `return` is already the first statement in its
block. Do not add an artificial leading blank line at the start of a function or
block.

## Local type declarations

A function-local type is appropriate when it is small, passive, method-free,
used only by that function, and makes the local algorithm easier to understand.

Prefer a package-level unexported type when it represents a meaningful domain or
algorithm concept, spans a substantial part of a complex function, benefits
from being visible outside control flow, may gain methods, or is likely to be
reused by nearby helpers.

Do not promote a tiny throwaway struct merely for consistency, and do not hide a
meaningful concept inside a function merely to avoid a package-level type.
