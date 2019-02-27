package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
)

func run(q string) ([]byte, error) {
	comp := compiler.New()
	program := comp.MustCompile(q)

	// create a root context
	ctx := context.Background()

	// we inform the driver to keep cookies between queries
	ctx = drivers.WithContext(
		ctx,
		cdp.NewDriver(cdp.WithKeepCookies()),
		drivers.AsDefault(),
	)

	return program.Run(ctx)
}