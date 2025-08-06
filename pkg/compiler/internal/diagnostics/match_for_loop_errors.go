package diagnostics

import (
	"github.com/MontFerret/ferret/pkg/file"
)

func matchForLoopErrors(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	prev := offending.Prev()

	if is(prev, "IN") {
		span := spanFromTokenSafe(prev.Token(), src)
		span.Start = span.End + 1
		span.End = span.Start + 1
		err.Message = "Expected expression after 'IN'"
		err.Hint = "Each FOR loop must iterate over a collection or range."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing value"),
		}

		return true
	}

	if is(prev, "FOR") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1
		err.Message = "Expected 'IN' after loop variable"
		err.Hint = "Use 'FOR x IN [iterable]' syntax."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing keyword"),
		}

		return true
	}

	if is(offending, "FOR") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1
		err.Message = "Expected loop variable before 'IN'"
		err.Hint = "FOR must declare a variable."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing variable"),
		}

		return true
	}

	if is(offending, "COLLECT") {
		msg := err.Message

		if has(msg, "COLLECT =") {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start = span.End
			span.End = span.Start + 1

			err.Message = "Expected variable before '=' in COLLECT"
			err.Hint = "COLLECT must group by a variable."
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing variable"),
			}

			return true
		}
	}

	if is(prev, "FILTER") {
		span := spanFromTokenSafe(prev.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		err.Message = "Expected condition after 'FILTER'"
		err.Hint = "FILTER requires a boolean expression."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing expression"),
		}

		return true
	}

	if is(prev, "LIMIT") {
		span := spanFromTokenSafe(prev.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		err.Message = "Expected number after 'LIMIT'"
		err.Hint = "LIMIT requires a numeric value."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing expression"),
		}

		return true
	}

	if isExtraneous(err.Message) {
		input := parseExtraneousInput(err.Message)

		if input != "','" {
			return false
		}

		var steps int

		// We walk back two tokens to find if the keyword is LIMIT.
		for ; steps < 2 && prev != nil; steps++ {
			prev = prev.Prev()
		}

		if is(prev, "LIMIT") {
			limitSpan := spanFromTokenSafe(prev.Token(), src)
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start = limitSpan.End + 1
			span.End += 4

			err.Message = "Too many arguments provided to LIMIT clause"
			err.Hint = "LIMIT accepts at most two arguments: offset and count."
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "unexpected third argument"),
			}

			return true
		}
	}

	return false
}
