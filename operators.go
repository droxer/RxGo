package RxGo

type Operator interface {
	Call(sub Subscriber) Subscriber
}

type opObserveOn struct {
	scheduler Scheduler
}

func (op *opObserveOn) Call(sub Subscriber) Subscriber {
	return &observeOnSubscriber{
		worker: op.scheduler.CreateWorker(),
		child:  sub,
	}
}

type observeOnSubscriber struct {
	worker Worker
	child  Subscriber
}

func (o *observeOnSubscriber) OnNext(next interface{}) {

}

func (o *observeOnSubscriber) OnError(e error) {

}

func (o *observeOnSubscriber) OnCompleted() {

}

func (o *observeOnSubscriber) UnSubscribe() {

}

func (o *observeOnSubscriber) IsSubscribed() bool {
	return false
}

func (o *observeOnSubscriber) Add(sub Subscription) {

}
