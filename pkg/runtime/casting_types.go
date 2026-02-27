package runtime

import "reflect"

var (
	typeBinary = reflect.TypeOf(Binary(nil))
	typeArray  = reflect.TypeOf((*Array)(nil))
	typeObject = reflect.TypeOf((*Object)(nil))

	typeObjectLike = reflect.TypeOf((*ObjectLike)(nil)).Elem()
	typeCollection = reflect.TypeOf((*Collection)(nil)).Elem()
	typeList       = reflect.TypeOf((*List)(nil)).Elem()
	typeMap        = reflect.TypeOf((*Map)(nil)).Elem()

	typeIndexReadable  = reflect.TypeOf((*IndexReadable)(nil)).Elem()
	typeIndexWritable  = reflect.TypeOf((*IndexWritable)(nil)).Elem()
	typeIndexRemovable = reflect.TypeOf((*IndexRemovable)(nil)).Elem()
	typeKeyReadable    = reflect.TypeOf((*KeyReadable)(nil)).Elem()
	typeKeyWritable    = reflect.TypeOf((*KeyWritable)(nil)).Elem()
	typeKeyRemovable   = reflect.TypeOf((*KeyRemovable)(nil)).Elem()
	typeValueRemovable = reflect.TypeOf((*ValueRemovable)(nil)).Elem()
	typeAppendable     = reflect.TypeOf((*Appendable)(nil)).Elem()
	typeContainable    = reflect.TypeOf((*Containable)(nil)).Elem()

	typeIterable     = reflect.TypeOf((*Iterable)(nil)).Elem()
	typeIterator     = reflect.TypeOf((*Iterator)(nil)).Elem()
	typeMeasurable   = reflect.TypeOf((*Measurable)(nil)).Elem()
	typeComparable   = reflect.TypeOf((*Comparable)(nil)).Elem()
	typeCloneable    = reflect.TypeOf((*Cloneable)(nil)).Elem()
	typeSortable     = reflect.TypeOf((*Sortable)(nil)).Elem()
	typeDispatchable = reflect.TypeOf((*Dispatchable)(nil)).Elem()
	typeObservable   = reflect.TypeOf((*Observable)(nil)).Elem()
	typeQueryable    = reflect.TypeOf((*Queryable)(nil)).Elem()
)

// expectedTypeOf resolves the expected runtime Type for T, even when T's zero value is nil.
func expectedTypeOf[T Value]() Type {
	var zero T
	t := reflect.TypeOf((*T)(nil)).Elem()

	switch t.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return expectedTypeFromReflect(t)
	default:
		return TypeOf(zero)
	}
}

func expectedTypeFromReflect(t reflect.Type) Type {
	if t == nil {
		return TypeNone
	}

	switch t {
	case typeBinary:
		return TypeBinary
	case typeArray:
		return TypeArray
	case typeObject:
		return TypeObject
	}

	if t.Kind() == reflect.Interface {
		switch {
		case t.AssignableTo(typeObjectLike):
			return TypeObject
		case t.AssignableTo(typeList):
			return TypeList
		case t.AssignableTo(typeMap):
			return TypeMap
		case t.AssignableTo(typeCollection):
			return TypeCollection
		case t.AssignableTo(typeIndexReadable):
			return TypeIndexReadable
		case t.AssignableTo(typeIndexWritable):
			return TypeIndexWritable
		case t.AssignableTo(typeIndexRemovable):
			return TypeIndexRemovable
		case t.AssignableTo(typeKeyReadable):
			return TypeKeyReadable
		case t.AssignableTo(typeKeyWritable):
			return TypeKeyWritable
		case t.AssignableTo(typeKeyRemovable):
			return TypeKeyRemovable
		case t.AssignableTo(typeValueRemovable):
			return TypeValueRemovable
		case t.AssignableTo(typeAppendable):
			return TypeAppendable
		case t.AssignableTo(typeContainable):
			return TypeContainable
		case t.AssignableTo(typeIterable):
			return TypeIterable
		case t.AssignableTo(typeIterator):
			return TypeIterator
		case t.AssignableTo(typeMeasurable):
			return TypeMeasurable
		case t.AssignableTo(typeComparable):
			return TypeComparable
		case t.AssignableTo(typeCloneable):
			return TypeCloneable
		case t.AssignableTo(typeSortable):
			return TypeSortable
		case t.AssignableTo(typeDispatchable):
			return TypeDispatchable
		case t.AssignableTo(typeObservable):
			return TypeObservable
		case t.AssignableTo(typeQueryable):
			return TypeQueryable
		}
	}

	return typeFromReflect(t)
}
