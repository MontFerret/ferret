package common

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
)

func CollectFrames(ctx context.Context, receiver *internal.Array, doc drivers.HTMLDocument) error {
	receiver.Push(doc)

	children, err := doc.GetChildDocuments(ctx)

	if err != nil {
		return err
	}

	children.ForEach(func(value core.Value, idx int) bool {
		childDoc, ok := value.(drivers.HTMLDocument)

		if !ok {
			err = core.TypeError(value, drivers.HTMLDocumentType)

			return false
		}

		return CollectFrames(ctx, receiver, childDoc) == nil
	})

	return nil
}
