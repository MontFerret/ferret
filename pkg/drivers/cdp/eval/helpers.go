package eval

import (
	runtime2 "github.com/MontFerret/ferret/pkg/runtime/core"
	"strings"

	"github.com/mafredri/cdp/protocol/runtime"

	"github.com/MontFerret/ferret/pkg/drivers"
)

func parseRuntimeException(details *runtime.ExceptionDetails) error {
	if details == nil || details.Exception == nil {
		return nil
	}

	desc := *details.Exception.Description

	if strings.Contains(desc, drivers.ErrNotFound.Error()) {
		return drivers.ErrNotFound
	}

	return runtime2.Error(
		runtime2.ErrUnexpected,
		desc,
	)
}
