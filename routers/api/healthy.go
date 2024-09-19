package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nacos-service-discovery-controller/pkg/app"
	"nacos-service-discovery-controller/service/healthy"
)

// Ready
// @Summary Ready api for kubernetes readinessProbe
// @ID Ready
// @Accept json
// @Produce json
// @Success 200 {object} healthy.Response
// @Router /ready [get]
func Ready(c *gin.Context) {
	appG := app.Gin{Context: c}
	appG.ResponseSuccess(http.StatusOK, healthy.New())
}

// Healthy
// @Summary Healthy api for kubernetes readinessProbe
// @ID Healthy
// @Accept json
// @Produce json
// @Success 200 {object} healthy.Response
// @Router /healthy [get]
func Healthy(c *gin.Context) {
	appG := app.Gin{Context: c}
	appG.ResponseSuccess(http.StatusOK, healthy.New())
}
