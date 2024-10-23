package compiler_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/parser"
)

func TestMember(t *testing.T) {
	RunUseCases(t, []UseCase{
		{
			"LET arr = [1,2,3,4] RETURN arr[10]",
			nil,
			nil,
		},
		{
			"LET arr = [1,2,3,4] RETURN arr[1]",
			2,
			nil,
		},
		{
			"LET arr = [1,2,3,4] LET idx = 1 RETURN arr[idx]",
			2,
			nil,
		},
		{
			`LET obj = { foo: "bar", qaz: "wsx"} RETURN obj["qaz"]`,
			"wsx",
			nil,
		},
		{
			fmt.Sprintf(`
								LET obj = { "foo": "bar", %s: "wsx"}
		
								RETURN obj["qaz"]
							`, "`qaz`"),
			"wsx",
			nil,
		},
		{
			fmt.Sprintf(`
								LET obj = { "foo": "bar", %s: "wsx"}
		
								RETURN obj["let"]
							`, "`let`"),
			"wsx",
			nil,
		},
		{
			`LET obj = { foo: "bar", qaz: "wsx"} LET key = "qaz" RETURN obj[key]`,
			"wsx",
			nil,
		},
		{
			`RETURN { foo: "bar" }.foo`,
			"bar",
			nil,
		},
		{
			`LET inexp = 1 IN {'foo': [1]}.foo
			LET ternaryexp = FALSE ? TRUE : {foo: TRUE}.foo
			RETURN inexp && ternaryexp`,
			true,
			nil,
		},
		{
			`RETURN ["bar", "foo"][0]`,
			"bar",
			nil,
		},
		{
			`LET inexp = 1 IN [[1]][0]
								LET ternaryexp = FALSE ? TRUE : [TRUE][0]
								RETURN inexp && ternaryexp`,
			true,
			nil,
		},
		{
			`LET obj = {
							first: {
								second: {
									third: {
										fourth: {
											fifth: {
												bottom: true
											}
										}
									}
								}
							}
						}
		
						RETURN obj.first.second.third.fourth.fifth.bottom`,
			true,
			nil,
		},
		{
			`LET o1 = {
		  first: {
		      second: {
		          ["third"]: {
		              fourth: {
		                  fifth: {
		                      bottom: true
		                  }
		              }
		          }
		      }
		  }
		}
		
		LET o2 = { prop: "third" }
		
		RETURN o1["first"]["second"][o2.prop]["fourth"]["fifth"].bottom`,

			true,
			nil,
		},
		{
			`LET o1 = {
		  first: {
		      second: {
		          third: {
		              fourth: {
		                  fifth: {
		                      bottom: true
		                  }
		              }
		          }
		      }
		  }
		}
		
		LET o2 = { prop: "third" }
		
		RETURN o1.first["second"][o2.prop].fourth["fifth"]["bottom"]`,

			true,
			nil,
		},
		{
			`LET obj = {
							attributes: {
								'data-index': 1
							}
						}
		
						RETURN obj.attributes['data-index']`,
			1,
			nil,
		},
		{
			`LET obj = NONE RETURN obj.foo`,
			nil,
			ShouldBeError,
		},
		{
			`LET obj = NONE RETURN obj?.foo`,
			nil,
			nil,
		},
	})

	//		Convey("ObjectDecl by literal passed to func call", func() {
	//			c := compiler.New()
	//
	//			p, err := c.Compile(`
	//				RETURN KEEP_KEYS({first: {second: "third"}}.first, "second")
	//			`)
	//			So(err, ShouldBeNil)
	//
	//			out, err := p.Run(context.Background())
	//			So(err, ShouldBeNil)
	//
	//			So(string(out), ShouldEqual, `{"second":"third"}`)
	//		})
	//
	//		Convey("ObjectDecl by literal as forSource", func() {
	//			c := compiler.New()
	//
	//			p, err := c.Compile(`
	//				FOR v, k IN {f: {foo: "bar"}}.f
	//					RETURN [k, v]
	//			`)
	//			So(err, ShouldBeNil)
	//
	//			out, err := p.Run(context.Background())
	//			So(err, ShouldBeNil)
	//
	//			So(string(out), ShouldEqual, `[["foo","bar"]]`)
	//		})
	//

	//
	//		Convey("ArrayDecl by literal passed to func call", func() {
	//			c := compiler.New()
	//
	//			p, err := c.Compile(`
	//				RETURN FIRST([[1, 2]][0])
	//			`)
	//			So(err, ShouldBeNil)
	//
	//			out, err := p.Run(context.Background())
	//			So(err, ShouldBeNil)
	//
	//			So(string(out), ShouldEqual, `1`)
	//		})
	//
	//		Convey("ArrayDecl by literal as forSource", func() {
	//			c := compiler.New()
	//
	//			p, err := c.Compile(`
	//				FOR i IN [[1, 2]][0]
	//					RETURN i
	//			`)
	//			So(err, ShouldBeNil)
	//
	//			out, err := p.Run(context.Background())
	//			So(err, ShouldBeNil)
	//
	//			So(string(out), ShouldEqual, `[1,2]`)
	//		})
	//

	//
	//		Convey("Prop after a func call", func() {
	//			c := compiler.New()
	//
	//			p, err := c.Compile(`
	//				LET arr = [{ name: "Bob" }]
	//
	//				RETURN FIRST(arr).name
	//			`)
	//
	//			So(err, ShouldBeNil)
	//
	//			out, err := p.Run(context.Background())
	//
	//			So(err, ShouldBeNil)
	//
	//			So(string(out), ShouldEqual, `"Bob"`)
	//		})
	//
	//		Convey("Computed prop after a func call", func() {
	//			c := compiler.New()
	//
	//			p, err := c.Compile(`
	//				LET arr = [{ name: { first: "Bob" } }]
	//
	//				RETURN FIRST(arr)['name'].first
	//			`)
	//
	//			So(err, ShouldBeNil)
	//
	//			out, err := p.Run(context.Background())
	//
	//			So(err, ShouldBeNil)
	//
	//			So(string(out), ShouldEqual, `"Bob"`)
	//		})
	//

	//
	//	Convey("Optional chaining", t, func() {
	//		Convey("Object", func() {
	//			Convey("When value does not exist", func() {
	//				c := compiler.New()
	//
	//				p, err := c.Compile(`
	//				LET obj = { foo: None }
	//
	//				RETURN obj.foo?.bar
	//			`)
	//
	//				So(err, ShouldBeNil)
	//
	//				out, err := p.Run(context.Background())
	//
	//				So(err, ShouldBeNil)
	//
	//				So(string(out), ShouldEqual, `null`)
	//			})
	//
	//			Convey("When value does exists", func() {
	//				c := compiler.New()
	//
	//				p, err := c.Compile(`
	//				LET obj = { foo: { bar: "bar" } }
	//
	//				RETURN obj.foo?.bar
	//			`)
	//
	//				So(err, ShouldBeNil)
	//
	//				out, err := p.Run(context.Background())
	//
	//				So(err, ShouldBeNil)
	//
	//				So(string(out), ShouldEqual, `"bar"`)
	//			})
	//		})
	//
	//		Convey("Array", func() {
	//			Convey("When value does not exist", func() {
	//				c := compiler.New()
	//
	//				p, err := c.Compile(`
	//				LET obj = { foo: None }
	//
	//				RETURN obj.foo?.bar?.[0]
	//			`)
	//
	//				So(err, ShouldBeNil)
	//
	//				out, err := p.Run(context.Background())
	//
	//				So(err, ShouldBeNil)
	//
	//				So(string(out), ShouldEqual, `null`)
	//			})
	//
	//			Convey("When value does exists", func() {
	//				c := compiler.New()
	//
	//				p, err := c.Compile(`
	//				LET obj = { foo: { bar: ["bar"] } }
	//
	//				RETURN obj.foo?.bar?.[0]
	//			`)
	//
	//				So(err, ShouldBeNil)
	//
	//				out, err := p.Run(context.Background())
	//
	//				So(err, ShouldBeNil)
	//
	//				So(string(out), ShouldEqual, `"bar"`)
	//			})
	//		})
	//
	//		Convey("Function", func() {
	//			Convey("When value does not exist", func() {
	//				c := compiler.New()
	//
	//				p, err := c.Compile(`
	//				RETURN FIRST([])?.foo
	//			`)
	//
	//				So(err, ShouldBeNil)
	//
	//				out, err := p.Run(context.Background())
	//
	//				So(err, ShouldBeNil)
	//
	//				So(string(out), ShouldEqual, `null`)
	//			})
	//
	//			Convey("When value does exists", func() {
	//				c := compiler.New()
	//
	//				p, err := c.Compile(`
	//				RETURN FIRST([{ foo: "bar" }])?.foo
	//			`)
	//
	//				So(err, ShouldBeNil)
	//
	//				out, err := p.Run(context.Background())
	//
	//				So(err, ShouldBeNil)
	//
	//				So(string(out), ShouldEqual, `"bar"`)
	//			})
	//
	//			Convey("When function returns error", func() {
	//				c := compiler.New()
	//				c.RegisterFunction("ERROR", func(ctx context.visitor, args ...core.Value) (core.Value, error) {
	//					return nil, core.ErrNotImplemented
	//				})
	//
	//				p, err := c.Compile(`
	//				RETURN ERROR()?.foo
	//			`)
	//
	//				So(err, ShouldBeNil)
	//
	//				out, err := p.Run(context.Background())
	//
	//				So(err, ShouldBeNil)
	//
	//				So(string(out), ShouldEqual, `null`)
	//			})
	//		})
	//	})
	//
	Convey("Reserved words as property name", t, func() {
		p := parser.New("RETURN TRUE")

		r := regexp.MustCompile(`\w+`)

		for idx, l := range p.GetLiteralNames() {
			if r.MatchString(l) {
				query := strings.Builder{}
				query.WriteString("LET o = {\n")
				query.WriteString(l[1 : len(l)-1])
				query.WriteString(":")
				query.WriteString(strconv.Itoa(idx))
				query.WriteString(",\n")
				query.WriteString("}\n")
				query.WriteString("RETURN o")

				expected := strings.Builder{}
				expected.WriteString("{")
				expected.WriteString(strings.ReplaceAll(l, "'", "\""))
				expected.WriteString(":")
				expected.WriteString(strconv.Itoa(idx))
				expected.WriteString("}")

				c := compiler.New()
				prog, err := c.Compile(query.String())

				So(err, ShouldBeNil)

				out, err := Exec(prog, true, runtime.WithFunctions(c.Functions().Unwrap()))

				So(err, ShouldBeNil)
				So(out, ShouldEqual, expected.String())
			}
		}
	})
}

//func BenchmarkMemberArray(b *testing.B) {
//	p := compiler.New().MustCompile(`
//				LET arr = [[[[1]]]]
//
//				RETURN arr[0][0][0][0]
//			`)
//
//	for n := 0; n < b.N; n++ {
//		p.Run(context.Background())
//	}
//}
//
//func BenchmarkMemberObject(b *testing.B) {
//	p := compiler.New().MustCompile(`
//				LET obj = {
//					first: {
//						second: {
//							third: {
//								fourth: {
//									fifth: {
//										bottom: true
//									}
//								}
//							}
//						}
//					}
//				}
//
//				RETURN obj.first.second.third.fourth.fifth.bottom
//			`)
//
//	for n := 0; n < b.N; n++ {
//		p.Run(context.Background())
//	}
//}
//
//func BenchmarkMemberObjectComputed(b *testing.B) {
//	p := compiler.New().MustCompile(`
//				LET obj = {
//					first: {
//						second: {
//							third: {
//								fourth: {
//									fifth: {
//										bottom: true
//									}
//								}
//							}
//						}
//					}
//				}
//
//				RETURN obj["first"]["second"]["third"]["fourth"]["fifth"]["bottom"]
//			`)
//
//	for n := 0; n < b.N; n++ {
//		p.Run(context.Background())
//	}
//}
