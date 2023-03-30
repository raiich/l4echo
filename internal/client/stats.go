package client

import (
	"sync"
	"sync/atomic"
	"time"
)

type Stats struct {
	Name string

	mu  sync.Mutex
	min time.Duration
	max time.Duration

	dropped uint64
	delayed uint64

	sending  uint64
	received uint64
}

func (s *Stats) OnReceived(sentAt time.Time) {
	atomic.AddUint64(&s.received, 1)

	d := time.Now().Sub(sentAt)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.min > d || s.min == 0 {
		s.min = d
	}
	if s.max < d || s.max == 0 {
		s.max = d
	}
}

func (s *Stats) Dropped(count uint64) {
	atomic.AddUint64(&s.dropped, count)
}

func (s *Stats) Delayed() {
	atomic.AddUint64(&s.delayed, 1)
}

func (s *Stats) OnSending() {
	atomic.AddUint64(&s.sending, 1)
}
