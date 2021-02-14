package workqueue

import "time"

type (
	RateLimiter interface {
		// how long item should wait
		When(item interface{}) time.Duration

		// stop tracking
		Forget(item interface{}) // todo 内存释放有问题否？调用点？

		// failure times the item has had
		NumRequeue(item interface{}) int
	}
)
