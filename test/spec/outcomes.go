package spec

type Outcomes struct {
	Result Expectation
	Error  Expectation
}

func (o Outcomes) Merge(other Outcomes) Outcomes {
	return Outcomes{
		Result: o.Result.Merge(other.Result),
		Error:  o.Error.Merge(other.Error),
	}
}
