package spec

import "github.com/MontFerret/ferret/v2/pkg/vm"

type ExecInfo struct {
	Outcomes
	Env       []vm.EnvironmentOption
	VM        []vm.Option
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

	if len(other.VM) > 0 {
		out.VM = append(out.VM, other.VM...)
	}

	out.Outcomes = out.Outcomes.Merge(other.Outcomes)

	return out
}
