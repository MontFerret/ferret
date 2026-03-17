package mem

import "io"

// AliasTracker maintains per-closer reference counts for register slots.
// It answers "how many live registers point at this closer?" in O(1),
// replacing the O(registers) scan in hasLiveRegisterAliasExcept.
//
// Keys are pre-validated io.Closer values (comparable, non-nil).
// Callers must extract the closer via TrackedCloserOf before calling.
type AliasTracker struct {
	counts map[io.Closer]int
}

// Inc increments the alias count for closer.
func (t *AliasTracker) Inc(closer io.Closer) {
	if t.counts == nil {
		t.counts = make(map[io.Closer]int)
	}

	t.counts[closer]++
}

// Dec decrements the alias count for closer and returns the new count.
// If the count reaches zero the entry is removed. Returns 0 if the
// closer was not tracked.
func (t *AliasTracker) Dec(closer io.Closer) int {
	if t.counts == nil {
		return 0
	}

	n, ok := t.counts[closer]
	if !ok {
		return 0
	}

	n--
	if n <= 0 {
		delete(t.counts, closer)

		if len(t.counts) == 0 {
			t.counts = nil
		}

		return 0
	}

	t.counts[closer] = n

	return n
}

// Count returns the current alias count for closer.
func (t *AliasTracker) Count(closer io.Closer) int {
	if t.counts == nil {
		return 0
	}

	return t.counts[closer]
}

// Reset clears all tracked aliases.
func (t *AliasTracker) Reset() {
	t.counts = nil
}
