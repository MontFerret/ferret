package bytecode

import (
	"fmt"

	"github.com/goccy/go-json"

	"github.com/MontFerret/ferret/v2/pkg/internal/hostfunction"
)

// HostFunction identifies one compiled host call signature.
// Its position in Functions.Host is the host binding ID.
type HostFunction struct {
	// Name is the canonical, already-qualified host function name.
	Name string `json:"name"`
	// ArgCount is the exact number of arguments at compiled callsites for this binding.
	ArgCount int `json:"argCount"`
}

// MarshalJSON writes the canonical host-function name to persisted metadata.
func (f HostFunction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string `json:"name"`
		ArgCount int    `json:"argCount"`
	}{
		Name:     hostfunction.CanonicalName(f.Name),
		ArgCount: f.ArgCount,
	})
}

// UnmarshalJSON requires explicit argument-count metadata so legacy host tables
// cannot be interpreted as zero-argument signatures.
func (f *HostFunction) UnmarshalJSON(data []byte) error {
	var decoded struct {
		ArgCount *int   `json:"argCount"`
		Name     string `json:"name"`
	}

	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	if decoded.ArgCount == nil {
		return fmt.Errorf("bytecode.HostFunction: missing argCount")
	}

	f.Name = hostfunction.CanonicalName(decoded.Name)
	f.ArgCount = *decoded.ArgCount

	return nil
}
