package encoding

import "github.com/MontFerret/api/result"

type (
	// Output is a caller-owned, one-shot encoded output handle. Consume or Collect
	// observes its terminal outcome; Close abandons it and reports cleanup errors.
	Output = result.Output
	// Content is detached encoded data owned by its recipient.
	Content = result.Content
	// Metadata is the immutable descriptor of the complete encoded representation.
	Metadata = result.Metadata
	// Consumer receives synchronous, borrowed, read-only encoded chunks.
	Consumer = result.Consumer
)

var (
	// ErrOutputInUse identifies a competing consumption attempt.
	ErrOutputInUse = result.ErrInUse
	// ErrOutputClosed identifies consumption after stopping, or interrupted by Close.
	ErrOutputClosed = result.ErrClosed
)
