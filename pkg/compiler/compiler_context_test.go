package compiler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestCompileRejectsInvalidContexts(t *testing.T) {
	c, err := compiler.New()
	if err != nil {
		t.Fatal(err)
	}

	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	expired, cancelDeadline := context.WithDeadline(t.Context(), time.Unix(1, 0))
	defer cancelDeadline()
	for _, tc := range []struct {
		name string
		ctx  context.Context
		want error
	}{
		{"nil", nil, runtime.ErrInvalidArgument},
		{"canceled", canceled, context.Canceled},
		{"deadline", expired, context.DeadlineExceeded},
	} {
		t.Run(tc.name, func(t *testing.T) {
			program, err := c.Compile(tc.ctx, source.NewAnonymous("RETURN 1"))
			if program != nil || !errors.Is(err, tc.want) {
				t.Fatalf("program = %v, error = %v", program, err)
			}
		})
	}
}
