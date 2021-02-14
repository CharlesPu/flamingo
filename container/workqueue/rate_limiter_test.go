package workqueue

import (
	"reflect"
	"testing"
	"time"
)

func TestRateLimiter_ExponentialBackOff(t *testing.T) {
	limiter := NewExponentialBackOffRateLimiter(1*time.Millisecond, 10*time.Millisecond)

	type (
		args struct {
			item     interface{}
			isForget bool
		}
	)
	tests := []struct {
		args args

		when       time.Duration
		numRequeue int
	}{
		{
			args: args{
				item: "A",
			},
			when:       1 * time.Millisecond,
			numRequeue: 1,
		},
		{
			args: args{
				item: "A",
			},
			when:       2 * time.Millisecond,
			numRequeue: 2,
		},
		{
			args: args{
				item: "A",
			},
			when:       4 * time.Millisecond,
			numRequeue: 3,
		},
		{
			args: args{
				item: "A",
			},
			when:       8 * time.Millisecond,
			numRequeue: 4,
		},
		{
			args: args{
				item: "A",
			},
			when:       10 * time.Millisecond,
			numRequeue: 5,
		},
		{
			args: args{
				item: "B",
			},
			when:       1 * time.Millisecond,
			numRequeue: 1,
		},
		{
			args: args{
				item: "B",
			},
			when:       2 * time.Millisecond,
			numRequeue: 2,
		},
		{
			args: args{
				item:     "A",
				isForget: true,
			},
			when:       1 * time.Millisecond,
			numRequeue: 1,
		},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			if tc.args.isForget {
				limiter.Forget(tc.args.item)
			}
			w := limiter.When(tc.args.item)
			n := limiter.NumRequeue(tc.args.item)
			if !reflect.DeepEqual(tc.when, w) {
				t.Fatalf("fail When: %+v, got: %+v", tc.when, w)
			}
			if !reflect.DeepEqual(tc.numRequeue, n) {
				t.Fatalf("fail NumRequeue, expect: %+v, got: %+v", tc.numRequeue, n)
			}
		})
	}

}
