package core

type Comparable interface {
	Compare(other Value) int64
}
