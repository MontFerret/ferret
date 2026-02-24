package runtime

var (
	// Actual types
	TypeNone = newBuiltinType("None", func(v Value) bool {
		return v == None || v == nil
	})
	TypeBoolean = newBuiltinType("Boolean", func(v Value) bool {
		_, ok := v.(Boolean)
		return ok
	})
	TypeInt = newBuiltinType("Int", func(v Value) bool {
		_, ok := v.(Int)
		return ok
	})
	TypeFloat = newBuiltinType("Float", func(v Value) bool {
		_, ok := v.(Float)
		return ok
	})
	TypeString = newBuiltinType("String", func(v Value) bool {
		_, ok := v.(String)
		return ok
	})
	TypeDateTime = newBuiltinType("DateTime", func(v Value) bool {
		_, ok := v.(DateTime)
		return ok
	})
	TypeArray = newBuiltinType("Array", func(v Value) bool {
		_, ok := v.(*Array)
		return ok
	})
	TypeObject = newBuiltinType("Object", func(v Value) bool {
		_, ok := v.(ObjectLike)
		return ok
	})
	TypeBinary = newBuiltinType("Binary", func(v Value) bool {
		_, ok := v.(Binary)
		return ok
	})

	// Interfaces
	TypeCollection = newBuiltinType("Collection", func(v Value) bool {
		_, ok := v.(Collection)
		return ok
	})
	TypeList = newBuiltinType("List", func(v Value) bool {
		_, ok := v.(List)
		return ok
	})
	TypeMap = newBuiltinType("Map", func(v Value) bool {
		_, ok := v.(Map)
		return ok
	})

	// Capabilities
	TypeIndexReadable = newBuiltinType("IndexReadable", func(v Value) bool {
		_, ok := v.(IndexReadable)
		return ok
	})
	TypeIndexRemovable = newBuiltinType("IndexRemovable", func(v Value) bool {
		_, ok := v.(IndexRemovable)
		return ok
	})
	TypeIndexWritable = newBuiltinType("IndexWritable", func(v Value) bool {
		_, ok := v.(IndexWritable)
		return ok
	})
	TypeKeyReadable = newBuiltinType("KeyReadable", func(v Value) bool {
		_, ok := v.(KeyReadable)
		return ok
	})
	TypeKeyWritable = newBuiltinType("KeyWritable", func(v Value) bool {
		_, ok := v.(KeyWritable)
		return ok
	})
	TypeKeyRemovable = newBuiltinType("KeyRemovable", func(v Value) bool {
		_, ok := v.(KeyRemovable)
		return ok
	})
	TypeValueRemovable = newBuiltinType("ValueRemovable", func(v Value) bool {
		_, ok := v.(ValueRemovable)
		return ok
	})
	TypeAppendable = newBuiltinType("Appendable", func(v Value) bool {
		_, ok := v.(Appendable)
		return ok
	})
	TypeContainable = newBuiltinType("Containable", func(v Value) bool {
		_, ok := v.(Containable)
		return ok
	})
	TypeIterable = newBuiltinType("Iterable", func(v Value) bool {
		_, ok := v.(Iterable)
		return ok
	})
	TypeIterator = newBuiltinType("Iterator", func(v Value) bool {
		_, ok := v.(Iterator)
		return ok
	})
	TypeMeasurable = newBuiltinType("Measurable", func(v Value) bool {
		_, ok := v.(Measurable)
		return ok
	})
	TypeComparable = newBuiltinType("Comparable", func(v Value) bool {
		_, ok := v.(Comparable)
		return ok
	})
	TypeCloneable = newBuiltinType("Cloneable", func(v Value) bool {
		_, ok := v.(Cloneable)
		return ok
	})
	TypeSortable = newBuiltinType("Sortable", func(v Value) bool {
		_, ok := v.(Sortable)
		return ok
	})
	TypeDispatchable = newBuiltinType("Dispatchable", func(v Value) bool {
		_, ok := v.(Dispatchable)
		return ok
	})
	TypeObservable = newBuiltinType("Observable", func(v Value) bool {
		_, ok := v.(Observable)
		return ok
	})
	TypeQueryable = newBuiltinType("Queryable", func(v Value) bool {
		_, ok := v.(Queryable)
		return ok
	})
)
