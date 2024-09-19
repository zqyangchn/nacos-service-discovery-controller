package errcode

import (
	"fmt"
	"sort"
)

type ErrorMessage struct {
	Code    string
	Message string
	Details []string
}

var codeSet = make(map[string]ErrorMessage, 100)

// New 新建错误
func New(errorCode string, errorMessage string) *ErrorMessage {
	if _, ok := codeSet[errorCode]; ok {
		msg := fmt.Sprintf("errorCode exist: %s, use other\n", errorCode)
		panic(msg)
	}

	eMsg := &ErrorMessage{
		Code:    errorCode,
		Message: errorMessage,
	}
	codeSet[errorCode] = *eMsg

	return eMsg
}

// Error 接口实现
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("error code: %d, error message :%s\n", e.Code, e.Message)
}

// WithDetails 添加错误详细描述信息
func (e *ErrorMessage) WithDetails(details ...string) *ErrorMessage {
	eMsg := *e
	eMsg.Details = []string{}
	for _, d := range details {
		eMsg.Details = append(eMsg.Details, d)
	}
	return &eMsg
}

type ErrorMessages []ErrorMessage

// Len 实现 sort 接口
func (e ErrorMessages) Len() int {
	return len(e)
}
func (e ErrorMessages) Less(i, j int) bool {
	return e[i].Code < e[j].Code
}
func (e ErrorMessages) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func GetAllErrorMessage() ErrorMessages {
	e := make([]ErrorMessage, 0, len(codeSet))

	for _, eMsg := range codeSet {
		e = append(e, eMsg)
	}

	sort.Sort(ErrorMessages(e))

	return e
}
