package base

import (
	"fmt"

	"github.com/pkg/errors"
)

func Try[T any](f func() T) (T, error) {
	var v T
	var e error

	func() {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case error:
					e = x
				case string:
					e = errors.New(x)
				default:
					e = errors.New(fmt.Sprintf("unknown error: %v", x))
				}
			}
		}()

		f()
	}()

	return v, e
}
