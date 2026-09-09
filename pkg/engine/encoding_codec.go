package engine

import (
	"github.com/MontFerret/ferret/v2/pkg/encoding"
)

type encodingCodecAlias struct {
	encoding.Codec
	contentType string
}

func (c encodingCodecAlias) ContentType() string {
	return c.contentType
}
