package runtime

import (
	"context"
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
