package RxGo

func Create(on OnSubscriber) *Observable {
	return &Observable{on}
}

type OnSubscriber func(sub Subscriber)

type Observable struct {
	onSubscriber OnSubscriber
}

func (o *Observable) Subscribe(sub Subscriber) {
	o.onSubscriber(sub)
}

func (o *Observable) ObserveOn(sch Scheduler) *Observable {
	return o.lift(&opObserveOn{sch})
}

func (o *Observable) lift(op Operator) *Observable {
	return &Observable{
		onSubscriber: func(sub Subscriber) {
			st := op.Call(sub)
			o.onSubscriber(st)
		},
	}
}
