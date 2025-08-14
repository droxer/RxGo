package streams

// BackpressureStrategy defines how to handle buffer overflow scenarios
type BackpressureStrategy int

const (
	// Buffer keeps all items in a bounded buffer, blocking when full
	Buffer BackpressureStrategy = iota

	// Drop discards new items when buffer is full
	Drop

	// Latest keeps only the latest item, discarding older ones
	Latest

	// Error signals an error when buffer overflows
	Error
)

// String returns string representation of backpressure strategy
func (bs BackpressureStrategy) String() string {
	switch bs {
	case Buffer:
		return "Buffer"
	case Drop:
		return "Drop"
	case Latest:
		return "Latest"
	case Error:
		return "Error"
	default:
		return "Unknown"
	}
}

// BackpressureConfig configures backpressure behavior for publishers
type BackpressureConfig struct {
	Strategy   BackpressureStrategy
	BufferSize int64
}
