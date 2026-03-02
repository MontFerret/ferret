package bytecode

type CollectorType int

const (
	CollectorTypeCounter CollectorType = iota
	CollectorTypeKey
	CollectorTypeKeyCounter
	CollectorTypeKeyGroup
	CollectorTypeAggregate
	CollectorTypeAggregateGroup
)
