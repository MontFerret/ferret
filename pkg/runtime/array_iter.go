package runtime

type ArrayIterator struct {
	values *Array
	length int
	pos    int
}

func NewArrayIterator(values *Array) Iterator {
	return &ArrayIterator{values: values, length: len(values.data), pos: 0}
}

func (iter *ArrayIterator) HasNext(_ Context) (bool, error) {
	return iter.length > iter.pos, nil
}

func (iter *ArrayIterator) Next(_ Context) (value Value, key Value, err error) {
	if iter.pos >= iter.length {
		return None, None, Error(ErrInvalidOperation, "no more elements")
	}

	value = iter.values.data[iter.pos]
	key = NewInt(iter.pos)
	iter.pos++

	return
}
