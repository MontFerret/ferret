package sdk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestProxyRemovable(t *testing.T) {
	Convey("Proxy should delegate removable interfaces", t, func() {
		ctx := t.Context()

		Convey("Should remove by index", func() {
			var _ runtime.IndexRemovable = (*Proxy[*runtime.Array])(nil)

			arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3))
			proxy := NewProxy(arr)

			removed, err := proxy.RemoveAt(ctx, 1)

			So(err, ShouldBeNil)
			So(removed, ShouldEqual, runtime.NewInt(2))

			after, _ := arr.At(ctx, 1)
			So(after, ShouldEqual, runtime.NewInt(3))
		})

		Convey("Should remove indexed targets by integer key", func() {
			arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3))
			proxy := NewProxy(arr)

			err := proxy.RemoveKey(ctx, runtime.NewInt(1))

			So(err, ShouldBeNil)

			after, _ := arr.At(ctx, 1)
			So(after, ShouldEqual, runtime.NewInt(3))
		})

		Convey("Should remove by key", func() {
			var _ runtime.KeyRemovable = (*Proxy[*runtime.Object])(nil)

			obj := runtime.NewObjectWith(map[string]runtime.Value{
				"one": runtime.NewInt(1),
				"two": runtime.NewInt(2),
			})
			proxy := NewProxy(obj)

			err := proxy.RemoveKey(ctx, runtime.NewString("one"))

			So(err, ShouldBeNil)

			_, found, lookupErr := obj.Lookup(ctx, runtime.NewString("one"))
			So(lookupErr, ShouldBeNil)
			So(found, ShouldBeFalse)
		})

		Convey("Should remove by value", func() {
			var _ runtime.ValueRemovable = (*Proxy[*runtime.Array])(nil)

			arr := runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3))
			proxy := NewProxy(arr)

			err := proxy.Remove(ctx, runtime.NewInt(2))

			So(err, ShouldBeNil)

			after, _ := arr.At(ctx, 1)
			So(after, ShouldEqual, runtime.NewInt(3))
		})
	})
}

func TestProxyMapRemovable(t *testing.T) {
	Convey("ProxyMap should implement removable interfaces", t, func() {
		ctx := t.Context()

		var _ runtime.KeyRemovable = (*ProxyMap[string, int])(nil)
		var _ runtime.ValueRemovable = (*ProxyMap[string, int])(nil)

		Convey("Should remove by key", func() {
			data := map[string]int{
				"one": 1,
				"two": 2,
			}
			proxy := NewProxyMap(data)

			err := proxy.RemoveKey(ctx, runtime.NewString("one"))

			So(err, ShouldBeNil)
			So(proxy.Target(), ShouldResemble, map[string]int{"two": 2})
		})

		Convey("Should remove by value", func() {
			data := map[string]int{
				"one": 1,
				"two": 2,
			}
			proxy := NewProxyMap(data)

			err := proxy.Remove(ctx, runtime.NewInt(2))

			So(err, ShouldBeNil)
			So(proxy.Target(), ShouldResemble, map[string]int{"one": 1})
		})
	})
}

func TestProxySliceRemovable(t *testing.T) {
	Convey("ProxySlice should implement removable interfaces", t, func() {
		ctx := t.Context()

		var _ runtime.IndexRemovable = (*ProxySlice[int])(nil)
		var _ runtime.ValueRemovable = (*ProxySlice[int])(nil)

		Convey("Should remove by index", func() {
			proxy := NewProxySlice([]int{1, 2, 3})

			removed, err := proxy.RemoveAt(ctx, 1)

			So(err, ShouldBeNil)
			So(removed.String(), ShouldEqual, "2")
			So(proxy.Target(), ShouldResemble, []int{1, 3})
		})

		Convey("Should no-op on missing positive index", func() {
			proxy := NewProxySlice([]int{1, 2, 3})

			removed, err := proxy.RemoveAt(ctx, 100)

			So(err, ShouldBeNil)
			So(removed, ShouldEqual, runtime.None)
			So(proxy.Target(), ShouldResemble, []int{1, 2, 3})
		})

		Convey("Should remove by value", func() {
			proxy := NewProxySlice([]int{1, 2, 3})

			err := proxy.Remove(ctx, runtime.NewInt(2))

			So(err, ShouldBeNil)
			So(proxy.Target(), ShouldResemble, []int{1, 3})
		})

		Convey("Should remove proxy-wrapped runtime values", func() {
			first := NewProxy(1)
			second := NewProxy(2)
			proxy := NewProxySlice([]runtime.Value{first, second})

			err := proxy.Remove(ctx, runtime.NewInt(2))

			So(err, ShouldBeNil)
			So(proxy.Target(), ShouldResemble, []runtime.Value{first})
		})

		Convey("Should remove by integer key", func() {
			proxy := NewProxySlice([]int{1, 2, 3})

			err := proxy.RemoveKey(ctx, runtime.NewInt(1))

			So(err, ShouldBeNil)
			So(proxy.Target(), ShouldResemble, []int{1, 3})
		})

		Convey("Should reject non-index keys", func() {
			proxy := NewProxySlice([]int{1, 2, 3})

			err := proxy.RemoveKey(ctx, runtime.NewString("1"))

			So(err, ShouldNotBeNil)
			So(proxy.Target(), ShouldResemble, []int{1, 2, 3})
		})
	})
}
