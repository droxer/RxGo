package RxGo

type OnSubscribe func(sub Subscriber)

type Observable struct {
	on OnSubscribe
}

func Create(on OnSubscribe) *Observable {
	return &Observable{on}
}

func (o *Observable) Subscribe(sub Subscriber) {
	o.on(sub)
}
