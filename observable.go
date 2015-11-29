package RxGo

type OnSubscribe func(sub Subscriber)

type Observable struct {
	on OnSubscribe
}

func Create(fn OnSubscribe) *Observable {
	return &Observable{fn}
}

func (o *Observable) Subscribe(sub Subscriber) {
	o.on(sub)
}
