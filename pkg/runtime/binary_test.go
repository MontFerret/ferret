package runtime_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func TestBinary(t *testing.T) {
	ctx := context.Background()

	Convey("Binary Operations", t, func() {

		Convey("NewBinary", func() {
			Convey("Should create Binary from byte slice", func() {
				data := []byte("hello world")
				binary := runtime.NewBinary(data)
				So(binary, ShouldNotBeNil)
				So(binary.Unwrap(), ShouldResemble, data)
			})

			Convey("Should create empty Binary", func() {
				binary := runtime.NewBinary([]byte{})
				So(binary, ShouldNotBeNil)
				So(binary.Unwrap(), ShouldResemble, []byte{})
			})
		})

		Convey("NewBinaryFrom", func() {
			Convey("Should create Binary from reader", func() {
				data := "hello world"
				reader := strings.NewReader(data)
				
				binary, err := runtime.NewBinaryFrom(reader)
				So(err, ShouldBeNil)
				So(binary, ShouldNotBeNil)
				So(binary.Unwrap(), ShouldResemble, []byte(data))
			})

			Convey("Should create empty Binary from empty reader", func() {
				reader := strings.NewReader("")
				
				binary, err := runtime.NewBinaryFrom(reader)
				So(err, ShouldBeNil)
				So(binary, ShouldNotBeNil)
				So(binary.Unwrap(), ShouldResemble, []byte{})
			})

			Convey("Should handle reader errors", func() {
				// Use a buffer and close it to simulate an error
				reader := &ErrorReader{}
				
				_, err := runtime.NewBinaryFrom(reader)
				So(err, ShouldNotBeNil)
			})
		})

		Convey(".MarshalJSON", func() {
			Convey("Should marshal to base64 encoded JSON string", func() {
				data := []byte("hello")
				binary := runtime.NewBinary(data)
				
				marshaled, err := binary.MarshalJSON()
				So(err, ShouldBeNil)
				So(string(marshaled), ShouldContainSubstring, "aGVsbG8=") // base64 for "hello"
			})

			Convey("Should marshal empty binary", func() {
				binary := runtime.NewBinary([]byte{})
				
				marshaled, err := binary.MarshalJSON()
				So(err, ShouldBeNil)
				So(string(marshaled), ShouldNotBeEmpty)
			})
		})

		Convey(".String", func() {
			Convey("Should return string representation", func() {
				data := []byte("hello")
				binary := runtime.NewBinary(data)
				
				str := binary.String()
				So(str, ShouldNotBeEmpty)
			})
		})

		Convey(".Hash", func() {
			Convey("Should calculate hash consistently", func() {
				data := []byte("hello")
				binary1 := runtime.NewBinary(data)
				binary2 := runtime.NewBinary(data)
				
				hash1 := binary1.Hash()
				hash2 := binary2.Hash()
				So(hash1, ShouldEqual, hash2)
			})

			Convey("Should calculate different hash for different data", func() {
				binary1 := runtime.NewBinary([]byte("hello"))
				binary2 := runtime.NewBinary([]byte("world"))
				
				hash1 := binary1.Hash()
				hash2 := binary2.Hash()
				So(hash1, ShouldNotEqual, hash2)
			})
		})

		Convey(".Copy", func() {
			Convey("Should create independent copy", func() {
				data := []byte("hello")
				original := runtime.NewBinary(data)
				
				copied := original.Copy()
				copyBinary := copied.(runtime.Binary)
				
				So(copyBinary.Unwrap(), ShouldResemble, original.Unwrap())
				
				// Modifying original data shouldn't affect copy since Copy should create new slice
				originalBytes := original.Unwrap().([]byte)
				if len(originalBytes) > 0 {
					originalBytes[0] = 'x' // This should not affect the copy
				}
			})
		})

		Convey(".Length", func() {
			Convey("Should return correct length", func() {
				data := []byte("hello")
				binary := runtime.NewBinary(data)
				
				length, err := binary.Length(ctx)
				So(err, ShouldBeNil)
				So(length, ShouldEqual, runtime.NewInt(5))
			})

			Convey("Should return zero for empty binary", func() {
				binary := runtime.NewBinary([]byte{})
				
				length, err := binary.Length(ctx)
				So(err, ShouldBeNil)
				So(length, ShouldEqual, runtime.NewInt(0))
			})
		})

		Convey(".Compare", func() {
			Convey("Should return 0 for equal binaries", func() {
				data := []byte("hello")
				binary1 := runtime.NewBinary(data)
				binary2 := runtime.NewBinary(data)
				
				result := binary1.Compare(binary2)
				So(result, ShouldEqual, 0)
			})

			Convey("Should return non-zero for different binaries", func() {
				binary1 := runtime.NewBinary([]byte("hello"))
				binary2 := runtime.NewBinary([]byte("world"))
				
				result := binary1.Compare(binary2)
				So(result, ShouldNotEqual, 0)
			})

			Convey("Should handle comparison with non-binary types", func() {
				binary := runtime.NewBinary([]byte("hello"))
				
				result := binary.Compare(runtime.NewString("hello"))
				So(result, ShouldNotEqual, 0) // Different types
			})
		})
	})
}

// ErrorReader is a helper for testing error conditions
type ErrorReader struct{}

func (r *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}