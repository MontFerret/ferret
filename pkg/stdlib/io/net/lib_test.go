package net_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/io/net"
)

func TestRegisterLib(t *testing.T) {
	Convey("Should register NET namespace functions", t, func() {
		ns := runtime.NewRootNamespace()

		err := net.RegisterLib(ns)

		So(err, ShouldBeNil)

		// Verify that functions were registered by checking registered function names
		functions := ns.Functions().Build()
		So(len(functions.Names()), ShouldBeGreaterThan, 0)

		// Check that HTTP functions are registered
		hasGet := false
		hasPost := false
		hasPut := false
		hasDelete := false
		hasDo := false

		for _, fn := range functions.Names() {
			if fn == "NET::HTTP::GET" {
				hasGet = true
			}
			if fn == "NET::HTTP::POST" {
				hasPost = true
			}
			if fn == "NET::HTTP::PUT" {
				hasPut = true
			}
			if fn == "NET::HTTP::DELETE" {
				hasDelete = true
			}
			if fn == "NET::HTTP::DO" {
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
