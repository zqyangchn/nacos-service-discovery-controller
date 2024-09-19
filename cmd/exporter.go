package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"nacos-service-discovery-controller/pkg/logger"
	"nacos-service-discovery-controller/pkg/nacos"
	"nacos-service-discovery-controller/pkg/nacos-service-exporter"
	"nacos-service-discovery-controller/pkg/scrapes-exporter"
	"nacos-service-discovery-controller/routers"
)

var exporterCmd = &cobra.Command{
	Use:   "exporter",
	Short: "prometheus exporter",
	Run: func(cmd *cobra.Command, args []string) {
		exporterCommand()
	},
}

var (
	ServerRunMode         string
	ServerHttpPort        string
	ServerReadTimeout     time.Duration
	ServerWriteTimeout    time.Duration
	ServerShutdownTimeout time.Duration
)

func init() {
	logger.Setup()

	exporterCmd.Flags().StringVarP(&NacosIPAddr, "nacosIpAddr", "", "nacos-cs", "nacos ip address")
	exporterCmd.Flags().StringVarP(&NacosScheme, "nacosScheme", "", "http", "nacos schema, http or https")
	exporterCmd.Flags().Uint64VarP(&NacosPort, "nacosPort", "", 8848, "nacos port")
	exporterCmd.Flags().StringVarP(&NacosContextPath, "nacosContextPath", "", "/nacos", "nacos context path")
	exporterCmd.Flags().StringVarP(&NacosUsername, "nacosUsername", "", "nacos", "nacos user")
	exporterCmd.Flags().StringVarP(&NacosPassword, "nacosPassword", "", "nacos", "nacos password")
	exporterCmd.Flags().StringVarP(&NacosNamespaceId, "nacosNamespaceId", "", "", "nacos namespace id")

	exporterCmd.Flags().StringVarP(&ServerRunMode, "serverRunMode", "", "release", "server run mode")
	exporterCmd.Flags().StringVarP(&ServerHttpPort, "serverHttpPort", "", "8428", "server http port")
	exporterCmd.Flags().DurationVarP(&ServerReadTimeout, "serverReadTimeout", "", 1*time.Minute, "server read timeout")
	exporterCmd.Flags().DurationVarP(&ServerWriteTimeout, "serverWriteTimeout", "", 1*time.Minute, "server write timeout")
	exporterCmd.Flags().DurationVarP(&ServerShutdownTimeout, "serverShutdownTimeout", "", 1*time.Minute, "server shutdown timeout")

	rootCmd.AddCommand(exporterCmd)
}

// @title API swagger
// @version 1.0
// @description message channel.
// @termsOfService http://github.com/zqyangchn

// @contact.name API Support
// @contact.url http://github.com/zqyangchn
// @contact.email zqyangchn@gmail.com
func exporterCommand() {
	// init scrapes prometheus exporter client
	prometheus.MustRegister(scrapesexp.New())
	logger.Info("scrapes prometheus exporter 启动 ...")

	// init nacos service status prometheus exporter client
	nacosConfig := nacos.NewConfig().SetScheme(NacosScheme).SetIpAddr(
		NacosIPAddr).SetPort(NacosPort).SetContextPath(NacosContextPath).SetUsername(
		NacosUsername).SetPassword(NacosPassword).SetNamespaceId(NacosNamespaceId)
	prometheus.MustRegister(nacos_service_exporter.New(nacosConfig))
	logger.Info("nacos service discovery prometheus exporter 启动 ...")

	gin.SetMode(ServerRunMode)
	router := routers.NewRouter()
	srv := &http.Server{
		Addr:           ":" + ServerHttpPort,
		Handler:        router,
		ReadTimeout:    ServerReadTimeout,
		WriteTimeout:   ServerWriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		logger.Info("web server 启动 ...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("web server 启动失败 !", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 15 秒的超时时间）
	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt)
	<-osSignal

	// 启动服务器关闭流程
	logger.Info("关闭 web server ...")
	ctx, cancel := context.WithTimeout(context.Background(), ServerShutdownTimeout)
	defer cancel()
	// srv.Shutdown(ctx) 关闭服务器监听端口, 不再接受新的请求
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("web server 关闭失败", zap.Error(err))
	}
	logger.Info("web server 关闭 !")
}
