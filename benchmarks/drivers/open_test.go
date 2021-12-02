package drivers_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"testing"

	"github.com/MontFerret/ferret"
)

var c *ferret.Instance

func init() {
	c = ferret.New()
	c.Drivers().Register(cdp.NewDriver(), drivers.AsDefault())
}

func Benchmark_Open_CDP(b *testing.B) {
	ctx := context.Background()

	p, err := c.Compile(`
	LET doc = DOCUMENT("https://www.montferret.dev")

	RETURN TRUE
`)

	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		if _, err := c.Run(ctx, p); err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Navigate_CDP(b *testing.B) {
	ctx := context.Background()

	p, err := c.Compile(`
LET doc = DOCUMENT('https://www.theverge.com/tech', {
    driver: "cdp",
    ignore: {
        resources: [
            {
                url: "*",
                type: "image"
            }
        ]
    }
})

WAIT_ELEMENT(doc, '.c-compact-river__entry', 5000)
LET articles = ELEMENTS(doc, '.c-entry-box--compact__image-wrapper')
LET links = (
    FOR article IN articles
        FILTER article.attributes?.href LIKE 'https://www.theverge.com/*'
        RETURN article.attributes.href
)

FOR link IN links
	LIMIT 10
    // The Verge has pretty heavy pages, so let's increase the navigation wait time
    NAVIGATE(doc, link, 20000)
    WAIT_ELEMENT(doc, '.c-entry-content', 15000)
    LET texter = ELEMENT(doc, '.c-entry-content')
    RETURN texter.innerText
`)

	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		if _, err := c.Run(ctx, p); err != nil {
			b.Fatal(err)
		}
	}
}
