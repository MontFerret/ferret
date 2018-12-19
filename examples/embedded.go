package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
	"os"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
)

type Topic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func main() {
	topics, err := getTopTenTrendingTopics()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, topic := range topics {
		fmt.Println(fmt.Sprintf("%s: %s %s", topic.Name, topic.Description, topic.URL))
	}
}

func getTopTenTrendingTopics() ([]*Topic, error) {
	query := `
		LET doc = DOCUMENT("https://github.com/topics")

		FOR el IN ELEMENTS(doc, ".py-4.border-bottom")
			LIMIT 10
			LET url = ELEMENT(el, "a")
			LET name = ELEMENT(el, ".f3")
			LET desc = ELEMENT(el, ".f5")

			RETURN {
				name: TRIM(name.innerText),
				description: TRIM(desc.innerText),
				url: "https://github.com" + url.attributes.href
			}
	`

	comp := compiler.New()

	program, err := comp.Compile(query)

	if err != nil {
		return nil, err
	}

	// create a root context
	ctx := context.Background()

	// enable HTML drivers
	// by default, Ferret Runtime knows nothing about HTML drivers
	// all HTML manipulations are done via functions from standard library
	ctx = drivers.WithDynamic(ctx, cdp.NewDriver())
	ctx = drivers.WithStatic(ctx, http.NewDriver())

	out, err := program.Run(ctx)

	if err != nil {
		return nil, err
	}

	res := make([]*Topic, 0, 10)

	err = json.Unmarshal(out, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
