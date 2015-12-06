package RxGo

func Create(on OnSubscribe) *Observable {
	return &Observable{on}
}

type OnSubscribe func(sub Subscriber)

type Observable struct {
	onSubscribe OnSubscribe
}

func (o *Observable) Subscribe(sub Subscriber) {
	o.onSubscribe(sub)
}

func (o *Observable) ObserveOn(sch Scheduler) *Observable {
	return o.lift(&opObserveOn{sch})
}

func (o *Observable) lift(op Operator) *Observable {
	return &Observable{
		onSubscribe: func(sub Subscriber) {
			st := op.Call(sub)
			o.onSubscribe(st)
		},
	}
}
