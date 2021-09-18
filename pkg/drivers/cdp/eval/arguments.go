package eval

import (
	"github.com/mafredri/cdp/protocol/runtime"
	"github.com/rs/zerolog"
)

type FunctionArguments []runtime.CallArgument

func (args FunctionArguments) MarshalZerologArray(a *zerolog.Array) {
	for _, arg := range args {
		if arg.ObjectID != nil {
			a.Str(string(*arg.ObjectID))
		} else {
			a.RawJSON(arg.Value)
		}
	}
}
