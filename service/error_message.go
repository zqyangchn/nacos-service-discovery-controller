package service

import (
	"nacos-service-discovery-controller/pkg/errcode"
)

// ErrorMessageResponse Response struct
type ErrorMessageResponse struct {
	*errcode.ErrorMessage
	Data errcode.ErrorMessages
}

func GetAllErrorMessage() *ErrorMessageResponse {
	return &ErrorMessageResponse{
		ErrorMessage: errcode.Success,
		Data:         errcode.GetAllErrorMessage(),
	}
}
