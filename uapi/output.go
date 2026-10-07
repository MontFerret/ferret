package uapi

import (
	"context"

	"github.com/MontFerret/api"
)

// outputAdapter projects terminal diagnostics without owning native lifecycle.
type outputAdapter struct {
	native api.Output
}

var _ api.Output = (*outputAdapter)(nil)

func (o *outputAdapter) Metadata() api.Metadata {
	return o.native.Metadata()
}

func (o *outputAdapter) Consume(ctx context.Context, consumer api.Consumer) error {
	return wrapDiagnosticError(o.native.Consume(ctx, consumer))
}

func (o *outputAdapter) Collect(ctx context.Context) (*api.Content, error) {
	content, err := o.native.Collect(ctx)

	return content, wrapDiagnosticError(err)
}

func (o *outputAdapter) Close() error {
	return wrapDiagnosticError(o.native.Close())
}
