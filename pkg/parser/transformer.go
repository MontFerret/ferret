package parser

import "github.com/antlr4-go/antlr/v4"

type TokenStreamTransformer func(antlr.TokenStream) antlr.TokenStream
