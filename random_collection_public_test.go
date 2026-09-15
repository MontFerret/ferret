package ferret_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/rnd"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestRandomFQLCollectionContracts(t *testing.T) {
	eng, err := ferret.New()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = eng.Close() })
	trivial, err := json.Marshal([]any{nil, "only", []any{}, []string{"only"}, rnd.NewSeed(0).Float64()})
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		query string
		want  string
	}{
		{
			query: `RETURN [random::choice([]), random::choice(["only"]), random::shuffle([]), random::shuffle(["only"]), random::float()]`,
			want:  string(trivial),
		},
		{
			query: `LET original = [1, 2, 3, 4] LET shuffled = random::shuffle(original) RETURN [original, shuffled]`,
			want:  `[[1,2,3,4],[1,2,4,3]]`,
		},
	} {
		out, err := eng.Run(t.Context(), ferret.NewAnonymousSource(test.query), ferret.WithSessionRandomSeed(0))
		if err != nil || out == nil || string(out.Content) != test.want {
			t.Fatalf("%s: output %v, error %v; want %s", test.query, out, err, test.want)
		}
	}

	for _, name := range []string{"choice", "shuffle"} {
		for _, input := range []string{"none", "1", `"text"`, "{}"} {
			query := "RETURN random::" + name + "(" + input + ")"
			_, err := eng.Run(t.Context(), ferret.NewAnonymousSource(query), ferret.WithSessionRandomSeed(0))
			if !errors.Is(err, runtime.ErrInvalidType) {
				t.Fatalf("%s: error %v, want invalid type", query, err)
			}
		}
	}
}
