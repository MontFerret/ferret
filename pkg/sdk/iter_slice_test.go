package sdk

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/smartystreets/goconvey/convey"
)

func TestSliceIterator(t *testing.T) {
	convey.Convey("SliceIterator", t, func() {
		data := []string{"a", "b", "c"}

		iter := NewSliceIterator(data)
		values := make([]string, 0, len(data))

		err := runtime.ForEachIter(t.Context(), iter, func(ctx context.Context, value, idx runtime.Value) (runtime.Boolean, error) {
			values = append(values, value.String())
			return runtime.True, nil
		})

		convey.So(err, convey.ShouldBeNil)
		convey.So(values, convey.ShouldResemble, []string{"a", "b", "c"})
	})
}
