package debugger

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func formatValue(value runtime.Value, access vm.DebugValueAccess, options FormatOptions) string {
	info, _ := access.DebugInfo(value)

	return formatValueWithInfo(value, info, access, options)
}

func formatValueWithInfo(value runtime.Value, info runtime.DebugInfo, access vm.DebugValueAccess, options FormatOptions) string {
	if info.Display != "" {
		return boundedText(info.Display, options.MaxBytes)
	}

	var b strings.Builder

	writeValue(&b, value, info, access, options, 0)
	out := b.String()

	if len(out) <= options.MaxBytes {
		return out
	}

	return boundedText(out, options.MaxBytes)
}

func writeValue(b *strings.Builder, value runtime.Value, info runtime.DebugInfo, access vm.DebugValueAccess, options FormatOptions, depth int) {
	if info.Display != "" {
		b.WriteString(boundedText(info.Display, options.MaxBytes))

		return
	}

	if value == nil || reflect.TypeOf(value) == reflect.TypeOf(runtime.None) {
		b.WriteString("NONE")
		return
	}

	switch value := value.(type) {
	case runtime.Boolean:
		if value {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case runtime.Int:
		b.WriteString(strconv.FormatInt(int64(value), 10))
	case runtime.Float:
		b.WriteString(strconv.FormatFloat(float64(value), 'g', -1, 64))
	case runtime.String:
		b.WriteString(strconv.Quote(boundedText(value.String(), options.MaxBytes)))
	case runtime.DateTime:
		b.WriteString(boundedText(value.String(), options.MaxBytes))
	case runtime.Binary:
		fmt.Fprintf(b, "Binary(%d)", len(value))
	default:
		inspection, ok := access.Inspect(value, options.MaxItems)

		if !ok {
			typeName := info.TypeName

			if typeName == "" {
				typeName = access.TypeName(value)
			}

			fmt.Fprintf(b, "HostValue(%s)", boundedText(typeName, options.MaxBytes))
			return
		}

		if depth >= options.MaxDepth || inspection.Length > options.MaxItems || !inspection.Complete {
			writeCollectionSummary(b, inspection)
			return
		}

		writeInspection(b, inspection, access, options, depth)
	}
}

func writeInspection(b *strings.Builder, inspection vm.DebugValueInspection, access vm.DebugValueAccess, options FormatOptions, depth int) {
	if inspection.Kind == vm.DebugValueArray {
		b.WriteByte('[')

		for i, item := range inspection.Items {
			if i > 0 {
				b.WriteString(", ")
			}

			info, _ := access.DebugInfo(item.Value)
			writeValue(b, item.Value, info, access, options, depth+1)
		}

		b.WriteByte(']')

		return
	}

	items := append([]vm.DebugValueItem(nil), inspection.Items...)
	sort.Slice(items, func(i, j int) bool { return items[i].Key < items[j].Key })
	b.WriteByte('{')

	for i, item := range items {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString(strconv.Quote(boundedText(item.Key, options.MaxBytes)))
		b.WriteString(": ")
		info, _ := access.DebugInfo(item.Value)
		writeValue(b, item.Value, info, access, options, depth+1)
	}

	b.WriteByte('}')
}

func writeCollectionSummary(b *strings.Builder, inspection vm.DebugValueInspection) {
	if inspection.Kind == vm.DebugValueArray {
		fmt.Fprintf(b, "Array(%d)", inspection.Length)
		return
	}

	fmt.Fprintf(b, "Object(%d)", inspection.Length)
}

func boundedText(value string, max int) string {
	if value == "" {
		return value
	}

	if max <= 0 {
		return "..."
	}

	if !utf8.ValidString(value) {
		value = strings.ToValidUTF8(value, "\uFFFD")
	}

	if len(value) <= max {
		return value
	}

	for max > 0 && !utf8.RuneStart(value[max]) {
		max--
	}

	return value[:max] + "..."
}
