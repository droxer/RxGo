package RxGo

func Create(on OnSubscribe) *Observable {
	return &Observable{on}
}

type OnSubscribe func(sub *Subscriber)

type Observable struct {
	on OnSubscribe
}

func (o *Observable) Subscribe(sub *Subscriber) {
	o.on(sub)
}
