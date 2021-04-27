package core

import (
	"context"
	"math"
	"math/rand"
	"reflect"
	"time"
)

func IsNil(input interface{}) bool {
	val := reflect.ValueOf(input)
	kind := val.Kind()

	switch kind {
	case reflect.Ptr,
		reflect.Array,
		reflect.Slice,
		reflect.Map,
		reflect.Func,
		reflect.Interface,
		reflect.Chan:
		return val.IsNil()
	case reflect.Struct,
		reflect.UnsafePointer:
		return false
	case reflect.Invalid:
		return true
	default:
		return false
	}
}

func NumberBoundaries(input float64) (max float64, min float64) {
	min = input / 2
	max = input * 2

	return
}

func NumberUpperBoundary(input float64) float64 {
	return input * 2
}

func NumberLowerBoundary(input float64) float64 {
	return input / 2
}

func RandomDefault() float64 {
	rand.Seed(time.Now().UnixNano())

	return rand.Float64()
}

func Random(max float64, min float64) float64 {
	r := RandomDefault()
	i := r * (max - min + 1)
	out := math.Floor(i) + min

	return out
}

func Random2(mid float64) float64 {
	max, min := NumberBoundaries(mid)

	return Random(max, min)
}

func ForEach(ctx context.Context, iter Iterator, predicate func(value Value, key Value) bool) error {
	for {
		value, key, err := iter.Next(ctx)

		if err != nil {
			if IsNoMoreData(err) {
				return nil
			}

			return err
		}

		if !predicate(value, key) {
			return nil
		}
	}
}
