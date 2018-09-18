package app

import (
	"fmt"
	"time"
)

type Timer struct {
	start    time.Time
	duration time.Duration
}

func NewTimer() *Timer {
	return &Timer{}
}

func (t *Timer) Start() {
	t.start = time.Now()
}

func (t *Timer) Stop() {
	t.duration = time.Since(t.start)
}

func (t *Timer) Print() string {
	return fmt.Sprintf("%f seconds", t.duration.Seconds())
}
