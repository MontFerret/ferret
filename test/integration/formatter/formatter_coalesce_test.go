package formatter_test

import (
	"bytes"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/pkg/source"
	. "github.com/MontFerret/ferret/v2/test/spec/format"
)

func TestFormatterCoalesceInline(t *testing.T) {
	RunSpecs(t, []Spec{
		S(`let name=user.name??"Unknown" return name`, `let name = user.name ?? "Unknown"
return name`),
		S(`return (none??1)??2`, `return (none ?? 1) ?? 2`),
	})
}

func TestFormatterCoalesceChainContinuation(t *testing.T) {
	RunSpecsWith(t, mustNewFormatter(t, formatter.WithPrintWidth(50)), []Spec{
		S(`let name=user?.profile?.displayName??user?.nickname??user?.name??"Anonymous" return name`, `let name = user?.profile?.displayName
    ?? user?.nickname
    ?? user?.name
    ?? "Anonymous"
return name`),
		S(`return primaryValue??(secondaryValue??tertiaryValue)`, `return primaryValue
    ?? (secondaryValue
        ?? tertiaryValue)`),
	})
}

func TestFormatterCoalesceRoundTrip(t *testing.T) {
	format := mustNewFormatter(t, formatter.WithPrintWidth(50))
	inputs := []string{
		`let name=user?.profile?.displayName??user?.nickname??user?.name??"Anonymous" return name`,
		`return primaryValue??(secondaryValue??tertiaryValue)`,
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			var first bytes.Buffer
			if err := format.Format(&first, source.NewAnonymous(input)); err != nil {
				t.Fatalf("first format failed: %v", err)
			}

			var second bytes.Buffer
			if err := format.Format(&second, source.NewAnonymous(first.String())); err != nil {
				t.Fatalf("second format failed: %v", err)
			}

			if got, want := second.String(), first.String(); got != want {
				t.Fatalf("coalescing format was not stable:\nfirst:\n%s\nsecond:\n%s", want, got)
			}
		})
	}
}
