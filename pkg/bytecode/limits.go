package bytecode

import "math/bits"

// MaxCollectionPreallocation limits the preallocated capacity accepted for
// persisted OpLoadArray and OpLoadObject instructions. This prevents untrusted
// artifacts from forcing excessive upfront allocations during execution.
const MaxCollectionPreallocation = 1 << 20

// MaxEncodedSortDirections limits the number of sort directions that can be
// encoded into OpDataSetMultiSorter. The directions are bit-packed into an int,
// so counts beyond the machine word size are not representable and must be
// rejected during persisted-program validation.
const MaxEncodedSortDirections = bits.UintSize
