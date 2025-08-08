package RxGo

type Observer[T any] interface {
	OnCompleted()
	OnError(e error)
	OnNext(next T)
}

type Subscriber[T any] interface {
	Observer[T]
	Start()
}
