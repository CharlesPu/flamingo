package workerpoolv2

import "time"

type (
	workerPoolOption struct {
		idleDuration time.Duration // 0 = disable
	}

	OptionFunc func(*workerPoolOption)
)

const (
	defaultIdleDuration = time.Second * 5
)

var (
	defaultOptions = workerPoolOption{
		idleDuration: defaultIdleDuration,
	}
)

func WithIdleDuration(d time.Duration) OptionFunc {
	return func(option *workerPoolOption) {
		option.idleDuration = d
	}
}
