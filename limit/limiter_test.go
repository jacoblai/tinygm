package limit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiter_Allow(t *testing.T) {
	var limiter = NewLimiter(uint64(5))
	for i := 0; i < 10; i++ {
		t.Log(i, limiter.Allow())
	}
	for i := 0; i < 70; i++ {
		t.Log(i, limiter.Allow())
		time.Sleep(1 * time.Second)
	}
}

func BenchmarkLimiter_Allow(b *testing.B) {
	var limiter = NewLimiter(uint64(50000))
	gw := sync.WaitGroup{}
	ct := int64(0)
	for i := 0; i < b.N; i++ {
		if limiter.Allow() {
			gw.Add(1)
			atomic.AddInt64(&ct, 1)
			b.Log(i)
			gw.Done()
		}
	}
	gw.Wait()
	b.Log("----", ct)
}
