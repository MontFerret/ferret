package bytecode

// MaxCollectionPreallocation limits the preallocated capacity accepted for
// persisted OpLoadArray and OpLoadObject instructions. This prevents untrusted
// artifacts from forcing excessive upfront allocations during execution.
const MaxCollectionPreallocation = 1 << 20
