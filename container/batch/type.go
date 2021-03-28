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
)
