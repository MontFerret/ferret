package rnd

import "testing"

func BenchmarkSourceConstruction(b *testing.B) {
	for _, test := range []struct {
		new  func() *Source
		name string
	}{
		{name: "default", new: New},
		{name: "seeded", new: func() *Source { return NewSeed(42) }},
	} {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				sourceBenchmarkSink = test.new()
			}
		})
	}
}

func BenchmarkSourceConstructionFirstDraw(b *testing.B) {
	for _, test := range []struct {
		new  func() *Source
		name string
	}{
		{name: "default", new: New},
		{name: "seeded", new: func() *Source { return NewSeed(42) }},
	} {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				sourceBenchmarkSink = test.new()
				sourceBenchmarkSink.Float64()
			}
		})
	}
}

var sourceBenchmarkSink *Source
