package stack

type Stack struct {
	q     []interface{}
	depth int
}

func (s *Stack) Init(size uint32) *Stack {
	if s == nil {
		return s
	}

	s.q = make([]interface{}, 0, size)
	s.depth = int(size)
	return s
}

func New(size uint32) *Stack {
	return new(Stack).Init(size)
}

func (s *Stack) Push(e interface{}) {
	if s.Full() {
		return
	}
	s.q = append(s.q, e)
}

func (s *Stack) Pop() interface{} {
	if s.Empty() {
		return nil
	}
	tail := s.Len() - 1
	ret := s.q[tail]
	s.q = s.q[:tail]
	return ret
}

func (s *Stack) Empty() bool {
	return len(s.q) == 0
}

func (s *Stack) Full() bool {
	return len(s.q) == s.depth
}

func (s *Stack) Len() int {
	return len(s.q)
}

func (s *Stack) Size() int {
	return s.depth
}

func (s *Stack) PushInt(e int) {
	if s.Full() {
		return
	}
	s.q = append(s.q, e)
}

func (s *Stack) PopInt() int {
	if s.Empty() {
		return 0
	}
	tail := s.Len() - 1
	ret := s.q[tail]
	s.q = s.q[:tail]
	return ret.(int)
}
