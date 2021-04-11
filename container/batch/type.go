package batch

type (
	BatchProcessor interface {
		Put(interface{})
		TryPut(interface{}) bool

		Run()
		Shutdown()

		Num() int

		flush([]interface{})
	}

	Consumer interface {
		// must be reentrancy
		Consume([]interface{})
	}

	KeyBatchProcessor interface {
		TryPut(string, interface{}) bool

		Run()
		Shutdown()

		KeyNum() int
		BufferNum() int

		flush(string, []interface{})
	}

	KeyConsumer interface {
		// must be reentrancy
		Consume(string, []interface{})
	}
)
