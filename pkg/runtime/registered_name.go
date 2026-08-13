package runtime

type (
	registeredFunction[T FunctionConstraint] struct {
		function T
		name     string
	}

	registeredDisplayName struct {
		name  string
		count int
	}

	registeredDisplayNames struct {
		entries map[string]registeredDisplayName
	}
)

func newRegisteredDisplayNames() *registeredDisplayNames {
	return &registeredDisplayNames{
		entries: make(map[string]registeredDisplayName),
	}
}

func (n *registeredDisplayNames) Declare(key, name string) string {
	if declared, exists := n.entries[key]; exists {
		return declared.name
	}

	n.entries[key] = registeredDisplayName{name: name}

	return name
}

func (n *registeredDisplayNames) Add(key, name string) {
	entry, exists := n.entries[key]
	if !exists {
		entry.name = name
	}

	entry.count++
	n.entries[key] = entry
}

func (n *registeredDisplayNames) Remove(key string) {
	entry := n.entries[key]
	if entry.count <= 1 {
		delete(n.entries, key)

		return
	}

	entry.count--
	n.entries[key] = entry
}

func (n *registeredDisplayNames) Name(key string) (string, bool) {
	entry, exists := n.entries[key]

	return entry.name, exists
}

func (n *registeredDisplayNames) Seed(other *registeredDisplayNames) {
	if other == nil {
		return
	}

	for key, entry := range other.entries {
		n.Declare(key, entry.name)
	}
}
