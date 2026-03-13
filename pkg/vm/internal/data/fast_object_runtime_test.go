package data_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func TestFastObjectGetMissingReturnsNoneWithoutError(t *testing.T) {
	obj := data.NewFastObject(nil, 0)

	val, err := obj.Get(context.Background(), runtime.NewString("missing"))
	if err != nil {
		t.Fatalf("expected missing fast-object key to return no error, got %v", err)
	}

	if val != runtime.None {
		t.Fatalf("expected runtime.None for missing fast-object key, got %v", val)
	}
}
