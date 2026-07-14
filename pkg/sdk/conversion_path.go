package sdk

import "strconv"

type conversionPath struct {
	data []byte
}

func newConversionPath() conversionPath {
	data := make([]byte, 1, 64)
	data[0] = '$'

	return conversionPath{data: data}
}

func (path *conversionPath) PushField(name string) int {
	mark := len(path.data)
	path.data = append(path.data, '.')
	path.data = append(path.data, name...)

	return mark
}

func (path *conversionPath) PushIndex(index int) int {
	mark := len(path.data)
	path.data = append(path.data, '[')
	path.data = strconv.AppendInt(path.data, int64(index), 10)
	path.data = append(path.data, ']')

	return mark
}

func (path *conversionPath) PushKey(key string) int {
	mark := len(path.data)
	path.data = append(path.data, '[')
	path.data = strconv.AppendQuote(path.data, key)
	path.data = append(path.data, ']')

	return mark
}

func (path *conversionPath) Restore(mark int) {
	path.data = path.data[:mark]
}

func (path conversionPath) String() string {
	return string(path.data)
}
