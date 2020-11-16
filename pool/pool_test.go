package pool

import (
	"testing"
)

func TestPoolManager(t *testing.T) {
	p := NewPoolManager(2, func(i interface{}) error {
		// for {
		// 	fmt.Println(i)
		// 	time.Sleep(time.Second)
		// }
		return nil
	})

	p.Do(2)
	p.Do(2)

	p.Destroy()
	for {

	}
}
