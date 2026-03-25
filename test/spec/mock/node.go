package mock

import (
	"context"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type Node struct {
	kind string
}

func NewNode(kind string) *Node {
	return &Node{kind: kind}
}

func (n *Node) MarshalJSON() ([]byte, error) {
	return encodingjson.Default.Encode(runtime.NewString(n.kind))
}

func (n *Node) String() string {
	return n.kind
}

func (n *Node) Unwrap() interface{} {
	return n.kind
}

func (n *Node) Hash() uint64 {
	return runtime.NewString(n.kind).Hash()
}

func (n *Node) Copy() runtime.Value {
	return n
}

func (n *Node) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return runtime.NewArrayWith(n).Iterate(ctx)
}

func (n *Node) Query(_ context.Context, q runtime.Query) (runtime.List, error) {
	switch q.Kind.String() {
	case "css":
		switch q.Payload.String() {
		case ".product":
			return runtime.NewArrayWith(NewNode("product")), nil
		case ".title":
			return runtime.NewArrayWith(NewNode("title")), nil
		case ".price":
			return runtime.NewArrayWith(NewNode("price")), nil
		default:
			return runtime.NewArrayWith(NewNode("node")), nil
		}
	case "text":
		return runtime.NewArrayWith(NewText(n.kind)), nil
	default:
		return runtime.NewArray(0), nil
	}
}
