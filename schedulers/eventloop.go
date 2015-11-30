package schedulers

type eventLoopScheduler struct {
}

func (e *eventLoopScheduler) CreateWorker() Worker {
	return &eventLoopWorker{}
}

type eventLoopWorker struct {
}

func (e *eventLoopWorker) UnSubscribe() {

}

func (e *eventLoopWorker) IsSubscribed() bool {
	return false
}

func (e *eventLoopWorker) Schedule(ac action0) {

}
