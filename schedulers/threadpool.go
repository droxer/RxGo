package schedulers

type threadPoolScheduler struct {
}

func (t *threadPoolScheduler) CreateWorker() Worker {
	return &threadWorker{}
}

type threadWorker struct {
}

func (t *threadWorker) UnSubscribe() {

}

func (t *threadWorker) IsSubscribed() bool {
	return false
}

func (t *threadWorker) Schedule(ac action0) {

}
