package optimization

// mapsEqual checks if two register maps are equal.
func mapsEqual(a, b map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}

	for k := range a {
		if !b[k] {
			return false
		}
	}

	return true
}
