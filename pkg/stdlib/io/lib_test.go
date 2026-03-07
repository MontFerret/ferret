package io_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io"
)

func TestRegisterLib(t *testing.T) {
	Convey("Should register IO namespace functions", t, func() {
		ns := runtime.NewLibrary()

		io.RegisterLib(ns)

		funcs, err := ns.Build()
		So(err, ShouldBeNil)
		// Verify that functions were registered by checking registered function names
		functions := funcs.List()
		So(len(functions), ShouldBeGreaterThan, 0)

		// Check that FS functions are registered
		hasRead := false
		hasWrite := false
		hasGet := false
		hasPost := false
		hasPut := false
		hasDelete := false
		hasDo := false

		for _, fn := range functions {
			if fn == "IO::FS::READ" {
				hasRead = true
			}
			if fn == "IO::FS::WRITE" {
				hasWrite = true
			}
			if fn == "IO::NET::HTTP::GET" {
				hasGet = true
			}
			if fn == "IO::NET::HTTP::POST" {
				hasPost = true
			}
			if fn == "IO::NET::HTTP::PUT" {
				hasPut = true
			}
			if fn == "IO::NET::HTTP::DELETE" {
				hasDelete = true
			}
			if fn == "IO::NET::HTTP::DO" {
				hasDo = true
			}
		}

		So(hasRead, ShouldBeTrue)
		So(hasWrite, ShouldBeTrue)
		So(hasGet, ShouldBeTrue)
		So(hasPost, ShouldBeTrue)
		So(hasPut, ShouldBeTrue)
		So(hasDelete, ShouldBeTrue)
		So(hasDo, ShouldBeTrue)
	})
}
