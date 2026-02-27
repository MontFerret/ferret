package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type SessionOption = vm.EnvironmentOption

var (
	WithSessionParams = vm.WithParams
	WithSessionParam  = vm.WithParam
)
