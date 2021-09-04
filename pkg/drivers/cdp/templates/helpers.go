package templates

import (
	"bytes"
	"github.com/MontFerret/ferret/pkg/drivers"
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func Param(input core.Value) string {
	switch value := input.(type) {
	case values.String:
		return ParamString(value)
	case values.Float:
		return ParamFloat(value)
	case values.Int:
		return ParamInt(value)
	default:
		if value != values.None {
			return value.String()
		}

		return "null"
	}
}

func ParamList(value []core.Value) string {
	var buf bytes.Buffer
	lastIndex := len(value) - 1

	for i, input := range value {
		switch v := input.(type) {
		case values.String:
			buf.WriteString(EscapeString(string(v)))
		default:
			buf.WriteString(v.String())
		}

		if i != lastIndex {
			buf.WriteString(",")
		}
	}

	return buf.String()
}

func ParamStringList(value []values.String) string {
	var buf bytes.Buffer
	lastIndex := len(value) - 1

	for i, input := range value {
		buf.WriteString(EscapeString(string(input)))

		if i != lastIndex {
			buf.WriteString(",")
		}
	}

	return buf.String()
}

func ParamString(value values.String) string {
	return EscapeString(string(value))
}

func ParamErr(err error) string {
	return EscapeString(err.Error())
}

func ParamFloat(value values.Float) string {
	return strconv.FormatFloat(float64(value), 'f', 6, 64)
}

func ParamInt(value values.Int) string {
	return strconv.Itoa(int(value))
}

func EscapeString(value string) string {
	return "`" + value + "`"
}

func flipWhen(when drivers.WaitEvent) drivers.WaitEvent {
	return drivers.WaitEvent((int(when) + 1) % 2)
}
