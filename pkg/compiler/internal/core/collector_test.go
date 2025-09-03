package core_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestCollectSelector(t *testing.T) {
	Convey("CollectSelector", t, func() {
		Convey("NewCollectSelector", func() {
			Convey("Should create a new collect selector", func() {
				name := runtime.NewString("testField")
				
				selector := core.NewCollectSelector(name)
				
				So(selector, ShouldNotBeNil)
				So(selector.Name(), ShouldEqual, name)
			})

			Convey("Should handle empty string", func() {
				name := runtime.NewString("")
				
				selector := core.NewCollectSelector(name)
				
				So(selector, ShouldNotBeNil)
				So(selector.Name(), ShouldEqual, name)
			})

			Convey("Should handle different string values", func() {
				testCases := []string{
					"fieldName",
					"field.nested",
					"field[0]",
					"complex.field[0].nested",
					"123",
					"field_name",
					"FIELD_NAME",
				}

				for _, testCase := range testCases {
					name := runtime.NewString(testCase)
					selector := core.NewCollectSelector(name)
					
					So(selector.Name(), ShouldEqual, name)
					So(string(selector.Name()), ShouldEqual, testCase)
				}
			})
		})

		Convey(".Name", func() {
			Convey("Should return the selector name", func() {
				name := runtime.NewString("fieldName")
				selector := core.NewCollectSelector(name)
				
				retrievedName := selector.Name()
				
				So(retrievedName, ShouldEqual, name)
				So(string(retrievedName), ShouldEqual, "fieldName")
			})
		})
	})
}

func TestCollector(t *testing.T) {
	Convey("Collector", t, func() {
		Convey("NewCollector", func() {
			Convey("Should create a new collector", func() {
				dst := vm.Operand(1)
				projection := core.NewCollectorGroupProjection("testGroup")
				selectors := []*core.CollectSelector{
					core.NewCollectSelector(runtime.NewString("field1")),
					core.NewCollectSelector(runtime.NewString("field2")),
				}
				aggregation := core.NewCollectorAggregation(vm.Operand(2), nil)
				
				collector := core.NewCollector(
					core.CollectorTypeKeyGroup,
					dst,
					projection,
					selectors,
					aggregation,
				)
				
				So(collector, ShouldNotBeNil)
				So(collector.Type(), ShouldEqual, core.CollectorTypeKeyGroup)
				So(collector.Destination(), ShouldEqual, dst)
				So(collector.Projection(), ShouldEqual, projection)
				So(collector.GroupSelectors(), ShouldHaveLength, 2)
				So(collector.Aggregation(), ShouldEqual, aggregation)
			})

			Convey("Should handle nil parameters", func() {
				dst := vm.Operand(1)
				
				collector := core.NewCollector(
					core.CollectorTypeCounter,
					dst,
					nil,
					nil,
					nil,
				)
				
				So(collector, ShouldNotBeNil)
				So(collector.Type(), ShouldEqual, core.CollectorTypeCounter)
				So(collector.Destination(), ShouldEqual, dst)
				So(collector.Projection(), ShouldBeNil)
				So(collector.GroupSelectors(), ShouldBeNil)
				So(collector.Aggregation(), ShouldBeNil)
			})

			Convey("Should handle empty selectors", func() {
				collector := core.NewCollector(
					core.CollectorTypeKey,
					vm.Operand(1),
					nil,
					[]*core.CollectSelector{},
					nil,
				)
				
				So(collector.GroupSelectors(), ShouldHaveLength, 0)
			})
		})

		Convey("DetermineCollectorType", func() {
			Convey("Should return correct type for different combinations", func() {
				// Test all combinations
				testCases := []struct {
					withGrouping    bool
					withAggregation bool
					withProjection  bool
					withCounter     bool
					expected        core.CollectorType
				}{
					{true, false, false, true, core.CollectorTypeKeyCounter},
					{true, false, false, false, core.CollectorTypeKeyGroup},
					{true, true, false, false, core.CollectorTypeKeyGroup},
					{false, true, false, false, core.CollectorTypeKeyGroup},
					{false, false, true, false, core.CollectorTypeCounter},
					{false, false, false, false, core.CollectorTypeCounter},
					{false, false, false, true, core.CollectorTypeCounter},
				}

				for _, tc := range testCases {
					result := core.DetermineCollectorType(
						tc.withGrouping,
						tc.withAggregation,
						tc.withProjection,
						tc.withCounter,
					)
					
					So(result, ShouldEqual, tc.expected)
				}
			})

			Convey("Should prioritize grouping over aggregation", func() {
				result := core.DetermineCollectorType(true, true, false, false)
				
				So(result, ShouldEqual, core.CollectorTypeKeyGroup)
			})

			Convey("Should prioritize counter when grouping is present", func() {
				result := core.DetermineCollectorType(true, false, false, true)
				
				So(result, ShouldEqual, core.CollectorTypeKeyCounter)
			})
		})

		Convey("Collector Types", func() {
			Convey("Should have correct type constants", func() {
				// Just verify the constants exist and have expected values
				So(core.CollectorTypeCounter, ShouldEqual, 0)
				So(core.CollectorTypeKey, ShouldEqual, 1)
				So(core.CollectorTypeKeyCounter, ShouldEqual, 2)
				So(core.CollectorTypeKeyGroup, ShouldEqual, 3)
			})
		})

		Convey("Integration", func() {
			Convey("Should work with real selector data", func() {
				// Create selectors
				selectors := []*core.CollectSelector{
					core.NewCollectSelector(runtime.NewString("category")),
					core.NewCollectSelector(runtime.NewString("status")),
					core.NewCollectSelector(runtime.NewString("priority")),
				}
				
				// Determine collector type
				collectorType := core.DetermineCollectorType(true, false, true, false)
				
				// Create collector
				collector := core.NewCollector(
					collectorType,
					vm.Operand(10),
					nil,
					selectors,
					nil,
				)
				
				// Verify
				So(collector.Type(), ShouldEqual, core.CollectorTypeKeyGroup)
				So(collector.GroupSelectors(), ShouldHaveLength, 3)
				
				// Check selector names
				So(string(collector.GroupSelectors()[0].Name()), ShouldEqual, "category")
				So(string(collector.GroupSelectors()[1].Name()), ShouldEqual, "status")
				So(string(collector.GroupSelectors()[2].Name()), ShouldEqual, "priority")
			})

			Convey("Should handle complex collector configurations", func() {
				// Create multiple collectors with different types
				collectors := []*core.Collector{
					core.NewCollector(core.CollectorTypeCounter, vm.Operand(1), nil, nil, nil),
					core.NewCollector(core.CollectorTypeKey, vm.Operand(2), nil, []*core.CollectSelector{
						core.NewCollectSelector(runtime.NewString("key1")),
					}, nil),
					core.NewCollector(core.CollectorTypeKeyCounter, vm.Operand(3), nil, []*core.CollectSelector{
						core.NewCollectSelector(runtime.NewString("key2")),
						core.NewCollectSelector(runtime.NewString("key3")),
					}, nil),
				}
				
				// Verify each collector
				So(collectors[0].Type(), ShouldEqual, core.CollectorTypeCounter)
				So(collectors[1].Type(), ShouldEqual, core.CollectorTypeKey)
				So(collectors[2].Type(), ShouldEqual, core.CollectorTypeKeyCounter)
				
				// Verify destinations are different
				So(collectors[0].Destination(), ShouldEqual, vm.Operand(1))
				So(collectors[1].Destination(), ShouldEqual, vm.Operand(2))
				So(collectors[2].Destination(), ShouldEqual, vm.Operand(3))
			})
		})
	})
}