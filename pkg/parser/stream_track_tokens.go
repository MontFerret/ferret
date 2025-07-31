package parser

import "github.com/antlr4-go/antlr/v4"

type TrackingTokenStream struct {
	antlr.TokenStream
	history *TokenHistory
}

func NewTrackingTokenStream(stream antlr.TokenStream, history *TokenHistory) antlr.TokenStream {
	return &TrackingTokenStream{
		TokenStream: stream,
		history:     history,
	}
}

func (s *TrackingTokenStream) Consume() {
	// Get current token before advancing
	tok := s.LT(1)
	s.TokenStream.Consume()
	s.history.Add(tok)
}
