package cli

import (
	"strings"

	"github.com/derekparker/trie"
)

// AutoCompleter autocompletes queries
// into the REPL.
// Implements AutoCompleter interface from
// github.com/chzyer/readline
type AutoCompleter struct {
	coreFuncs *trie.Trie
}

func NewAutoCompleter(functions []string) *AutoCompleter {
	coreFuncs := trie.New()

	for _, function := range functions {
		coreFuncs.Add(function, function)
	}

	return &AutoCompleter{
		coreFuncs: coreFuncs,
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

	for _, fn := range ac.coreFuncs.PrefixSearch(token) {
		// cuts a piece of word that is already written
		// in the repl
		withoutPre := []rune(fn)[len(token):]
		newLine = append(newLine, withoutPre)
	}

	return newLine, pos
}
