# DateTime library contracts

The DateTime stdlib group registers the canonical `datetime::` namespace and
deprecated global migration functions. Registration metadata is published through
the Core API generator. Deprecation is documentation metadata, without compiler
or execution warnings.

## Canonical surface

`now`, `parse`, `format`, `year`, `month`, `day`, `hour`, `minute`, `second`,
`millisecond`, `day_of_week`, `day_of_year`, `quarter`, `days_in_month`,
`is_leap_year`, `add`, `subtract`, `same`, and `diff` live in `datetime::`.
The Go package exports the corresponding PascalCase functions and `RegisterLib`.
`Parse` accepts one or two arguments; all other functions have fixed arities.
Unit parsing and deprecated FQL adapters are private implementation details.

Parsing accepts a String and optional Go time layout, defaulting to RFC3339.
Formatting accepts a DateTime and Go layout. Accessors inspect the value's own
local calendar. Weekdays remain Sunday=0 through Saturday=6; day-of-year starts
at 1, quarters at 1, and milliseconds range from 0 to 999. Month lengths use Go
calendar arithmetic in UTC with the input's local year and month, independent
of timezone transitions.

## Comparison and arithmetic

`same(a, b, unit)` compares local calendar fields through the specified
precision: year, year/month, calendar date, date/hour, date/hour/minute,
date/hour/minute/second, or those fields plus the truncated millisecond.
Week is a separate case comparing ISO week-year and ISO week number. Offset and
location identity are not compared, and inputs are not converted to UTC. Thus
two occurrences of the same local hour during a DST overlap can compare equal,
while the same instant represented on different local dates need not.

`diff(a, b, unit)` returns a signed Float for `b - a`, including fractional
units. Only milliseconds, seconds, minutes, and hours are supported. It uses
runtime-owned checked DateTime subtraction, propagating `runtime.ErrRange` when
the interval cannot fit a Duration (roughly 292 years). Day, week, month, and
year differences are rejected with an argument error, even for identical inputs.
Calendar-aware differences are deferred; they are never duration approximations.

`add` and `subtract` retain integer amounts and existing Go time arithmetic:
subday units use elapsed duration arithmetic, while day/week/month/year use
`time.Time.AddDate` in the input's location. A calendar day can span 23 or 25
elapsed hours. Month-end values normalize rather than clamp: January 31, 2023
plus one month becomes March 3, and February 29, 2024 plus one year becomes
March 1, 2025.

All unit-taking functions accept descriptive singular/plural names,
case-insensitively. Existing one-letter unit spellings remain accepted for
migration compatibility but are not the canonical vocabulary.

## Migration compatibility

All previously registered globals remain deprecated. `now` maps to
`datetime::now`; there is no `date_now`. `date` maps to `datetime::parse`.
Other globals drop `date_`, with `date_dayofweek`, `date_dayofyear`, and
`date_leapyear` mapping to `day_of_week`, `day_of_year`, and `is_leap_year`.
Thin forwarding declarations carry legacy metadata without deprecating the
canonical implementations.

Two compatibility adapters retain different signatures:

- `date_compare(a, b, start[, end])` requires **every** selected component to
  match. The inclusive component order is year, month, week, day, hour, minute,
  second, millisecond; the default end is millisecond. Week includes ISO
  week-year. Reversed ranges fail. This fixes the old any-component bug and
  remains distinct from canonical precision equality.
- `date_diff(a, b, unit[, asFloat])` shares the canonical checked duration
  calculation and supported units. Omitted/false returns an Int truncated
  toward zero by integer division; true returns the canonical Float. Negative
  results, calendar-unit rejection, and range errors intentionally replace the
  previous absolute, approximate, and saturating behavior.

Obsolete Go `Date*` entry points, exported unit constants, `Unit`,
`UnitFromString`, `AddUnit`, and `IsDatesEqual` are removed. Go callers migrate
to the canonical functions; legacy compatibility exists only in FQL.

Runtime DateTime representation and native operators are unchanged. Timezone
APIs, Unix timestamp APIs, boundary operations, numeric datetime construction,
and calendar differences are outside this library refactor.
