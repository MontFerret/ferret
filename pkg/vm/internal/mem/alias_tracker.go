package mem

// AliasTracker maintains per-resource reference counts for register slots.
// It answers "how many live registers point at this resource?" in O(1),
// replacing the O(registers) scan in hasLiveRegisterAliasExcept.
//
// Keys are ResourceKey values derived via ResourceKeyOf. Callers must resolve
// the key before calling Inc/Dec/Count.
type AliasTracker struct {
	counts map[ResourceKey]int
}

// Inc increments the alias count for key.
func (t *AliasTracker) Inc(key ResourceKey) {
	if t.counts == nil {
		t.counts = make(map[ResourceKey]int)
	}

	t.counts[key]++
}

// Dec decrements the alias count for key and returns the new count.
// If the count reaches zero the entry is removed. Returns 0 if the
// key was not tracked.
func (t *AliasTracker) Dec(key ResourceKey) int {
	if t.counts == nil {
		return 0
	}

	n, ok := t.counts[key]
	if !ok {
		return 0
	}

	n--
	if n <= 0 {
		delete(t.counts, key)

		if len(t.counts) == 0 {
			t.counts = nil
		}

		return 0
	}

	t.counts[key] = n

	return n
}

// Count returns the current alias count for key.
func (t *AliasTracker) Count(key ResourceKey) int {
	if t.counts == nil {
		return 0
	}

	return t.counts[key]
}

// Delete removes all alias tracking for key.
func (t *AliasTracker) Delete(key ResourceKey) {
	if t.counts == nil {
		return
	}

	delete(t.counts, key)

	if len(t.counts) == 0 {
		t.counts = nil
	}
}

// Reset clears all tracked aliases.
func (t *AliasTracker) Reset() {
	t.counts = nil
}
