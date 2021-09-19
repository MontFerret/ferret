package eval

import (
	"strings"

	"github.com/mafredri/cdp/protocol/runtime"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func parseRuntimeException(details *runtime.ExceptionDetails) error {
	if details == nil || details.Exception == nil {
		return nil
	}

	desc := *details.Exception.Description

	if strings.Contains(desc, drivers.ErrNotFound.Error()) {
		return drivers.ErrNotFound
	}

	return core.Error(
		core.ErrUnexpected,
		desc,
	)
}
