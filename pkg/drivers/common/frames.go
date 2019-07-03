package common

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func CollectFrames(ctx context.Context, receiver *values.Array, doc drivers.HTMLDocument) error {
	receiver.Push(doc)

	children, err := doc.GetChildDocuments(ctx)

	if err != nil {
		return err
	}

	children.ForEach(func(value core.Value, idx int) bool {
		childDoc, ok := value.(drivers.HTMLDocument)

		if !ok {
			err = core.TypeError(value.Type(), drivers.HTMLDocumentType)

			return false
		}

		return CollectFrames(ctx, receiver, childDoc) == nil
	})

	return nil
}
