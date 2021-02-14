package workqueue

import (
	"reflect"
	"testing"
)

func TestNewRingBuffer(t *testing.T) {
	q := NewRingQueue(5)

	type (
		args struct {
			add    []interface{}
			getNum int
		}
		want struct {
			lenBefore, lenAfter int
			get                 []interface{}
		}
	)
	tests := []struct {
		args args
		want want
	}{
		{
			want: want{},
		},
		{
			args: args{
				add: []interface{}{1},
			},
			want: want{
				lenAfter: 1,
			},
		},
		{
			args: args{
				add: []interface{}{1, 2, 3, 4, 5},
			},
			want: want{
				lenBefore: 1,
				lenAfter:  5,
			},
		},
		{
			args: args{
				getNum: 2,
			},
			want: want{
				lenBefore: 5,
				lenAfter:  3,
				get:       []interface{}{1, 1},
			},
		},
		{
			args: args{
				getNum: 4,
			},
			want: want{
				lenBefore: 3,
				lenAfter:  0,
				get:       []interface{}{2, 3, 4},
			},
		},
		{
			args: args{
				getNum: 100,
			},
			want: want{
				lenBefore: 0,
				lenAfter:  0,
			},
		},
		{
			args: args{
				add:    []interface{}{14, 15},
				getNum: 1,
			},
			want: want{
				lenBefore: 0,
				lenAfter:  1,
				get:       []interface{}{14},
			},
		},
		{
			args: args{
				add: []interface{}{11, 12, 13, 16},
			},
			want: want{
				lenBefore: 1,
				lenAfter:  5,
			},
		},
		{
			args: args{
				getNum: 5,
			},
			want: want{
				lenBefore: 5,
				lenAfter:  0,
				get:       []interface{}{15, 11, 12, 13, 16},
			},
		},
	}
	for _, tt := range tests {
		got := want{
			lenBefore: q.Len(),
		}
		for _, v := range tt.args.add {
			q.Add(v)
		}
		for i := 0; i < tt.args.getNum; i++ {
			val := q.Get()
			if val != nil {
				got.get = append(got.get, val)
			}
		}
		got.lenAfter = q.Len()

		t.Logf("ring bugger: %+v", q.(*ringQueue))
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("NewRingQueue() = %+v, want %+v", got, tt.want)
		}
	}
}
