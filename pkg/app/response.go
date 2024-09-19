package app

import (
	"github.com/gin-gonic/gin"

	"nacos-service-discovery-controller/pkg/errcode"
)

type Gin struct {
	Context *gin.Context
}

type ResponseError struct {
	*errcode.ErrorMessage
}

type ResponseSuccess struct {
	*errcode.ErrorMessage
}

// ResponseError response error
func (g *Gin) ResponseError(httpCode int, eMsg *errcode.ErrorMessage) {
	g.Context.JSON(httpCode, ResponseError{
		ErrorMessage: eMsg,
	})
	return
}

// ResponseSuccess response success
func (g *Gin) ResponseSuccess(httpCode int, data interface{}) {
	g.Context.JSON(httpCode, data)
	return
}
