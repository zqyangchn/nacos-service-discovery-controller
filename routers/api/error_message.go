package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nacos-service-discovery-controller/pkg/app"
	"nacos-service-discovery-controller/service"
)

// ErrorMessages
// @Summary 获取错误码
// @ID ErrorMessages
// @Accept json
// @Produce json
// @Success 200 {object} service.ErrorMessageResponse
// @Router /error/message [get]
func ErrorMessages(c *gin.Context) {
	appG := app.Gin{Context: c}
	appG.ResponseSuccess(http.StatusOK, service.GetAllErrorMessage())
}
