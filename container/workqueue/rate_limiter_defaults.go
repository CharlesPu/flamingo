package workqueue

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

type (
	exponentialBackOffRateLimiter struct {
		failLock *sync.Mutex
		failCnt  map[interface{}]int

		base time.Duration
		max  time.Duration
	}
)

// more failures, more wait time
// delay = base * (2^failCnt)
func NewExponentialBackOffRateLimiter(base time.Duration, max time.Duration) RateLimiter {
	return &exponentialBackOffRateLimiter{
		failLock: &sync.Mutex{},
		failCnt:  make(map[interface{}]int),
		base:     base,
		max:      max,
	}
}

func (e *exponentialBackOffRateLimiter) When(item interface{}) time.Duration {
	e.failLock.Lock()

	exp := e.failCnt[item]
	e.failCnt[item] = e.failCnt[item] + 1

	backoff := float64(e.base.Nanoseconds()) * math.Pow(2, float64(exp))
	if backoff > math.MaxInt64 || time.Duration(backoff) > e.max {
		e.failLock.Unlock()
		return e.max
	}
	e.failLock.Unlock()
	return time.Duration(backoff)
}

func (e *exponentialBackOffRateLimiter) Forget(item interface{}) {
	e.failLock.Lock()
	delete(e.failCnt, item)
	e.failLock.Unlock()
}

func (e *exponentialBackOffRateLimiter) NumRequeue(item interface{}) int {
	e.failLock.Lock()
	r := e.failCnt[item]
	e.failLock.Unlock()

	return r
}

type (
	randomRateLimiter struct {
		failLock *sync.Mutex
		failCnt  map[interface{}]int

		max time.Duration
	}
)

// random [0, max)
func NewRandomRateLimiter(max time.Duration) RateLimiter {
	if max >= math.MaxInt64 {
		max = math.MaxInt64 - 1
	}
	r := &randomRateLimiter{
		failLock: &sync.Mutex{},
		failCnt:  make(map[interface{}]int),
		max:      max,
	}

	return r
}
func (r *randomRateLimiter) When(item interface{}) time.Duration {
	r.failLock.Lock()
	rand.Seed(time.Now().UnixNano())
	r.failCnt[item]++

	backoff := rand.Intn(int(r.max.Nanoseconds()))
	r.failLock.Unlock()
	return time.Duration(backoff)
}

func (r *randomRateLimiter) Forget(item interface{}) {
	r.failLock.Lock()
	delete(r.failCnt, item)
	r.failLock.Unlock()
}

func (r *randomRateLimiter) NumRequeue(item interface{}) int {
	r.failLock.Lock()
	ret := r.failCnt[item]
	r.failLock.Unlock()

	return ret
}
