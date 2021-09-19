package dom

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/runtime"
)

type NodeLoader struct {
	dom *Manager
}

func NewNodeLoader(dom *Manager) eval.ValueLoader {
	return &NodeLoader{dom}
}

func (n *NodeLoader) Load(ctx context.Context, frameID page.FrameID, _ eval.RemoteObjectType, _ eval.RemoteClassName, id runtime.RemoteObjectID) (core.Value, error) {
	return n.dom.ResolveElement(ctx, frameID, id)
}
