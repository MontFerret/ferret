package bytecode

import "math/bits"

// MaxCollectionPreallocation limits the preallocated capacity accepted for
// persisted OpLoadArray and OpLoadObject instructions. This prevents untrusted
// artifacts from forcing excessive upfront allocations during execution.
const MaxCollectionPreallocation = 1 << 20

// MaxEncodedSortDirections limits the number of sort directions that can be
// encoded into OpDataSetMultiSorter using the current signed-int bit-packing.
// The highest bit must remain clear so the encoded operand stays non-negative,
// which caps the supported direction count at one less than the machine word
// size during persisted-program validation.
const MaxEncodedSortDirections = bits.UintSize - 1
