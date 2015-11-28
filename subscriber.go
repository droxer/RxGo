package RxGo

type Observer interface {
	OnCompleted()
	OnError(e error)
	OnNext(next interface{})
}

type Subscriber interface {
	Observer
}
