package bytecode

import (
	"fmt"

	"github.com/goccy/go-json"
)

// HostFunction identifies one compiled host call signature.
// Its position in Functions.Host is the host binding ID.
type HostFunction struct {
	// Name is the canonical, already-qualified host function name.
	Name string `json:"name"`
	// ArgCount is the exact number of arguments at compiled callsites for this binding.
	ArgCount int `json:"argCount"`
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

	f.Name = decoded.Name
	f.ArgCount = *decoded.ArgCount

	return nil
}
