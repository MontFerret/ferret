package mem

type interfaceCloser struct {
	marker any
	closes *int
}

func (c interfaceCloser) Close() error {
	*c.closes++

	return nil
}
