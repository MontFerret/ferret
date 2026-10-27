# Codec I/O

`pkg/encoding` owns content-type registration and the runtime-value codec
contract. JSON and MessagePack implement context-aware synchronous I/O:

```go
Encode(ctx context.Context, dst io.Writer, value runtime.Value) error
Decode(ctx context.Context, src io.Reader) (runtime.Value, error)
```

The caller lends I/O only for the operation. Codecs never close it, retain it, or
use it after returning. Write buffers may be reused when Write returns. Success
completes codec-owned writes and flushing; callers flush, finalize, and close
their own wrappers. Failure may leave a written prefix or consumed input.
Non-nil I/O errors retain their identity. A short write without an error returns
`io.ErrShortWrite` and stops without retrying.

`EncodeBytes(ctx, enc, value)` uses a bytes.Buffer and returns its written prefix
with any encoding error. `DecodeBytes(ctx, dec, data)` uses bytes.NewReader.
Atomic publication requires buffering and checking encoding success before
publishing those bytes. There is no compatibility interface for the former
byte methods. SDK host-value codecs and bytecode serialization are independent.

## Hooks and contexts

Encoder hooks run recursively in depth-first order for visited runtime values:

```go
type PreEncoderHook func(context.Context, runtime.Value) error
type PostEncoderHook func(context.Context, runtime.Value, error) error
type PreDecoderHook func(context.Context) error
type PostDecoderHook func(context.Context, runtime.Value, error) error
```

Decoder hooks run once per operation. Their post-hook receives the materialized
value on decoding success and runtime.None on decoding failure. Same-kind hooks
run in registration order and stop at the first hook error. A pre-hook failure
skips the operation and its post-hooks; ancestor encoder post-hooks still see
descendant failures. Joined errors preserve operation and hook causes.

Contexts reach hooks and runtime capabilities. Built-ins reject nil contexts and
nil or typed-nil readers/writers with errors wrapping runtime.ErrInvalidArgument.
Argument validation precedes ownership discovery and cancellation, so invalid
I/O remains an invalid-argument failure even with a canceled context.

After validation, encoding binds the optional ownership callback and completes
stored-root discovery before checking cancellation directly. This check also
observes cancellation during discovery and returns before user hooks,
serialization, or writes begin. Decoding checks cancellation directly after
validation, before hooks or reads; it has no source-root discovery phase.

Once an operation starts, built-ins check cancellation at bounded I/O boundaries,
every 256 traversal values/tokens, and between large payload chunks. They avoid
ctx.Err calls on every item. Arbitrary blocking
I/O and context-free custom marshalers cannot be forcibly interrupted.
Context-free Go MarshalJSON protocols use background context only at that
unavoidable boundary; stdlib and HTTP adapters use their caller's context.

## Implementations and materialization

Built-in codec values share immutable hook configuration through a private
pointer. The zero value needs no configuration allocation. Configurers publish
independent snapshots, and each operation borrows its configuration once;
recursive encoding does not copy hook slices at every value. This also avoids
allocating a copy of the codec when passing a concrete value to `EncodeBytes`.

After validation, root discovery, and the initial cancellation check, exact
`bytes.Buffer` destinations can encode a native String, Duration, or non-nil Regexp root
without borrowing traversal state. This requires no encoder hooks or adoption
callback. JSON accepts unescaped ASCII text up to 4,089 bytes; MessagePack accepts
strings shorter than 32 bytes using its fixstr representation. Both retain the
traversal step and final cancellation check. Other values and destinations use
normal traversal. JSON retains already-formatted native root text when its
shortcut declines it and uses the normal pooled string writer, avoiding repeated
duration formatting. The collecting helper has no separate codec dispatch API.

JSON writes numbers, containers, escaped strings, and Base64 progressively.
Numeric scratch and bounded 4 KiB payload scratch avoid container materialization
for generic iterables. Custom marshalers and generic Go-value fallbacks may
buffer. Its standard JSON tokenizer handles fragmented readers and strictly
rejects extra roots, malformed separators, and truncated containers; integer
width is preserved.

Exact `*bytes.Buffer` destinations use direct checked writes and available
capacity for numeric formatting. Hook-free native JSON ranges collect at most
4 KiB of numeric output between checked buffer writes, while retaining traversal
cancellation checks. Arbitrary writers, including wrappers around a buffer,
retain immediate checked writes so hooks and iterators observe failures during
traversal. JSON literals use reusable scratch without per-value allocations.
String encoding copies unescaped ASCII runs within the same payload checkpoints;
escaping, Unicode output, and invalid-UTF-8 replacement retain their existing
behavior.

MessagePack connects the dependency encoder/decoder to borrowed I/O through a
short-write and cancellation adapter. Its pooled encoding state owns both the
dependency encoder and checked writer. Instances are reset on every operation
and release I/O, contexts, discovery callbacks, and custom encoder settings on
every exit. Output bytes belong to the collecting caller, not the pool.
An exact-buffer adapter is selected once per operation and shares the base
writer's cancellation counters and sticky error. Normal strings and map keys
still use the dependency encoder, including settings changed by custom encoding
earlier in that operation. Wrappers around buffers retain generic checked I/O.
Hook-free MessagePack ranges without a discovery callback encode native integers
directly, retaining sampled cancellation; the same rule applies to exact integer
map values after key encoding. Native string keys avoid an interface method call,
while host-defined keys retain their String method. Other values and configured
operations use the ordinary visitor.
Generic iterables still materialize to
determine the length required by a MessagePack container header. Decoding in both
formats returns a fully materialized runtime.Value and may read ahead. This is
not a constant-memory guarantee or a lazy decoding API.

Both encoders borrow the source value. Their internal iterator handle owns only
a distinct acquired iterator, traversed with runtime.ForEachIter. An iterator
that is the source itself, or shares its runtime.Resource identity, remains with
the source's owner. Comparable identities are checked safely; non-comparable
resource aliases require runtime.Resource for stable identity. The encoders do
not change runtime.ForEach's cleanup contract for other callers.

Distinct acquired iterators close once after traversal and child encoding,
including on failure and cancellation. MessagePack collects the items, encodes
the collected list after successful traversal, then closes the iterator. Joined
non-EOF traversal failures remain errors even when another cause is io.EOF.

Native materialization lends one private, synchronous discovery callback through
the operation context. Both built-ins cache it once per encoding operation in a
private operation context wrapper, without replacing contexts passed to hooks or
host capabilities. Codecs report values; the callback registers them with
vm.Result. Concrete JSON/MessagePack
codec types, including configured and pointer instances, use this callback alone.
Custom codecs and decorators retain the injected ownership pre-hook for
compatibility; selection uses Go type rather than content type.

The root and successful iterator yields, host List.At results, and host map
callback values and keys acquire ownership before cancellation checkpoints,
key/separator writes, and user hooks can abandon them. Discovery scans only
concrete runtime.Array and runtime.Object storage at these acquisition boundaries
and around hooks or container failures. MessagePack also scans its collected
buffer before hooks and when hooks or failure can leave changed stored children.
Ordinary recursive visits register directly, keeping successful hook-free
traversal linear. Hook and failure rescans may repeat stored graph walks.
The iterative walker completes
independently of cancellation, resets cycle protection for each discovery, and
reuses scratch. Repeated containers are rescanned to discover later mutations.
It never evaluates host capabilities or lazy collections, runs marshalers, or
reads further iterator values merely to discover resources. Discovery does not
change which values receive user encoding hooks.

Result identity deduplication prevents comparable resources and ResourceID
aliases from closing twice. Plain dynamically non-comparable closers have no
safe identity and need runtime.Resource to deduplicate aliases. Direct codec
calls without the callback retain caller-owned cleanup. Callback and traversal
scratch are borrowed only for materialization; pools release operation state
before reuse.

Iterator close errors join encoding errors before the parent's post-hooks. The
VM result remains responsible for closing adopted source
and child resources. Direct codec callers retain that responsibility themselves.

Native and UAPI output remain eager and byte-backed. Native materialization
collects into an owned buffer with the original invocation context; before-run
replacement contexts remain execution-only. Debugger materialization receives
the active execution or resume-command context. Successful zero-byte encoding
produces present empty Content; failed zero-byte encoding leaves Content absent.
Nonempty error prefixes have LengthKnown=false. Encoding and cleanup error
identities and result/session cleanup ordering remain
intact. Collected data survives codec scratch reuse and result/session closure.
Query execution, record iteration, and the output lifecycle are unchanged.

## Performance evidence

JSON/MessagePack encoding and decoding, borrowed I/O, iterator ownership, and
native materialization benchmarks are recorded in the linked
[codec streaming migration report](codec-streaming-migration.md). The report
separates collecting buffers, cancellation, hooks, callback registration, and
stored discovery costs from the durable contracts above.

The [encoder performance recovery report](codec-encoder-performance.md) compares
the pre-migration revision, merged migration, and subsequent encoder
optimizations, including remaining overhead and allocation controls.
