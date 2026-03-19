package test

import "testing"

var zeroValueCloseCalls int

type zeroValueCloser struct{}

func (zeroValueCloser) Close() error {
	zeroValueCloseCalls++
	return nil
}

type valueCloser struct {
	closed *int
}

func (c valueCloser) Close() error {
	if c.closed != nil {
		*c.closed = *c.closed + 1
	}

	return nil
}

type pointerCloser struct {
	closed int
}

func (c *pointerCloser) Close() error {
	c.closed++
	return nil
}

func TestTestingCloseSkipsZeroValueBenchmark(t *testing.T) {
	zeroValueCloseCalls = 0
	t.Cleanup(func() {
		zeroValueCloseCalls = 0
	})

	instance := Testing[zeroValueCloser]{}
	instance.Close()

	if got := zeroValueCloseCalls; got != 0 {
		t.Fatalf("expected zero-value benchmark to be skipped, got %d closes", got)
	}

	if !isZero(instance.Benchmark) {
		t.Fatal("expected benchmark to remain reset after close")
	}
}

func TestTestingCloseClosesConfiguredValueBenchmarkAndResets(t *testing.T) {
	closed := 0
	instance := Testing[valueCloser]{}
	instance.SetBenchmark(valueCloser{closed: &closed})

	instance.Close()

	if got := closed; got != 1 {
		t.Fatalf("expected configured value benchmark to close once, got %d closes", got)
	}

	if !isZero(instance.Benchmark) {
		t.Fatal("expected value benchmark to be reset after close")
	}
}

func TestTestingCloseClosesConfiguredPointerBenchmarkAndResets(t *testing.T) {
	closer := &pointerCloser{}
	instance := Testing[*pointerCloser]{}
	instance.SetBenchmark(closer)

	instance.Close()
	instance.Close()

	if got := closer.closed; got != 1 {
		t.Fatalf("expected configured pointer benchmark to close once, got %d closes", got)
	}

	if instance.Benchmark != nil {
		t.Fatal("expected pointer benchmark to be reset after close")
	}
}
