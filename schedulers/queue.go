package schedulers

import (
	"container/list"
	"sync"
)

type Queue interface {
	Add(item interface{}) interface{}
	Remove(item interface{}) interface{}
	Poll() interface{}
	Peek() interface{}
	Contains(item interface{}) bool
	Len() int
}

type simple struct {
	list *list.List
	sync.Mutex
}

func NewQueue() Queue {
	return &simple{
		list: list.New(),
	}
}

func (s *simple) Add(item interface{}) interface{} {
	s.Lock()
	defer s.Unlock()

	return s.list.PushFront(item)
}

func (s *simple) Poll() interface{} {
	s.Lock()
	defer s.Unlock()

	if s.list.Len() == 0 {
		return nil
	}

	item := s.list.Back()
	return s.list.Remove(item)
}

func (s *simple) Remove(item interface{}) interface{} {
	return s.list.Remove(item.(*list.Element))
}

func (s *simple) Contains(item interface{}) bool {
	for el := s.list.Front(); el != nil; el = el.Next() {
		if el.Value == item {
			return true
		}
	}
	return false
}

func (s *simple) Peek() interface{} {
	if s.list.Len() == 0 {
		return nil
	}

	return s.list.Front().Value
}

func (s *simple) Len() int {
	return s.list.Len()
}
