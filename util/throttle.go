package util

import (
	"sync"
	"time"
)

type Throttler interface {
	Do(f func(), b bool)
}

func NewThrottler(duration time.Duration) Throttler {
	return &throttle{
		duration: duration,
	}
}

type throttle struct {
	duration time.Duration
	once     sync.Once
	m        sync.Mutex
}

// Do is Throttler implement
func (t *throttle) Do(f func(), isafter bool) {
	t.m.Lock()
	defer t.m.Unlock()
	t.once.Do(func() {
		go func() {
			time.Sleep(t.duration)
			t.m.Lock()
			defer t.m.Unlock()
			if isafter {
				f()
			}
			t.once = sync.Once{}
		}()
		if !isafter {
			f()
		}
	})
}
