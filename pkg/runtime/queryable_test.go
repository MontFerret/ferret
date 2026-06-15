package runtime

import (
	"context"
	"encoding/json"
	"testing"
)

func TestDefaultQueryOneReturnsFirstOrNone(t *testing.T) {
	tests := []struct {
		list List
		want Value
		name string
	}{
		{
			name: "nil-list",
			list: nil,
			want: None,
		},
		{
			name: "empty",
			list: NewArray(0),
			want: None,
		},
		{
			name: "one",
			list: NewArrayWith(NewString("only")),
			want: NewString("only"),
		},
		{
			name: "many",
			list: NewArrayWith(NewString("first"), NewString("second")),
			want: NewString("first"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DefaultQueryOne(context.Background(), Query{}, func(context.Context, Query) (List, error) {
				return tt.list, nil
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("unexpected value: got %v want %v", got, tt.want)
			}
		})
	}
}

func TestQueryJSONFieldNames(t *testing.T) {
	data, err := json.Marshal(Query{
		Options:    None,
		Params:     NewObjectWith(map[string]Value{"value": NewInt(1)}),
		Kind:       NewString("capture"),
		Expression: NewString("test"),
	})
	if err != nil {
		t.Fatalf("failed to marshal query: %v", err)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatalf("failed to unmarshal query fields: %v", err)
	}

	for _, key := range []string{"options", "params", "kind", "expression"} {
		if _, ok := fields[key]; !ok {
			t.Fatalf("expected query JSON field %q in %s", key, data)
		}
	}
	if len(fields) != 4 {
		t.Fatalf("expected exactly four query JSON fields, got %d in %s", len(fields), data)
	}
	if _, ok := fields["payload"]; ok {
		t.Fatalf("did not expect legacy payload field in %s", data)
	}
}
