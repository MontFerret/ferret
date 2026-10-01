package diagnostics

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	arrayEntryValidation struct {
		parser *fql.FqlParser
		stream *arrayEntryStream
		errors *arrayEntryErrors
	}

	// A bounded cursor over already buffered tokens. The original parser's
	// cursor and token coordinates remain untouched, including during prediction.
	arrayEntryStream struct {
		*antlr.CommonTokenStream
		eof antlr.Token
		end int
	}

	arrayEntryErrors struct {
		*antlr.DefaultErrorListener
		failed bool
	}
)

func newArrayEntryValidation(tokens *antlr.CommonTokenStream) *arrayEntryValidation {
	cursor := *tokens

	return &arrayEntryValidation{
		stream: &arrayEntryStream{CommonTokenStream: &cursor},
		errors: &arrayEntryErrors{DefaultErrorListener: antlr.NewDefaultErrorListener()},
	}
}

func (v *arrayEntryValidation) parse(start, end int) (*arrayEntryStream, bool) {
	v.stream.resetRange(start, end)
	v.errors.failed = false
	if v.parser == nil {
		v.parser = fql.NewFqlParser(v.stream)
		v.parser.BuildParseTrees = false
		v.parser.RemoveErrorListeners()
		v.parser.AddErrorListener(v.errors)
	} else {
		v.parser.SetTokenStream(v.stream)
	}

	v.parser.ArrayEntry()
	if v.errors.failed {
		// Failed rules may leave grammar-owned context flags unsettled. Only
		// reuse a parser after an entry parsed without any error recovery.
		v.parser = nil
	}

	return v.stream, !v.errors.failed
}

func (s *arrayEntryStream) resetRange(start, end int) {
	boundary := s.CommonTokenStream.Get(end)
	eof := antlr.NewCommonToken(&antlr.TokenSourceCharStreamPair{}, antlr.TokenEOF, antlr.TokenDefaultChannel, boundary.GetStart(), boundary.GetStart()-1)
	eof.SetTokenIndex(end)
	eof.SetText("<EOF>")
	s.eof = eof
	s.end = end
	s.CommonTokenStream.Seek(start)
}

func (s *arrayEntryStream) LT(k int) antlr.Token {
	token := s.CommonTokenStream.LT(k)
	if token != nil && token.GetTokenIndex() >= s.end {
		return s.eof
	}

	return token
}

func (s *arrayEntryStream) LA(k int) int {
	token := s.LT(k)
	if token == nil {
		return 0
	}

	return token.GetTokenType()
}

func (s *arrayEntryStream) Consume() {
	if s.Index() < s.end {
		s.CommonTokenStream.Consume()
	}
}

func (s *arrayEntryStream) Seek(index int) {
	s.CommonTokenStream.Seek(min(index, s.end))
}

func (s *arrayEntryStream) Get(index int) antlr.Token {
	if index >= s.end {
		return s.eof
	}

	return s.CommonTokenStream.Get(index)
}

func (s *arrayEntryStream) Size() int {
	return s.end + 1
}

func (e *arrayEntryErrors) SyntaxError(antlr.Recognizer, interface{}, int, int, string, antlr.RecognitionException) {
	e.failed = true
}
