package collections

type noopIterator struct{}

var NoopIterator = &noopIterator{}

func (iterator *noopIterator) HasNext() bool {
	return false
}

func (iterator *noopIterator) Next() (ResultSet, error) {
	return nil, ErrExhausted
}
