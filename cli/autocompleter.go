package cli

import (
	"regexp"
	"strings"

	"github.com/derekparker/trie"
)

// AutoCompleter autocompletes queries
// into the REPL.
// Implements AutoCompleter interface from
// github.com/chzyer/readline
type AutoCompleter struct {
	coreFuncs  *trie.Trie
	initTokens *trie.Trie
}

// any fql script start with
// one of the follow tokens:
var initTokens = [][]rune{
	[]rune("FOR"),
	[]rune("RETURN"),
	[]rune("LET"),
}

// check that entered token is first
var initTokensRegexp = regexp.MustCompile(`^\s*[A-Z]+\s*$`)

func NewAutoCompleter(functions []string) *AutoCompleter {
	coreFuncs := trie.New()

	for _, function := range functions {
		coreFuncs.Add(function, function)
	}

	inits := trie.New()
	initStr := ""

	for _, init := range initTokens {
		initStr = string(init)
		inits.Add(initStr, initStr)
	}

	return &AutoCompleter{
		coreFuncs:  coreFuncs,
		initTokens: inits,
	}
}

// Do implements method of AutoCompleter interface
func (ac *AutoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line)
	tokens := strings.Split(lineStr, " ")
	token := tokens[len(tokens)-1]

	// if remove this check, than
	// on any empty string will return
	// all available functions
	if token == "" {
		return newLine, pos
	}

	searcher := ac.coreFuncs
	if initTokensRegexp.MatchString(lineStr) {
		searcher = ac.initTokens
	}

	for _, fn := range searcher.PrefixSearch(token) {
		// cuts a piece of word that is already written
		// in the repl
		withoutPre := []rune(fn)[len(token):]
		newLine = append(newLine, withoutPre)
	}

	return newLine, pos
}
