package cli

// AutoCompleter autocompletes queries
// into the REPL.
// Implements AutoCompleter interface from
// github.com/chzyer/readline
type AutoCompleter struct{}

func NewAutoCompleter() *AutoCompleter {
	return &AutoCompleter{}
}

func (au *AutoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	if pos == 0 {
		return [][]rune{[]rune("FOR"), []rune("RETURN"), []rune("LET")}, pos
	}

	return [][]rune{[]rune("HELLO, WORLD")}, 0
}
