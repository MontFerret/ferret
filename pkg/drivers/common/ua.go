package common

import (
	"github.com/corpix/uarand"
)

const RandomUserAgent = "*"

func GetUserAgent(val string) string {
	if val == "" {
		return val
	}

	if val != RandomUserAgent {
		return val
	}

	// TODO: Change the implementation
	return uarand.GetRandom()
}
