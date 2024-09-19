package healthy

import (
	"strconv"
	"time"

	"nacos-service-discovery-controller/pkg/errcode"
)

type Response struct {
	*errcode.ErrorMessage
	Data Result `json:"data"`
}

type Result struct {
	Timestamp string `json:"timestamp"`
}

func New() *Response {
	return &Response{
		ErrorMessage: errcode.Success,
		Data: Result{
			Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		},
	}
}
