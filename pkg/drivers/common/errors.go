package common

import (
	"io"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/rs/zerolog"
)

var (
	ErrReadOnly    = core.Error(core.ErrInvalidOperation, "read only")
	ErrInvalidPath = core.Error(core.ErrInvalidOperation, "invalid path")
)

func CloseAll(logger zerolog.Logger, closers []io.Closer, msg string) {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			logger.Error().Err(err).Msg(msg)
		}
	}
}
