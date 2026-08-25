package ferret

type sessionOptionTarget struct {
	setters []nativeSessionOption
}

func newSessionOptionTarget(capacity int) *sessionOptionTarget {
	return &sessionOptionTarget{
		setters: make([]nativeSessionOption, 0, capacity),
	}
}

func (t *sessionOptionTarget) SetParam(name string, value any) error {
	t.setters = append(t.setters, sessionParamOption(name, value))

	return nil
}

func (t *sessionOptionTarget) SetParams(params map[string]any) error {
	t.setters = append(t.setters, sessionParamsOption(params))

	return nil
}

func (t *sessionOptionTarget) SetOutputContentType(contentType string) error {
	t.setters = append(t.setters, outputContentTypeOption(contentType))

	return nil
}
