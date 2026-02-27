package fs_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/fs"
)

func TestRegisterLib(t *testing.T) {
	Convey("Should register FS namespace functions", t, func() {
		ns := runtime.NewRootNamespace()

		fs.RegisterLib(ns)

		// Verify that functions were registered by checking registered function names
		functions, err := ns.Build()
		So(err, ShouldBeNil)
		names := functions.List()
		So(len(names), ShouldBeGreaterThan, 0)

		// Check that FS functions are registered
		hasRead := false
		hasWrite := false

		for _, fn := range names {
			if fn == "FS::READ" {
				hasRead = true
			}
			if fn == "FS::WRITE" {
				hasWrite = true
			}
		}

		So(hasRead, ShouldBeTrue)
		So(hasWrite, ShouldBeTrue)
	})
}
