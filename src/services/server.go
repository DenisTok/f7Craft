package services

import (
	"errors"
	"sync"
	"time"
)

type ServerService interface {
	AddToQuarry(name string)
	CheckAccess(name string) error
	GiveAccess(name string) error
}

type serverService struct {
	logged map[string]bool // [minecraft name] yes\no
	sync.RWMutex
}

func NewServerService() ServerService {
	return &serverService{
		logged: make(map[string]bool),
	}
}

func (ss *serverService) AddToQuarry(name string) {
	ss.Lock()
	ss.logged[name] = false
	ss.Unlock()
}

func (ss *serverService) CheckAccess(name string) error {
	ss.RLock()
	b, ok := ss.logged[name]
	ss.RUnlock()
	if !ok {
		return errors.New("no user")
	}

	if !b {
		return ErrWait
	}

	go func() {
		time.Sleep(time.Second * 2)
		ss.Lock()
		delete(ss.logged, name)
		ss.Unlock()
	}()

	return nil
}

func (ss *serverService) GiveAccess(name string) error {
	ss.RLock()
	b, ok := ss.logged[name]
	ss.RUnlock()
	if !ok {
		return errors.New("no user")
	}

	if !b {
		ss.Lock()
		ss.logged[name] = true
		ss.Unlock()
	}

	return nil
}

var ErrWait = errors.New("wait for access")
