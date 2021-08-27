package limit

import (
	"sync/atomic"
	"time"
)

type Limiter struct {
	max    atomic.Value
	tokens atomic.Value
	last   atomic.Value
}

func NewLimiter(max uint64) *Limiter {
	lt := Limiter{}
	lt.last.Store(time.Now())
	lt.max.Store(max)
	lt.tokens.Store(uint64(0))
	return &lt
}

func (l *Limiter) Allow() bool {
	if l.tokens.Load().(uint64) < l.max.Load().(uint64) {
		l.tokens.Store(l.tokens.Load().(uint64) + 1)
		return true
	} else {
		now := time.Now()
		if now.Sub(l.last.Load().(time.Time)).Minutes() >= 1 {
			l.last.Store(now)
			l.tokens.Store(uint64(0))
			return true
		}
		return false
	}
}
