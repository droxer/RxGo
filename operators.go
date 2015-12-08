package RxGo

type Operator interface {
	Call(sub Subscriber) Subscriber
}

type opObserveOn struct {
	scheduler Scheduler
}

type observeOnSubscriber struct {
	worker Worker
	child  Subscriber
}

func (op *opObserveOn) Call(sub Subscriber) Subscriber {
	return &observeOnSubscriber{
		worker: op.scheduler.CreateWorker(),
		child:  sub,
	}
}

func (o *observeOnSubscriber) Start() {
	o.worker.Start()
}

func (o *observeOnSubscriber) OnNext(next interface{}) {
	o.worker.Schedule(func() {
		o.child.OnNext(next)
	})
}

func (o *observeOnSubscriber) OnError(e error) {

}

func (o *observeOnSubscriber) OnCompleted() {
	o.worker.Stop()
}

func (o *observeOnSubscriber) UnSubscribe() {

}

func (o *observeOnSubscriber) IsSubscribed() bool {
	return false
}
