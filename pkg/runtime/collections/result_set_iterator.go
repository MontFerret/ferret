package collections

type ResultSetIterator struct {
	values []ResultSet
	pos    int
}

func NewResultSetIterator(values []ResultSet) Iterator {
	return &ResultSetIterator{values, 0}
}

func (iterator *ResultSetIterator) HasNext() bool {
	return len(iterator.values) > iterator.pos
}

func (iterator *ResultSetIterator) Next() (ResultSet, error) {
	if len(iterator.values) > iterator.pos {
		idx := iterator.pos
		val := iterator.values[idx]
		iterator.pos++

		return val, nil
	}

	return nil, ErrExhausted
}
