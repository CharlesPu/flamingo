package perror

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type (
	pError struct {
		stackInfo []string // 调用栈信息，由底层->顶层
		errCode   ErrCode
		err       error
	}
)

const (
	maxStackNum = 10
	pErrorTag   = "[PError]"
	sepNewLine  = "\n"

	NoAddMsg        = ""
	sepAfterObj     = ": "
	sepBetweenStack = " => "
)

func New(format string, a ...interface{}) PError {
	pErr := &pError{
		stackInfo: make([]string, 0, maxStackNum),
		err:       fmt.Errorf(format, a...),
	}
	object := getCaller()
	pErr.appendStackInfo(formatErrStr(object, format, a...))

	return pErr
}

func Wrap(err error, addMsgWithFormat string, a ...interface{}) PError {
	if err == nil {
		return nil
	}
	pErr, ok := err.(*pError)
	if !ok || pErr == nil {
		return New(addMsgWithFormat, a...)
	}
	object := getCaller()
	pErr.appendStackInfo(formatErrStr(object, addMsgWithFormat, a...))

	return pErr
}

func getCaller() string {
	pc, _, _, _ := runtime.Caller(2) // 往上找几层函数
	fc := runtime.FuncForPC(pc).Name()
	fcStrs := strings.Split(fc, "/")
	object := fcStrs[len(fcStrs)-1]
	return object
}

func (pe *pError) appendStackInfo(curInfo string) {
	if len(pe.stackInfo) < maxStackNum {
		pe.stackInfo = append(pe.stackInfo, curInfo)
	}
}

func formatErrStr(object, format string, a ...interface{}) string {
	var errStrBuffer bytes.Buffer

	errStrBuffer.WriteString(object)
	errStrBuffer.WriteString(sepAfterObj)

	if format != NoAddMsg {
		if len(a) > 0 {
			errStrBuffer.WriteString(fmt.Sprintf(format, a...))
		} else {
			errStrBuffer.WriteString(format)
		}
	}
	return errStrBuffer.String()
}

func (pe *pError) Error() string {
	if pe == nil {
		return ""
	}
	var buffer bytes.Buffer

	buffer.WriteString(pErrorTag)
	stackLen := len(pe.stackInfo)
	for idx := range pe.stackInfo {
		buffer.WriteString(pe.stackInfo[stackLen-idx-1])
		if idx < stackLen-1 {
			buffer.WriteString(sepBetweenStack)
		}
	}

	// buffer.WriteString(pe.err.Error())
	buffer.WriteString(sepNewLine)

	return buffer.String()
}

func (pe *pError) OriginError() error {
	if pe == nil {
		return nil
	}
	return pe.err
}
