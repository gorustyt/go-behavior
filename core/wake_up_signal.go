package core

import (
	"sync"
	"sync/atomic"
	"time"
)

type WakeUpSignal struct {
	mutex_ *sync.Mutex
	cv_    *sync.Cond
	ready_ atomic.Bool
}

func NewWakeUpSignal() *WakeUpSignal {
	mu := &sync.Mutex{}
	return &WakeUpSignal{
		mutex_: mu,
		cv_:    sync.NewCond(mu),
	}
}

// / Return true if the timeout was NOT reached and the
// / signal was received.

func (s *WakeUpSignal) WaitFor(usec time.Duration) bool {
	timer := time.NewTimer(usec)
	defer timer.Stop()
	defer s.ready_.Store(false)
	for {
		select {
		case <-timer.C:
			return false
		default:
			if s.ready_.Load() {
				s.cv_.Wait()
				return true
			}
		}
	}
}

func (s *WakeUpSignal) EmitSignal() {
	s.ready_.Store(true)
	s.cv_.Broadcast()
}
