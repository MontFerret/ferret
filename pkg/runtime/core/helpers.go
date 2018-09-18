package core

import "reflect"

func IsNil(input interface{}) bool {
	val := reflect.ValueOf(input)
	kind := val.Kind()

	switch kind {
	case reflect.Ptr,
		reflect.Array,
		reflect.Slice,
		reflect.Map,
		reflect.Struct,
		reflect.Func,
		reflect.Interface,
		reflect.Chan,
		reflect.UnsafePointer:
		return val.IsNil()
	case reflect.Invalid:
		return true
	default:
		return false
	}
}
