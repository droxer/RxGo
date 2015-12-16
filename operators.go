package RxGo

type Operator interface {
	Call(sub Subscriber) Subscriber
}

type opObserveOn struct {
	scheduler Scheduler
}

type observeOnSubscriber struct {
	scheduler Scheduler
	child     Subscriber
}

func (op *opObserveOn) Call(sub Subscriber) Subscriber {
	return &observeOnSubscriber{
		scheduler: op.scheduler,
		child:     sub,
	}
}

func (o *observeOnSubscriber) Start() {
	o.scheduler.Start()
}

func (o *observeOnSubscriber) OnNext(next interface{}) {
	o.scheduler.Schedule(func() {
		o.child.OnNext(next)
	})
}

func (o *observeOnSubscriber) OnError(e error) {

}

func (o *observeOnSubscriber) OnCompleted() {
	o.scheduler.Stop()
}

func (o *observeOnSubscriber) UnSubscribe() {

}

func (o *observeOnSubscriber) IsSubscribed() bool {
	return false
}
