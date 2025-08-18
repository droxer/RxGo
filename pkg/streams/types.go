package streams

type BackpressureStrategy int

const (
	Buffer BackpressureStrategy = iota

	Drop

	Latest

	Error
)

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

type BackpressureConfig struct {
	Strategy   BackpressureStrategy
	BufferSize int64
}
