package runtime

type (
	Allocator interface {
		Object(size int) *Object
		Array(cap int) *Array
	}

	allocatorImpl struct{}
)

func NewAllocator() Allocator {
	return &allocatorImpl{}
}

func (a *allocatorImpl) Object(size int) *Object {
	return newObjectOf(size)
}

func (a *allocatorImpl) Array(cap int) *Array {
	return newArray(cap)
}
