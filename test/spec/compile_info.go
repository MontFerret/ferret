package spec

type CompileInfo struct {
	Outcomes
}

func (c CompileInfo) Merge(other CompileInfo) CompileInfo {
	return CompileInfo{
		Outcomes: c.Outcomes.Merge(other.Outcomes),
	}
}
