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

type (
	programJSON struct {
		Source     *file.Source   `json:"source,omitempty"`
		Functions  Functions      `json:"functions,omitempty"`
		Bytecode   []Instruction  `json:"bytecode"`
		Constants  []constantJSON `json:"constants,omitempty"`
		CatchTable []Catch        `json:"catchTable,omitempty"`
		Params     []string       `json:"params,omitempty"`
		Metadata   Metadata       `json:"metadata"`
		ISAVersion int            `json:"isaversion,omitempty"`
		Registers  int            `json:"registers"`
	}

	constantJSON struct {
		Type  string             `json:"type"`
		Value stdjson.RawMessage `json:"value,omitempty"`
	}
)

func encodeConstant(value runtime.Value) (constantJSON, error) {
	if value == nil || value == runtime.None {
		return constantJSON{Type: encodeConstantType(runtime.TypeNone)}, nil
	}

	if value == AggregateKeyMarker {
		return constantJSON{Type: encodeConstantType(typeAggregateKeyMarker)}, nil
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
	return strings.ToLower(runtime.TypeName(typ))
}

func decodeConstant(frame constantJSON) (runtime.Value, error) {
	switch frame.Type {
	case encodeConstantType(runtime.TypeNone):
		return runtime.None, nil
	case encodeConstantType(typeAggregateKeyMarker):
		return AggregateKeyMarker, nil
	case encodeConstantType(runtime.TypeBoolean):
		raw := strings.TrimSpace(string(frame.Value))
		if raw == "" {
			return runtime.False, fmt.Errorf("empty bool value")
		}

		value, err := strconv.ParseBool(raw)
		if err != nil {
			return runtime.False, err
		}

		return runtime.NewBoolean(value), nil
	case encodeConstantType(runtime.TypeInt):
		raw := strings.TrimSpace(string(frame.Value))
		if raw == "" {
			return runtime.ZeroInt, fmt.Errorf("empty int value")
		}

		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return runtime.ZeroInt, err
		}

		return runtime.NewInt64(value), nil
	case encodeConstantType(runtime.TypeFloat):
		raw := strings.TrimSpace(string(frame.Value))
		if raw == "" {
			return runtime.ZeroFloat, fmt.Errorf("empty float value")
		}

		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return runtime.ZeroFloat, err
		}

		return runtime.NewFloat(value), nil
	case encodeConstantType(runtime.TypeString):
		var value string

		if err := json.Unmarshal(frame.Value, &value); err != nil {
			return runtime.EmptyString, err
		}

		return runtime.NewString(value), nil
	case encodeConstantType(runtime.TypeBinary):
		var encoded string

		if err := json.Unmarshal(frame.Value, &encoded); err != nil {
			return nil, err
		}

		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, err
		}

		return runtime.NewBinary(decoded), nil
	case encodeConstantType(runtime.TypeDateTime):
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
