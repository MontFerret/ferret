package fs_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/io/fs"
)

func TestRegisterLib(t *testing.T) {
	Convey("Should register FS namespace functions", t, func() {
		ns := runtime.NewRootNamespace()
		
		err := fs.RegisterLib(ns)
		
		So(err, ShouldBeNil)
		
		// Verify that functions were registered by checking registered function names
		functions := ns.RegisteredFunctions()
		So(len(functions), ShouldBeGreaterThan, 0)
		
		// Check that FS functions are registered
		hasRead := false
		hasWrite := false
		
		for _, fn := range functions {
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