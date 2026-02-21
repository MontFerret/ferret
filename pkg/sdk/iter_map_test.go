package sdk

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/smartystreets/goconvey/convey"
)

func TestMapIterator(t *testing.T) {
	convey.Convey("MapIterator", t, func() {
		data := map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		}

		iter := NewMapIterator(data)

		keys := make(map[string]bool)
		values := make(map[string]bool)

		err := runtime.ForEachIter(t.Context(), iter, func(ctx context.Context, value, idx runtime.Value) (runtime.Boolean, error) {
			keys[value.String()] = true
			values[idx.String()] = true

			return runtime.True, nil
		})

		convey.So(err, convey.ShouldBeNil)
		convey.So(keys, convey.ShouldResemble, map[string]bool{"a": true, "b": true, "c": true})
		convey.So(values, convey.ShouldResemble, map[string]bool{"1": true, "2": true, "3": true})
	})
}
