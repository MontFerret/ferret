package spec

import "github.com/MontFerret/ferret/v2/pkg/vm"

type ExecInfo struct {
	Outcomes
	Env       []vm.EnvironmentOption
	RawOutput bool
}

func (e ExecInfo) Merge(other ExecInfo) ExecInfo {
	out := e

	if other.RawOutput {
		out.RawOutput = true
	}

	if len(other.Env) > 0 {
		out.Env = append(out.Env, other.Env...)
	}

	out.Outcomes = out.Outcomes.Merge(other.Outcomes)

	return out
}
