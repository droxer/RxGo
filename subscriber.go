package RxGo

type Observer interface {
	OnCompleted()
	OnError(e error)
	OnNext(next interface{})
}

type Subscription interface {
	UnSubscribe()
	IsSubscribed() bool
}

type Subscriber struct {
	OnNext        func(next interface{})
	OnError       func(e error)
	OnCompleted   func()
	subscriptions []Subscription
}

func (s *Subscriber) UnSubscribe() {
	// TODO: TBD
}

func (s *Subscriber) IsSubscribed() bool {
	return false
}

func (s *Subscriber) Add(sub Subscription) {
	s.subscriptions = append(s.subscriptions, sub)
}
