package websockets

import (
	"container/list"
	"sync"
)

type Observer interface {
	Notify()
}

type Subject interface {
	AddObservers(observers ...Observer)
	NotifyObservers()
}

type broadcastSubject struct {
	observers map[Observer]bool
}

func (s *broadcastSubject) AddObservers(observers ...Observer) {
	for k := range observers {
		observer := observers[k]
		s.observers[observer] = true
	}
}

func (s *broadcastSubject) NotifyObservers() {
	for k := range s.observers {
		k.Notify()
	}
}

func NewBroadcastSubject() *broadcastSubject {
	return &broadcastSubject{
		observers: map[Observer]bool{},
	}
}

type unicastSubject struct {
	observers *list.List
	mutex     *sync.Mutex
}

func (u *unicastSubject) AddObservers(observers ...Observer) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	for k := range observers {
		observer := observers[k]
		u.observers.PushBack(observer)
	}
}

func (u *unicastSubject) NotifyObservers() {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if u.observers.Len() == 0 {
		return
	}

	element := u.observers.Front()
	if observer, ok := element.Value.(Observer); ok {
		observer.Notify()
	}

	u.observers.MoveToBack(element)
}

func NewUnicastSubject() *unicastSubject {
	return &unicastSubject{
		observers: list.New(),
		mutex:     &sync.Mutex{},
	}
}
