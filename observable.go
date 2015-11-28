package RxGo

type Observable struct {
	on OnSubscribe
}

func create(f OnSubscribe) *Observable {
	return &Observable{f}
}

func (o *Observable) Subscribe(sub Subscriber) {
	o.on.Call(sub)
}

type OnSubscribe interface {
	Call(sub Subscriber)
}
