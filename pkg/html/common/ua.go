package common

import (
	"github.com/MontFerret/ferret/pkg/runtime/env"
	"github.com/malisit/kolpa"
)

var k = kolpa.C()

func GetUserAgent(val string) string {
	if val != env.RandomUserAgent {
		return val
	}

	return k.UserAgent()
}
