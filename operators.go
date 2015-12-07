package RxGo

import (
	q "github.com/droxer/queue"
)

type Operator interface {
	Call(sub Subscriber) Subscriber
}

type opObserveOn struct {
	scheduler Scheduler
}

type observeOnSubscriber struct {
	worker Worker
	child  Subscriber
	queue  q.Queue
}

func (op *opObserveOn) Call(sub Subscriber) Subscriber {
	return &observeOnSubscriber{
		worker: op.scheduler.CreateWorker(),
		child:  sub,
		queue:  q.New(),
	}
}

func (o *observeOnSubscriber) OnNext(next interface{}) {
	o.queue.Add(next)

	o.worker.Schedule(func() {
		for {
			item := o.queue.Poll()
			if item == nil {
				return
			}
			o.child.OnNext(item)
		}
	})
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
