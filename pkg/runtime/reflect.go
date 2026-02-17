package runtime

import (
	"context"
	"fmt"
	"reflect"
)

const TagName = "ferret"

func ReadValueByKey(ctx context.Context, input any, key Value) (Value, error) {
	if input == nil {
		return None, nil
	}

	wrapper, ok := input.(Unwrappable)

	if ok {
		input = wrapper.Unwrap()
	}

	t := reflect.TypeOf(input)
	keyStr := key.String()

	switch t.Kind() {
	case reflect.Struct:
		v := reflect.ValueOf(input)

		var name string

		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)

			if f.Tag.Get(TagName) == keyStr {
				if !f.IsExported() {
					return None, fmt.Errorf("field %s is not exported", keyStr)
				}

				name = f.Name

				break
			}
		}

		if name == "" {
			return None, nil
		}

		field := v.FieldByName(name)

		if field.IsValid() {
			return Parse(field.Interface()), nil
		}

		return None, nil
	case reflect.Map:
		v := reflect.ValueOf(input)
		field := v.MapIndex(reflect.ValueOf(key))

		if field.IsValid() {
			return Parse(field.Interface()), nil
		}

		return None, nil
	case reflect.Ptr:
		v := reflect.ValueOf(input)

		if v.IsNil() {
			return None, nil
		}

		return ReadValueByKey(ctx, v.Elem().Interface(), key)
	default:
		break
	}

	return None, fmt.Errorf("cannot read field by key from type: %s", t.String())
}
