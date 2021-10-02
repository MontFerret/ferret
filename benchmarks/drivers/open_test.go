package drivers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/MontFerret/ferret"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
)

func Benchmark_Open_CDP(b *testing.B) {
	ctx := context.Background()
	d := cdp.NewDriver()

	c := ferret.New()
	if err := c.Drivers().Register(d); err != nil {
		b.Fatal(err)
	}

	p, err := c.Compile(fmt.Sprintf(`
	LET doc = DOCUMENT("%s")

	RETURN TRUE
`, "https://www.montferret.dev"))

	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		if _, err := p.Run(ctx); err != nil {
			b.Fatal(err)
		}
	}
}
