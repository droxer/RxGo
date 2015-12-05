package RxGo

type opSubscribeOn struct {
	scheduler Scheduler
}

func (op *opSubscribeOn) Call(sub Subscriber) *Subscriber {
	inner := op.scheduler.CreateWorker()
	sub.Add(inner)

	return &Subscriber{
		OnError: func(e error) {
			sub.OnError(e)
		},
		OnNext: func(next interface{}) {
			inner.Schedule(func() {

			})
		},
	}
}
