package RxGo

type Observer interface {
	OnCompleted()
	OnError(e error)
	OnNext(next interface{})
}

type Subscription interface {
	UnSubscribe()
	IsSubscribed() bool
}

type Subscriber interface {
	Observer
	Subscription
	Add(sub Subscription)
}
