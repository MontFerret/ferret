package runtime

var (
	// Actual types
	TypeNone = NewType("None", func(v Value) bool {
		return v == None || v == nil
	})
	TypeBoolean = NewType("Boolean", func(v Value) bool {
		_, ok := v.(Boolean)
		return ok
	})
	TypeInt = NewType("Int", func(v Value) bool {
		_, ok := v.(Int)
		return ok
	})
	TypeFloat = NewType("Float", func(v Value) bool {
		_, ok := v.(Float)
		return ok
	})
	TypeString = NewType("String", func(v Value) bool {
		_, ok := v.(String)
		return ok
	})
	TypeDateTime = NewType("DateTime", func(v Value) bool {
		_, ok := v.(DateTime)
		return ok
	})
	TypeArray = NewType("Array", func(v Value) bool {
		_, ok := v.(*Array)
		return ok
	})
	TypeObject = NewType("Object", func(v Value) bool {
		_, ok := v.(ObjectLike)
		return ok
	})
	TypeBinary = NewType("Binary", func(v Value) bool {
		_, ok := v.(Binary)
		return ok
	})

	// Interfaces
	TypeCollection = NewType("Collection", func(v Value) bool {
		_, ok := v.(Collection)
		return ok
	})
	TypeList = NewType("List", func(v Value) bool {
		_, ok := v.(List)
		return ok
	})
	TypeMap = NewType("Map", func(v Value) bool {
		_, ok := v.(Map)
		return ok
	})

	// Capabilities
	TypeIndexReadable = NewType("IndexReadable", func(v Value) bool {
		_, ok := v.(IndexReadable)
		return ok
	})
	TypeIndexRemovable = NewType("IndexRemovable", func(v Value) bool {
		_, ok := v.(IndexRemovable)
		return ok
	})
	TypeIndexWritable = NewType("IndexWritable", func(v Value) bool {
		_, ok := v.(IndexWritable)
		return ok
	})
	TypeKeyReadable = NewType("KeyReadable", func(v Value) bool {
		_, ok := v.(KeyReadable)
		return ok
	})
	TypeKeyWritable = NewType("KeyWritable", func(v Value) bool {
		_, ok := v.(KeyWritable)
		return ok
	})
	TypeKeyRemovable = NewType("KeyRemovable", func(v Value) bool {
		_, ok := v.(KeyRemovable)
		return ok
	})
	TypeValueRemovable = NewType("ValueRemovable", func(v Value) bool {
		_, ok := v.(ValueRemovable)
		return ok
	})
	TypeAppendable = NewType("Appendable", func(v Value) bool {
		_, ok := v.(Appendable)
		return ok
	})
	TypeContainable = NewType("Containable", func(v Value) bool {
		_, ok := v.(Containable)
		return ok
	})
	TypeIterable = NewType("Iterable", func(v Value) bool {
		_, ok := v.(Iterable)
		return ok
	})
	TypeIterator = NewType("Iterator", func(v Value) bool {
		_, ok := v.(Iterator)
		return ok
	})
	TypeMeasurable = NewType("Measurable", func(v Value) bool {
		_, ok := v.(Measurable)
		return ok
	})
	TypeComparable = NewType("Comparable", func(v Value) bool {
		_, ok := v.(Comparable)
		return ok
	})
	TypeCloneable = NewType("Cloneable", func(v Value) bool {
		_, ok := v.(Cloneable)
		return ok
	})
	TypeSortable = NewType("Sortable", func(v Value) bool {
		_, ok := v.(Sortable)
		return ok
	})
	TypeDispatchable = NewType("Dispatchable", func(v Value) bool {
		_, ok := v.(Dispatchable)
		return ok
	})
	TypeObservable = NewType("Observable", func(v Value) bool {
		_, ok := v.(Observable)
		return ok
	})
	TypeQueryable = NewType("Queryable", func(v Value) bool {
		_, ok := v.(Queryable)
		return ok
	})
)
