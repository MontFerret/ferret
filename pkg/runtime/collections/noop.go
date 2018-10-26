package collections

type noopIterator struct{}

var NoopIterator = &noopIterator{}

func (iterator *noopIterator) HasNext() bool {
	return false
}

func (iterator *noopIterator) Next() (DataSet, error) {
	return nil, ErrExhausted
}
