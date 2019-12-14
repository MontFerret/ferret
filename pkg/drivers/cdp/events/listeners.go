package events

import "sync"

type ListenerCollection struct {
	mu     sync.Mutex
	values map[ID]map[ListenerID]Listener
}

func NewListenerCollection() *ListenerCollection {
	lc := new(ListenerCollection)
	lc.values = make(map[ID]map[ListenerID]Listener)

	return lc
}

func (lc *ListenerCollection) Size(eventID ID) int {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	bucket, exists := lc.values[eventID]

	if !exists {
		return 0
	}

	return len(bucket)
}

func (lc *ListenerCollection) Add(listener Listener) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	bucket, exists := lc.values[listener.EventID]

	if !exists {
		bucket = make(map[ListenerID]Listener)
		lc.values[listener.EventID] = bucket
	}

	bucket[listener.ID] = listener
}

func (lc *ListenerCollection) Remove(eventID ID, listenerID ListenerID) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	bucket, exists := lc.values[eventID]

	if !exists {
		return
	}

	delete(bucket, listenerID)
}

func (lc *ListenerCollection) Values(eventID ID) []Listener {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	bucket, exists := lc.values[eventID]

	if !exists {
		return []Listener{}
	}

	snapshot := make([]Listener, 0, len(bucket))

	for _, listener := range bucket {
		snapshot = append(snapshot, listener)
	}

	return snapshot
}
