package clauses

type CollectProjection struct {
	selector *CollectSelector
}

func NewCollectProjection(selector *CollectSelector) *CollectProjection {
	return &CollectProjection{selector}
}
