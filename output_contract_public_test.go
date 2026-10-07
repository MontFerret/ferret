package ferret_test

import (
	"context"
	"testing"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2"
)

var (
	_ func(*ferret.Engine, context.Context, ferret.Source, ...ferret.SessionOption) (api.Output, error) = (*ferret.Engine).Run
	_ func(*ferret.Session, context.Context) (api.Output, error)                                        = (*ferret.Session).Run
	_ api.Content                                                                                       = ferret.Content{}
	_ ferret.Content                                                                                    = api.Content{}
	_ api.Metadata                                                                                      = ferret.Metadata{}
	_ ferret.Metadata                                                                                   = api.Metadata{}
	_ api.Consumer                                                                                      = ferret.Consumer(nil)
	_ ferret.Consumer                                                                                   = api.Consumer(nil)
)

func TestOutputSentinelIdentity(t *testing.T) {
	if ferret.ErrOutputInUse != api.ErrOutputInUse || ferret.ErrOutputClosed != api.ErrOutputClosed {
		t.Fatal("public output errors must retain canonical identity")
	}
}
