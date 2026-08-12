package net_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io/net"
)

func TestRegisterLib(t *testing.T) {
	Convey("Should register NET namespace functions", t, func() {
		ns := runtime.NewLibrary()

		net.RegisterLib(ns)

		funcs, err := ns.Build()
		So(err, ShouldBeNil)

		// Verify that functions were registered by checking registered function names
		functions := funcs.List()
		So(len(functions), ShouldBeGreaterThan, 0)

		// Check that HTTP functions are registered
		hasGet := false
		hasPost := false
		hasPut := false
		hasDelete := false
		hasDo := false

		for _, fn := range functions {
			if fn == "net::http::get" {
				hasGet = true
			}
			if fn == "net::http::post" {
				hasPost = true
			}
			if fn == "net::http::put" {
				hasPut = true
			}
			if fn == "net::http::delete" {
				hasDelete = true
			}
			if fn == "net::http::do" {
				hasDo = true
			}
		}

		So(hasGet, ShouldBeTrue)
		So(hasPost, ShouldBeTrue)
		So(hasPut, ShouldBeTrue)
		So(hasDelete, ShouldBeTrue)
		So(hasDo, ShouldBeTrue)
	})
}
