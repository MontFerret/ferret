package dom

import (
	"strings"

	"github.com/mafredri/cdp/protocol/page"
	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

var FrameIDType = core.NewType("ferret.drivers.cdp.dom.FrameID")

type FrameID page.FrameID

func NewFrameID(id page.FrameID) FrameID {
	return FrameID(id)
}

func (f FrameID) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(string(f), jettison.NoHTMLEscaping())
}

func (f FrameID) Type() core.Type {
	return FrameIDType
}

func (f FrameID) String() string {
	return string(f)
}

func (f FrameID) Compare(other core.Value) int64 {
	var s1 string
	var s2 string

	s1 = string(f)

	switch v := other.(type) {
	case FrameID:
		s2 = string(v)
	case *HTMLDocument:
		s2 = string(v.Frame().Frame.ID)
	case values.String:
		s2 = v.String()
	default:
		return -1
	}

	return int64(strings.Compare(s1, s2))
}

func (f FrameID) Unwrap() interface{} {
	return page.FrameID(f)
}

func (f FrameID) Hash() uint64 {
	return values.Hash(FrameIDType, []byte(f))
}

func (f FrameID) Copy() core.Value {
	return f
}
