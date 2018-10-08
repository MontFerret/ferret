package common

import (
	"github.com/MontFerret/ferret/pkg/runtime/env"
	"github.com/corpix/uarand"
)

func GetUserAgent(val string) string {
	if val == "" {
		return val
	}

	if val != env.RandomUserAgent {
		return val
	}

	// TODO: Change the implementation
	return uarand.GetRandom()
}
