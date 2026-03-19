package optimization

type pipelineTestPass struct {
	run      func(ctx *PassContext) (*PassResult, error)
	name     string
	requires []string
}

func newPipelineTestPass(name string, requires []string, run func(ctx *PassContext) (*PassResult, error)) Pass {
	return &pipelineTestPass{
		name:     name,
		requires: requires,
		run:      run,
	}
}

func (p *pipelineTestPass) Name() string {
	return p.name
}

func (p *pipelineTestPass) Requires() []string {
	return p.requires
}

func (p *pipelineTestPass) Run(ctx *PassContext) (*PassResult, error) {
	return p.run(ctx)
}
