package testing

import (
	"context"
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func equalityFailureMessage(ctx context.Context, descriptor assertion, args []runtime.Value, positive bool) string {
	actual := args[0]
	expected := args[1]

	var message string
	if positive {
		message = unequalValuesMessage(ctx, actual, expected)
	} else {
		message = unexpectedlyEqualValuesMessage(ctx, actual, expected)
	}

	if len(args) == descriptor.args.max {
		return args[descriptor.args.max-1].String() + "\n" + message
	}

	return message
}

func unequalValuesMessage(ctx context.Context, actual, expected runtime.Value) string {
	differences, ok := discoverStructuralDifferences(ctx, actual, expected)
	if !ok {
		return fmt.Sprintf(
			"values are not equal\nexpected: %s\nactual:   %s",
			formatValue(ctx, expected),
			formatValue(ctx, actual),
		)
	}

	var builder strings.Builder
	builder.WriteString("values are not equal")
	for _, item := range differences.items {
		builder.WriteByte('\n')
		builder.WriteString(item.path)
		builder.WriteString("\n  expected: ")
		if item.kind == differenceExpectedMissing {
			builder.WriteString("<missing>")
		} else {
			builder.WriteString(formatValue(ctx, item.expected))
		}

		builder.WriteString("\n  actual:   ")
		if item.kind == differenceActualMissing {
			builder.WriteString("<missing>")
		} else {
			builder.WriteString(formatValue(ctx, item.actual))
		}
	}

	if differences.limitReached {
		builder.WriteString("\n... additional differences may be omitted")
	}

	return builder.String()
}

func unexpectedlyEqualValuesMessage(ctx context.Context, first, second runtime.Value) string {
	firstValue := formatValue(ctx, first)
	secondValue := formatValue(ctx, second)
	if firstValue == secondValue {
		return "expected values to differ\nboth: " + firstValue
	}

	return fmt.Sprintf(
		"expected values to differ\nfirst:  %s\nsecond: %s",
		firstValue,
		secondValue,
	)
}
