package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "nacos-service-discovery-controller/docs"
	"nacos-service-discovery-controller/middleware/auth"
	"nacos-service-discovery-controller/middleware/zaplogger"
	"nacos-service-discovery-controller/pkg/logger"
	"nacos-service-discovery-controller/routers/api"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(zaplogger.GinZap(logger.GetZapRouterLogger()))
	r.Use(zaplogger.RecoveryWithZap(logger.GetZapRouterLogger(), true))

	gin.DebugPrintRouteFunc = logger.GinDebugPrintRouteZapLoggerFunc

	// swagger
	r.GET("/8428/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ready", api.Ready)
	r.GET("/healthy", api.Healthy)

	r.GET("/metrics", api.PrometheusHandler())

	authorized := r.Group("/", auth.Auth())
	{
		authorized.GET("/error/message", api.ErrorMessages)
	}

	return r
}
