package core

type CollectorProjection struct {
	groupsVariable string
	countVariable  string
}

func NewCollectorGroupProjection(groupsVariable string) *CollectorProjection {
	return &CollectorProjection{
		groupsVariable: groupsVariable,
		countVariable:  "",
	}
}

func NewCollectorCountProjection(countVariable string) *CollectorProjection {
	return &CollectorProjection{
		groupsVariable: "",
		countVariable:  countVariable,
	}
}

func (p *CollectorProjection) VariableName() string {
	if p.groupsVariable != "" {
		return p.groupsVariable
	}

	if p.countVariable != "" {
		return p.countVariable
	}

	return ""
}

func (p *CollectorProjection) IsGrouped() bool {
	return p.groupsVariable != ""
}

func (p *CollectorProjection) IsCounted() bool {
	return p.countVariable != ""
}
