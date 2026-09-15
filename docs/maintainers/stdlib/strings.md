# String, encoding, and crypto contracts

`stdlib.Strings` registers global text operations. `stdlib.Encoding` and
`stdlib.Crypto` register the `encoding::` and `crypto::` namespaces independently.
`stdlib.Path` registers the complete `path::` namespace. `Full()` and `Safe()`
include all four groups. The encoding function package delegates JSON
serialization to `pkg/encoding/json`; it does not own a second JSON codec.

## Late-alpha API migration

These changes intentionally remove legacy names and overloads without aliases.

| Former call | Current call |
| --- | --- |
| `concat_separator(separator, values)` | `join(values, separator)` |
| `substitute(text, search, replacement[, limit])` | `replace(text, search, replacement[, limit])` |
| `substitute(text, search)` | `replace(text, search, "")` |
| `contains(text, search, true)` | `find_first(text, search)` for an index |
| `regex_match(text, pattern)` | `regex_find(text, pattern)` or `regex_find_all(text, pattern)` |
| Regex Boolean case options | Inline pattern flags such as `(?i)` |
| `json_parse`, `json_stringify` | `encoding::json_parse`, `encoding::json_stringify` |
| `encode_uri_component`, `decode_uri_component` | `encoding::query_escape`, `encoding::query_unescape` |
| `to_base64`, `from_base64` | `encoding::base64_encode`, `encoding::base64_decode` |
| `escape_html`, `unescape_html` | `encoding::html_escape`, `encoding::html_unescape` |
| `md5`, `sha1`, `sha512`, `random_token` | The same names under `crypto::` |
| `base`, `clean`, `dir`, `ext`, `is_abs`, `separate`, `match` | The same names under `path::` |
| Global path `join(parts...)` | `path::join(parts...)` |

All eight path functions move into `path::` without global aliases. Their path
semantics and arities are unchanged; global `join(values, separator)` is string
joining. For example, `path::base(path::join("a", "b.txt"))` returns `"b.txt"`.

`starts_with`, `ends_with`, `repeat`, and `crypto::sha256` are new. Identifiers
allow underscores after digits, so names such as `base64_encode` are callable
in FQL; the first character must still be a letter.

## Text and bounds

Text, patterns, separators, and cutsets require String values. Supplied bounds,
lengths, and limits require Int values; supplying a wrong type never selects a
default. `like` requires a Boolean when its optional case mode is supplied.

`left`, `right`, `substring`, `find_first`, and `find_last` count Unicode runes,
not bytes or grapheme clusters. Left/right counts must be non-negative and clamp
to the rune count. Substring offsets outside the string, including negative
offsets, return an empty string. Its optional length is non-negative and clamps
without overflowing. Find bounds clamp to `[0, rune count]`; the end is exclusive,
reversed bounds return `-1`, and empty searches return the effective start or end.

All three trim functions default to Unicode whitespace. Explicit cutsets are
literal sets of runes; an empty cutset is a no-op.

`split` and `regex_split` accept an optional non-negative maximum result count.
Zero produces `[]`; a positive limit preserves the unsplit remainder. `replace`
accepts an optional non-negative replacement count, with zero leaving text
unchanged. Omission means unlimited; explicit `-1` is invalid. Replacement is
literal, including dollar signs. An empty search matches at rune boundaries.
`repeat` requires a non-negative count that fits in the host integer range and
limits its result to 64 MiB (67,108,864 UTF-8 bytes), including when the count is
one. It checks the byte limit before allocation without overflowing. Oversized
results return an argument error on the count argument, which `ON ERROR` can
catch. Zero count or empty text produces an empty string for any valid count.

`join` takes one List of Strings and a String separator. It preserves order and
empty elements, inserts one separator between adjacent elements, and propagates
iteration errors and cancellation. It neither flattens lists nor skips None.
`concat` retains Any-to-string conversion, single-List iteration, and None
contributing no text. These are explicit exceptions to strict text inputs.

## Formatting and patterns

`fmt` accepts either automatic `{}` placeholders or zero-based decimal `{n}`
placeholders. Modes cannot be mixed. Indexed placeholders may repeat, but every
supplied value must be used. `{{` and `}}` escape literal braces. Malformed or
unbalanced placeholders, missing values, overflowing indexes, and unused values
are errors. Values use their runtime string representations.

`like` matches the whole string with glob syntax (`*`, `?`, character classes,
and alternatives). Percent and underscore are literal; SQL-LIKE translation is
removed. The optional Boolean still enables case-insensitive matching.

Regex functions use Go regular expressions and inline flags. `regex_test`
returns Boolean. `regex_find` returns None or `{match, groups, named}`;
`regex_find_all` returns an array of those objects. `groups` excludes the full
match, preserves capture declaration order, and uses empty strings for unmatched
groups. `named` maps capture names to strings. No captures yields `groups: []`
and `named: {}`.

`regex_find` and `regex_find_all` reject duplicate non-empty capture names with
an invalid-argument error on the pattern (argument 2), identifying the duplicated
name. Names are compared exactly, so `value` and `Value` are distinct. Validation
uses the compiled expression's subexpression names before matching; duplicates
are invalid even across alternatives, in unmatched optional groups, or when no
match exists. For example, `(?P<value>a)(?P<value>b)` is rejected. These errors are
catchable with `ON ERROR`. `regex_test`, `regex_replace`, and `regex_split` retain
Go's duplicate-name behavior because they do not return named capture objects.

All-match operations use Go's non-overlapping and zero-width match rules.
`regex_replace` expands `$1`, `${name}`, and `$$` using Go replacement semantics.
Invalid expressions return normal argument errors. Ignored/reserved arguments
and Boolean regex options are removed.

## Encoding and crypto

Query escaping uses query-form semantics: space becomes `+`, and a literal plus
becomes `%2B`. Base64 encoding accepts String UTF-8 bytes or Binary and emits
standard padded Base64; decoding returns Binary, preserving arbitrary bytes.
HTML helpers escape/unescape HTML character references. JSON stringify accepts
Any; JSON parse requires String and preserves the existing codec's behavior.

Hash functions accept String or Binary raw bytes and return lowercase hex.
`crypto::random_token` accepts lengths from 1 through 65536 inclusive. It uses
`crypto/rand.Reader`, a local buffer, and rejection sampling over ASCII lowercase
letters, uppercase letters, and digits. Entropy errors propagate; no shared
mutable generator state is used. The internal reader seam supports deterministic
sampling, failure, and cancellation tests without replacing the process reader.

## Validation and publication

Package tests cover boundaries, strict types, capture shapes, sampling, and
iteration failures. FQL integration runs at None, Basic, and Full optimization.
Registry tests enforce capability isolation, removed aliases, fixed arities,
the complete PATH namespace, and coexistence of string `join` with `path::join`.
Regex tests cover duplicate rejection before matching and unchanged capture
shapes. `BenchmarkRegexNamedCaptures` measures both structured regex APIs with
named and unnamed groups on matching and nonmatching input.

Structured Go comments feed the Core API reference and category catalog, which
now include Encoding and Crypto. The website's generated reference consumes a
published Core version: after releasing these changes, update its runtime v2
pin and run `mage build`. Do not publish an unreleased artifact under an existing
version or hand-edit generated reference pages.
