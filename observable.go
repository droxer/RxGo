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
			var st Subscriber
			defer func() {
				if e := recover(); e != nil {
					st.OnError(e.(error))
				}
			}()
			st = op.Call(sub)
			st.Start()
			o.onSubscribe(st)
		},
	}
}
