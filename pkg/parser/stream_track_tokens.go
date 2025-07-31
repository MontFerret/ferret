package parser

import "github.com/antlr4-go/antlr/v4"

type TrackingTokenStream struct {
	antlr.TokenStream
	tokens *TokenHistory
}

func NewTrackingTokenStream(stream antlr.TokenStream, history *TokenHistory) antlr.TokenStream {
	return &TrackingTokenStream{
		TokenStream: stream,
		tokens:      history,
	}
}

func (ts *TrackingTokenStream) Tokens() *TokenHistory {
	return ts.tokens
}

func (ts *TrackingTokenStream) LT(i int) antlr.Token {
	tok := ts.TokenStream.LT(i)

	if i == 1 && tok != nil && tok.GetTokenType() != antlr.TokenEOF {
		ts.tokens.Add(tok)
	}

	return tok
}
