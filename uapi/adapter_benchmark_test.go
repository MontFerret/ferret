package uapi

import (
	"testing"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkCompile(b *testing.B) {
	for _, query := range []struct{ name, text string }{
		{"ASCII", "RETURN @value + 1"},
		{"Unicode", "LET prefix = \"é\"\nRETURN [prefix, @value]"},
	} {
		for _, adapter := range []string{"Native", "Universal"} {
			b.Run(query.name+"/"+adapter, func(b *testing.B) {
				native := newTestEngine(b)
				portable := Wrap(native)
				b.Cleanup(func() { _ = portable.Close() })
				b.ReportAllocs()
				b.ResetTimer()
				for b.Loop() {
					if adapter == "Native" {
						p, err := native.Compile(b.Context(), source.NewAnonymous(query.text))
						if err != nil {
							b.Fatal(err)
						}

						if err := p.Close(); err != nil {
							b.Fatal(err)
						}
					} else {
						p, err := portable.Compile(b.Context(), api.NewAnonymousSource(query.text))
						if err != nil {
							b.Fatal(err)
						}

						if err := p.Close(); err != nil {
							b.Fatal(err)
						}
					}
				}
			})
		}
	}
}

func BenchmarkSessionCreation(b *testing.B) {
	for _, debug := range []bool{false, true} {
		for _, options := range []bool{false, true} {
			for _, adapter := range []string{"Native", "Universal"} {
				name := map[bool]string{false: "ordinary", true: "debug"}[debug] + "/" + map[bool]string{false: "defaults", true: "options"}[options] + "/" + adapter
				b.Run(name, func(b *testing.B) {
					native := newTestEngine(b)
					compiled, err := native.CompileDebug(b.Context(), source.NewAnonymous("RETURN @value"))
					if err != nil {
						b.Fatal(err)
					}

					portable := newPlan(compiled)
					b.Cleanup(func() { _ = portable.Close() })
					var nativeOptions []engine.SessionOption
					var universalOptions []api.SessionOption
					if options {
						nativeOptions = []engine.SessionOption{engine.WithSessionParam("value", 42), engine.WithOutputContentType("application/json")}
						universalOptions = []api.SessionOption{api.WithParam("value", 42), api.WithOutputContentType("application/json")}
					}

					b.ReportAllocs()
					b.ResetTimer()
					for b.Loop() {
						switch {
						case adapter == "Native" && !debug:
							s, err := compiled.NewSession(b.Context(), nativeOptions...)
							if err != nil {
								b.Fatal(err)
							}

							if err := s.Close(); err != nil {
								b.Fatal(err)
							}
						case adapter == "Native":
							s, err := compiled.NewDebugSession(b.Context(), nativeOptions...)
							if err != nil {
								b.Fatal(err)
							}

							if err := s.Close(); err != nil {
								b.Fatal(err)
							}
						case !debug:
							s, err := portable.NewSession(b.Context(), universalOptions...)
							if err != nil {
								b.Fatal(err)
							}

							if err := s.Close(); err != nil {
								b.Fatal(err)
							}
						default:
							s, err := portable.NewDebugSession(b.Context(), universalOptions...)
							if err != nil {
								b.Fatal(err)
							}

							if err := s.Close(); err != nil {
								b.Fatal(err)
							}
						}
					}
				})
			}
		}
	}
}

func BenchmarkReusableSession(b *testing.B) {
	for _, adapter := range []string{"Native", "Universal"} {
		b.Run(adapter, func(b *testing.B) {
			native := newTestEngine(b)
			p, err := native.Compile(b.Context(), source.NewAnonymous("RETURN @value + 1"))
			if err != nil {
				b.Fatal(err)
			}

			b.Cleanup(func() { _ = p.Close() })
			s, err := p.NewSession(b.Context(), engine.WithSessionParam("value", 41))
			if err != nil {
				b.Fatal(err)
			}

			portable := &session{native: s}
			b.Cleanup(func() { _ = portable.Close() })
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				if adapter == "Native" {
					out, err := s.Run(b.Context())
					if err != nil || out == nil || string(out.Content) != "42" {
						b.Fatalf("output=%+v err=%v", out, err)
					}
				} else {
					out, err := portable.Run(b.Context())
					if err != nil || out == nil || string(out.Content) != "42" {
						b.Fatalf("output=%+v err=%v", out, err)
					}
				}
			}
		})
	}
}

func BenchmarkRuntimeRun(b *testing.B) {
	for _, adapter := range []string{"Native", "Universal"} {
		b.Run(adapter, func(b *testing.B) {
			native := newTestEngine(b)
			portable := Wrap(native)
			b.Cleanup(func() { _ = portable.Close() })
			nativeSource := source.NewAnonymous("RETURN @value + 1")
			portableSource := api.NewAnonymousSource("RETURN @value + 1")
			nativeOption := engine.WithSessionParam("value", 41)
			portableOption := api.WithParam("value", 41)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				if adapter == "Native" {
					out, err := native.Run(b.Context(), nativeSource, nativeOption)
					if err != nil || out == nil || string(out.Content) != "42" {
						b.Fatalf("output=%+v err=%v", out, err)
					}
				} else {
					out, err := portable.Run(b.Context(), portableSource, portableOption)
					if err != nil || out == nil || string(out.Content) != "42" {
						b.Fatalf("output=%+v err=%v", out, err)
					}
				}
			}
		})
	}
}

func BenchmarkRuntimeConstruction(b *testing.B) {
	for _, adapter := range []string{"Native", "Universal"} {
		b.Run(adapter, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if adapter == "Native" {
					native, err := engine.New()
					if err != nil {
						b.Fatal(err)
					}

					if err := native.Close(); err != nil {
						b.Fatal(err)
					}
				} else {
					portable, err := New()
					if err != nil {
						b.Fatal(err)
					}

					if err := portable.Close(); err != nil {
						b.Fatal(err)
					}
				}
			}
		})
	}
}
