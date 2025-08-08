package subscriber

// Observer defines the basic observer interface
type Observer[T any] interface {
	OnCompleted()
	OnError(e error)
	OnNext(next T)
}

// Subscriber defines the subscriber interface for legacy observable
// This maintains backward compatibility with the old API
type Subscriber[T any] interface {
	Observer[T]
	Start()
}
