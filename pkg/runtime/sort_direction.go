package runtime

// SortDirection represents the sorting direction, either ascending or descending.
type SortDirection = Int

const (
	SortDirectionAsc  SortDirection = iota // Ascending sort direction
	SortDirectionDesc                      // Descending sort direction
)

func NewSortDirection(direction Int) SortDirection {
	if direction == 0 {
		return SortDirectionAsc
	}

	return SortDirectionDesc
}

// EncodeSortDirections encodes a slice of SortDirection values into a single integer by combining their bit representations.
func EncodeSortDirections(directions []SortDirection) int {
	result := 0

	for _, dir := range directions {
		result = (result << 1) | int(dir)
	}

	return result
}

// DecodeSortDirections decodes an integer into a slice of SortDirection values representing sorting directions.
// The number of decoded directions is determined by the count argument.
// Each bit of the encoded integer corresponds to a SortDirection value in the resulting slice.
func DecodeSortDirections(encoded int, count int) []SortDirection {
	directions := make([]SortDirection, count)

	for i := count - 1; i >= 0; i-- {
		directions[i] = SortDirection(encoded & 1)
		encoded >>= 1
	}

	return directions
}
