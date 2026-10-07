package engine

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/debugger"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestMaterializationPreservesCallerContext(t *testing.T) {
	type key struct{}
	const contentType = "application/x-context-json"
	for _, mode := range []string{"native", "debugger"} {
		t.Run(mode, func(t *testing.T) {
			seen := 0
			want := "invocation"
			codec := aliasCodec{base: encodingjson.Default, contentType: contentType,
				pre: func(ctx context.Context, _ runtime.Value) error {
					seen++
					if ctx.Value(key{}) != want || ctx.Err() != nil {
						t.Fatalf("materialization context: value=%v err=%v", ctx.Value(key{}), ctx.Err())
					}
					return nil
				}}
			eng := mustNewEngine(t, WithEncodingCodec(contentType, codec),
				WithBeforeRunHook(func(ctx context.Context) (context.Context, error) {
					return context.WithValue(ctx, key{}, "execution"), nil
				}))
			ctx := context.WithValue(t.Context(), key{}, "invocation")
			if mode == "native" {
				plan := mustCompilePlan(t, eng, "RETURN 1")
				session := mustNewSession(t, plan, WithOutputContentType(contentType))
				if _, err := collectSession(session, ctx); err != nil {
					t.Fatal(err)
				}
			} else {
				plan, err := eng.CompileDebug(ctx, source.NewAnonymous("RETURN 1"))
				if err != nil {
					t.Fatal(err)
				}
				defer plan.Close()
				session, err := plan.NewDebugSession(ctx, WithOutputContentType(contentType))
				if err != nil {
					t.Fatal(err)
				}
				defer session.Close()
				if _, err := session.Start(ctx); err != nil {
					t.Fatal(err)
				}
				want = "resume"
				event, err := session.Continue(context.WithValue(t.Context(), key{}, want))
				if err != nil || event.Reason != debugger.ReasonCompleted || string(event.Output.Data) != "1" {
					t.Fatalf("completion: %v %v", event, err)
				}
			}
			if seen != 1 {
				t.Fatalf("hook calls=%d, want 1", seen)
			}
		})
	}
}
