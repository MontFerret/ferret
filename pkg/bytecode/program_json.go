package bytecode

import (
	"encoding/base64"
	stdjson "encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	aggregateKeyMarkerType = runtime.Type("agg_key_marker")
)

type (
	programJSON struct {
		Source     *file.Source   `json:"source,omitempty"`
		Registers  int            `json:"registers"`
		Bytecode   []Instruction  `json:"bytecode"`
		Constants  []constantJSON `json:"constants,omitempty"`
		CatchTable []Catch        `json:"catch_table,omitempty"`
		Params     []string       `json:"params,omitempty"`
		Metadata   Metadata       `json:"metadata"`
	}

	constantJSON struct {
		Type  string             `json:"type"`
		Value stdjson.RawMessage `json:"value,omitempty"`
	}
)

func (p *Program) MarshalJSON() ([]byte, error) {
	if p == nil {
		return []byte("null"), nil
	}

	constants := make([]constantJSON, len(p.Constants))

	for i, value := range p.Constants {
		encoded, err := encodeConstant(value)

		if err != nil {
			return nil, fmt.Errorf("bytecode.Program: encode constant %d: %w", i, err)
		}

		constants[i] = encoded
	}

	payload := programJSON{
		Source:     p.Source,
		Registers:  p.Registers,
		Bytecode:   p.Bytecode,
		Constants:  constants,
		CatchTable: p.CatchTable,
		Params:     p.Params,
		Metadata:   p.Metadata,
	}

	return json.Marshal(payload)
}

func (p *Program) UnmarshalJSON(data []byte) error {
	if p == nil {
		return fmt.Errorf("bytecode.Program: UnmarshalJSON on nil pointer")
	}

	var decoded programJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	constants := make([]runtime.Value, len(decoded.Constants))
	for i, value := range decoded.Constants {
		decodedValue, err := decodeConstant(value)
		if err != nil {
			return fmt.Errorf("bytecode.Program: decode constant %d: %w", i, err)
		}

		constants[i] = decodedValue
	}

	p.Source = decoded.Source
	p.Registers = decoded.Registers
	p.Bytecode = decoded.Bytecode
	p.Constants = constants
	p.CatchTable = decoded.CatchTable
	p.Params = decoded.Params
	p.Metadata = decoded.Metadata

	return nil
}

func encodeConstant(value runtime.Value) (constantJSON, error) {
	if value == nil || value == runtime.None {
		return constantJSON{Type: encodeConstantType(runtime.TypeNone)}, nil
	}

	if value == AggregateKeyMarker {
		return constantJSON{Type: encodeConstantType(aggregateKeyMarkerType)}, nil
	}

	switch v := value.(type) {
	case runtime.Boolean:
		if v {
			return constantJSON{Type: encodeConstantType(runtime.TypeBoolean), Value: stdjson.RawMessage("true")}, nil
		}

		return constantJSON{Type: encodeConstantType(runtime.TypeBoolean), Value: stdjson.RawMessage("false")}, nil
	case runtime.Int:
		return constantJSON{Type: encodeConstantType(runtime.TypeInt), Value: stdjson.RawMessage(strconv.FormatInt(int64(v), 10))}, nil
	case runtime.Float:
		return constantJSON{Type: encodeConstantType(runtime.TypeFloat), Value: stdjson.RawMessage(strconv.FormatFloat(float64(v), 'g', -1, 64))}, nil
	case runtime.String:
		encoded, err := json.Marshal(string(v))

		if err != nil {
			return constantJSON{}, err
		}

		return constantJSON{Type: encodeConstantType(runtime.TypeString), Value: encoded}, nil
	case runtime.Binary:
		encoded, err := json.Marshal(base64.StdEncoding.EncodeToString([]byte(v)))

		if err != nil {
			return constantJSON{}, err
		}

		return constantJSON{Type: encodeConstantType(runtime.TypeBinary), Value: encoded}, nil
	case runtime.DateTime:
		encoded, err := json.Marshal(v.Time.Format(time.RFC3339Nano))

		if err != nil {
			return constantJSON{}, err
		}

		return constantJSON{Type: encodeConstantType(runtime.TypeDateTime), Value: encoded}, nil
	default:
		return constantJSON{}, fmt.Errorf("unsupported constant type %T", value)
	}
}

func encodeConstantType(typ runtime.Type) string {
	return strings.ToLower(typ.String())
}

func decodeConstant(frame constantJSON) (runtime.Value, error) {
	dt, err := decodeConstantType(frame.Type)

	if err != nil {
		return nil, err
	}

	switch dt {
	case runtime.TypeNone:
		return runtime.None, nil
	case aggregateKeyMarkerType:
		return AggregateKeyMarker, nil
	case runtime.TypeBoolean:
		raw := strings.TrimSpace(string(frame.Value))
		if raw == "" {
			return runtime.False, fmt.Errorf("empty bool value")
		}

		value, err := strconv.ParseBool(raw)
		if err != nil {
			return runtime.False, err
		}

		return runtime.NewBoolean(value), nil
	case runtime.TypeInt:
		raw := strings.TrimSpace(string(frame.Value))
		if raw == "" {
			return runtime.ZeroInt, fmt.Errorf("empty int value")
		}

		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return runtime.ZeroInt, err
		}

		return runtime.NewInt64(value), nil
	case runtime.TypeFloat:
		raw := strings.TrimSpace(string(frame.Value))
		if raw == "" {
			return runtime.ZeroFloat, fmt.Errorf("empty float value")
		}

		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return runtime.ZeroFloat, err
		}

		return runtime.NewFloat(value), nil
	case runtime.TypeString:
		var value string

		if err := json.Unmarshal(frame.Value, &value); err != nil {
			return runtime.EmptyString, err
		}

		return runtime.NewString(value), nil
	case runtime.TypeBinary:
		var encoded string

		if err := json.Unmarshal(frame.Value, &encoded); err != nil {
			return nil, err
		}

		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, err
		}

		return runtime.NewBinary(decoded), nil
	case runtime.TypeDateTime:
		var encoded string
		if err := json.Unmarshal(frame.Value, &encoded); err != nil {
			return nil, err
		}

		parsed, err := time.Parse(time.RFC3339Nano, encoded)
		if err != nil {
			return nil, err
		}

		return runtime.NewDateTime(parsed), nil
	default:
		return nil, fmt.Errorf("unsupported constant type %q", frame.Type)
	}
}

func decodeConstantType(typ string) (runtime.Type, error) {
	switch typ {
	case "none":
		return runtime.TypeNone, nil
	case "agg_key_marker":
		return aggregateKeyMarkerType, nil
	case "bool":
		return runtime.TypeBoolean, nil
	case "int":
		return runtime.TypeInt, nil
	case "float":
		return runtime.TypeFloat, nil
	case "string":
		return runtime.TypeString, nil
	case "binary":
		return runtime.TypeBinary, nil
	case "datetime":
		return runtime.TypeDateTime, nil
	default:
		return "", fmt.Errorf("unsupported constant type %q", typ)
	}
}
