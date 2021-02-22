package workerpool

import "time"

type (
	workerPoolOption struct {
		cleanDuration time.Duration // 0 = disable
	}

	OptionFunc func(*workerPoolOption)
)

const (
	defaultIdleDuration = time.Minute * 10
)

var (
	defaultOptions = &workerPoolOption{
		cleanDuration: defaultIdleDuration,
	}
)

func WithCleanDuration(d time.Duration) OptionFunc {
	return func(option *workerPoolOption) {
		option.cleanDuration = d
	}
}
