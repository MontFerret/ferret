package internal

type resultUse uint8

const (
	resultRequired resultUse = iota
	resultDiscarded
)
