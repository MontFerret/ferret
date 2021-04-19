package compiler_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMember(t *testing.T) {
	Convey("Computed properties", t, func() {
		Convey("Array by literal", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET arr = [1,2,3,4]

				RETURN arr[1]
			`)

			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `2`)
		})

		Convey("Array by variable", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET arr = [1,2,3,4]
				LET idx = 1

				RETURN arr[idx]
			`)

			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `2`)
		})

		Convey("Object by literal", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET obj = { foo: "bar", qaz: "wsx"}

				RETURN obj["qaz"]
			`)

			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"wsx"`)
		})

		Convey("Object by literal with property defined as a string", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET obj = { "foo": "bar", "qaz": "wsx"}

				RETURN obj["qaz"]
			`)

			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"wsx"`)
		})

		Convey("Object by literal with property defined as a multi line string", func() {
			c := compiler.New()

			p, err := c.Compile(fmt.Sprintf(`
				LET obj = { "foo": "bar", %s: "wsx"}

				RETURN obj["qaz"]
			`, "`qaz`"))

			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"wsx"`)
		})

		Convey("Object by variable", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET obj = { foo: "bar", qaz: "wsx"}
				LET key = "qaz"

				RETURN obj[key]
			`)

			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"wsx"`)
		})

		Convey("ObjectDecl by literal", func() {
			c := compiler.New()

			p, err := c.Compile(`
				RETURN { foo: "bar" }.foo
			`)
			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())
			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"bar"`)
		})

		Convey("ObjectDecl by literal passed to func call", func() {
			c := compiler.New()

			p, err := c.Compile(`
				RETURN KEEP_KEYS({first: {second: "third"}}.first, "second")
			`)
			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())
			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `{"second":"third"}`)
		})

		Convey("ObjectDecl by literal as forSource", func() {
			c := compiler.New()

			p, err := c.Compile(`
				FOR v, k IN {f: {foo: "bar"}}.f
					RETURN [k, v]
			`)
			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())
			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `[["foo","bar"]]`)
		})

		Convey("ObjectDecl by literal as expression", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET inexp = 1 IN {'foo': [1]}.foo
				LET ternaryexp = FALSE ? TRUE : {foo: TRUE}.foo
				RETURN inexp && ternaryexp
			`)
			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())
			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `true`)
		})

		Convey("ArrayDecl by literal", func() {
			c := compiler.New()

			p, err := c.Compile(`
				RETURN ["bar", "foo"][0]
			`)
			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())
			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"bar"`)
		})

		Convey("ArrayDecl by literal passed to func call", func() {
			c := compiler.New()

			p, err := c.Compile(`
				RETURN FIRST([[1, 2]][0])
			`)
			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())
			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `1`)
		})

		Convey("ArrayDecl by literal as forSource", func() {
			c := compiler.New()

			p, err := c.Compile(`
				FOR i IN [[1, 2]][0]
					RETURN i
			`)
			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())
			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `[1,2]`)
		})

		Convey("ArrayDecl by literal as expression", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET inexp = 1 IN [[1]][0]
				LET ternaryexp = FALSE ? TRUE : [TRUE][0]
				RETURN inexp && ternaryexp
			`)
			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())
			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `true`)
		})

		Convey("Deep path", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET obj = {
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

				RETURN obj.first.second.third.fourth.fifth.bottom
			`)

			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `true`)
		})

		Convey("Deep computed path", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET obj = {
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

				RETURN obj["first"]["second"]["third"]["fourth"]["fifth"].bottom
			`)

			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `true`)
		})

		Convey("Prop after a func call", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET arr = [{ name: "Bob" }]

				RETURN FIRST(arr).name
			`)

			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"Bob"`)
		})

		Convey("Computed prop after a func call", func() {
			c := compiler.New()

			p, err := c.Compile(`
				LET arr = [{ name: { first: "Bob" } }]

				RETURN FIRST(arr)['name'].first
			`)

			So(err, ShouldBeNil)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, `"Bob"`)
		})

		Convey("Computed property with quotes", func() {
			c := compiler.New()

			p := c.MustCompile(`
				LET obj = {
					attributes: {
						'data-index': 1
					}
				}
				
				RETURN obj.attributes['data-index']
			`)

			out, err := p.Run(context.Background())

			So(err, ShouldBeNil)

			So(string(out), ShouldEqual, "1")
		})
	})
}

func BenchmarkMemberArray(b *testing.B) {
	p := compiler.New().MustCompile(`
				LET arr = [[[[1]]]]

				RETURN arr[0][0][0][0]
			`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkMemberObject(b *testing.B) {
	p := compiler.New().MustCompile(`
				LET obj = {
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

				RETURN obj.first.second.third.fourth.fifth.bottom
			`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkMemberObjectComputed(b *testing.B) {
	p := compiler.New().MustCompile(`
				LET obj = { "foo": "bar"}

				RETURN obj["foo"]
			`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
