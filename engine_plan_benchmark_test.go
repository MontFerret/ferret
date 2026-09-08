package ferret

import "testing"

func BenchmarkEngineCompilePlan(b *testing.B) {
	for _, debug := range []bool{false, true} {
		name := "ordinary"
		if debug {
			name = "debug"
		}

		b.Run(name, func(b *testing.B) {
			engine, err := New()
			if err != nil {
				b.Fatal(err)
			}

			defer func() { _ = engine.Close() }()
			compile := engine.Compile
			if debug {
				compile = engine.CompileDebug
			}

			src := NewAnonymousSource("RETURN @value + 1")
			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				plan, err := compile(b.Context(), src)
				if err != nil {
					b.Fatal(err)
				}

				if err := plan.Close(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkEngineLoadPlan(b *testing.B) {
	engine, err := New()
	if err != nil {
		b.Fatal(err)
	}

	defer func() { _ = engine.Close() }()
	plan, err := engine.Compile(b.Context(), NewAnonymousSource("RETURN @value + 1"))
	if err != nil {
		b.Fatal(err)
	}

	data, err := plan.Marshal()
	if err != nil {
		b.Fatal(err)
	}

	if err := plan.Close(); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		loaded, err := engine.Load(data)
		if err != nil {
			b.Fatal(err)
		}

		if err := loaded.Close(); err != nil {
			b.Fatal(err)
		}
	}
}
