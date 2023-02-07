package compiler_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/runtime"
)

func TestAggregate(t *testing.T) {
	Convey("Should aggregate values without grouping", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
  				COLLECT AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
  				RETURN {
    				minAge, 
    				maxAge 
  				}
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[{"maxAge":69,"minAge":25}]`)
	})

	Convey("Should aggregate values without grouping with multiple arguments", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
  				COLLECT AGGREGATE ages = UNION(u.age, u.age)
  				RETURN ages
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[[31,25,36,69,45,31,25,36,69,45]]`)
	})

	Convey("Should aggregate values with grouping", t, func() {
		c := compiler.New()

		prog, err := c.Compile(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
  				COLLECT ageGroup = FLOOR(u.age / 5) * 5 
  				AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
  				RETURN {
					ageGroup,
    				minAge, 
    				maxAge 
  				}
		`)

		So(err, ShouldBeNil)
		So(prog, ShouldHaveSameTypeAs, &runtime.Program{})

		out, err := prog.Run(context.Background())

		So(err, ShouldBeNil)
		So(string(out), ShouldEqual, `[{"ageGroup":25,"maxAge":25,"minAge":25},{"ageGroup":30,"maxAge":31,"minAge":31},{"ageGroup":35,"maxAge":36,"minAge":36},{"ageGroup":45,"maxAge":45,"minAge":45},{"ageGroup":65,"maxAge":69,"minAge":69}]`)
	})
}

func BenchmarkAggregate(b *testing.B) {
	p := compiler.New().MustCompile(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
  				COLLECT AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
  				RETURN {
    				minAge, 
    				maxAge 
  				}
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkAggregate2(b *testing.B) {
	p := compiler.New().MustCompile(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
  				COLLECT AGGREGATE ages = UNION(u.age, u.age)
  				RETURN ages
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}

func BenchmarkAggregate3(b *testing.B) {
	p := compiler.New().MustCompile(`
			LET users = [
				{
					active: true,
					married: true,
					age: 31,
					gender: "m"
				},
				{
					active: true,
					married: false,
					age: 25,
					gender: "f"
				},
				{
					active: true,
					married: false,
					age: 36,
					gender: "m"
				},
				{
					active: false,
					married: true,
					age: 69,
					gender: "m"
				},
				{
					active: true,
					married: true,
					age: 45,
					gender: "f"
				}
			]
			FOR u IN users
  				COLLECT ageGroup = FLOOR(u.age / 5) * 5 
  				AGGREGATE minAge = MIN(u.age), maxAge = MAX(u.age)
  				RETURN {
					ageGroup,
    				minAge, 
    				maxAge 
  				}
		`)

	for n := 0; n < b.N; n++ {
		p.Run(context.Background())
	}
}
