package eval

import "github.com/mafredri/cdp/protocol/runtime"

type RemoteValue interface {
	RemoteID() runtime.RemoteObjectID
}
