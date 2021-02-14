package workqueue

type (
	// rate limit + ring queue
	RateLimitQueue interface {
		DelayQueue

		// enqueue an item if the rate limiter says ok
		AddRateLimited(item interface{})

		Forget(item interface{})

		NumRequeue(item interface{}) int
	}

	rateLimitQueue struct {
		DelayQueue
		limiter RateLimiter
	}
)

func NewRateLimitQueue(qSize int, limiter RateLimiter) RateLimitQueue {
	return &rateLimitQueue{
		DelayQueue: NewDelayQueue(qSize),
		limiter:    limiter,
	}
}

func (r *rateLimitQueue) AddRateLimited(item interface{}) {
	r.DelayQueue.AddAfter(item, r.limiter.When(item))
}

func (r *rateLimitQueue) Forget(item interface{}) {
	r.limiter.Forget(item)
}

func (r *rateLimitQueue) NumRequeue(item interface{}) int {
	return r.limiter.NumRequeue(item)
}
