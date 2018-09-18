package compiler

type visitorFn func() (interface{}, error)

type result struct {
	data interface{}
	err  error
}

func newResult(data interface{}, err error) *result {
	return &result{data, err}
}

func newResultFrom(fn visitorFn) *result {
	out, err := fn()

	return &result{out, err}
}

func (res *result) Ok() bool {
	return res.err == nil
}

func (res *result) Data() interface{} {
	return res.data
}

func (res *result) Error() error {
	return res.err
}
