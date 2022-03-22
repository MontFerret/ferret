package runtime_test

import (
	"context"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestNewOptions(t *testing.T) {
	Convey("Should return a Option object with context, params and logs", t, func() {
		valueString := []core.Value{
			values.NewString("pseudo-user"),
		}

		optionTestWithParam := runtime.WithParam("test-param", valueString)
		optionTestWithParams := runtime.WithParams(make(map[string]interface{}))
		optionWithLog := runtime.WithLog(&io.PipeWriter{})
		optionWithLogLevel := runtime.WithLogLevel(logging.DebugLevel)
		optionWithLogFields := runtime.WithLogFields(make(map[string]interface{}))

		options := []runtime.Option{}

		options = append(options, optionTestWithParam)
		options = append(options, optionTestWithParams)
		options = append(options, optionWithLog)
		options = append(options, optionWithLogLevel)
		options = append(options, optionWithLogFields)

		opts := runtime.NewOptions(options)

		ctx := opts.WithContext(context.Background())
		So(ctx.Err(), ShouldBeNil)
	})

}
