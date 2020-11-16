package perror

type (
	// PError
	PError interface {
		Error() string      // 实现error接口，但需要自定义打印格式
		OriginError() error // 获得原始error
	}
	// ErrCode
	ErrCode uint32
)
