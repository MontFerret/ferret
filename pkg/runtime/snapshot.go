package runtime

import "context"

type (
	// ListSnapshotter provides a native in-memory snapshot of a List.
	// A successful call must return a non-nil Array containing the list's current
	// values. Snapshot is a bulk read optimization; it does not imply transactional
	// or atomic point-in-time consistency for remote sources, and it does not carry
	// the VM's terminal result-materialization or ownership-transfer semantics.
	// Implementing ListSnapshotter does not by itself make a value a List.
	ListSnapshotter interface {
		Snapshot(context.Context) (*Array, error)
	}

	// MapSnapshotter provides a native in-memory snapshot of a Map.
	// A successful call must return a non-nil Object containing the map's current
	// properties. Snapshot is a bulk read optimization; it does not imply
	// transactional or atomic point-in-time consistency for remote sources, and it
	// does not carry the VM's terminal result-materialization or ownership-transfer
	// semantics.
	// Implementing MapSnapshotter does not by itself make a value a Map.
	MapSnapshotter interface {
		Snapshot(context.Context) (*Object, error)
	}
)
