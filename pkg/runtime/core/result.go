package core

type (
	Result struct {
		err  <-chan error
		data <-chan Value
	}
)

func NewResult(err <-chan error, data <-chan Value) *Result {
	return &Result{err, data}
}

func (r *Result) Error() <-chan error {
	return r.err
}

func (r *Result) Data() <-chan Value {
	return r.data
}
