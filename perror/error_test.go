package perror

import "testing"

func TestError(t *testing.T) {
	t.Log(tt1())
}

func tt1() PError {
	return Wrap(tt2(), "haha")
}

func tt2() PError {
	return New("hehe")
}
