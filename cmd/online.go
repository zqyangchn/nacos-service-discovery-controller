package cmd

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"nacos-service-discovery-controller/pkg/basic-utils"
	"nacos-service-discovery-controller/pkg/logger"
	"nacos-service-discovery-controller/pkg/nacos"
)

const (
	nacosOnlineOperationTimeLimit = 15 * time.Second
)

var OnlineWaitTime time.Duration

var onlineCmd = &cobra.Command{
	Use:   "online",
	Short: "上线 nacos 注册的微服务",
	Run: func(cmd *cobra.Command, args []string) {
		onlineCommand()
	},
}

func init() {
	logger.Setup()

	onlineCmd.Flags().StringVarP(&NacosIPAddr, "nacosIpAddr", "", "nacos-cs", "nacos ip address")
	onlineCmd.Flags().StringVarP(&NacosScheme, "nacosScheme", "", "http", "nacos schema, http or https")
	onlineCmd.Flags().Uint64VarP(&NacosPort, "nacosPort", "", 8848, "nacos port")
	onlineCmd.Flags().StringVarP(&NacosContextPath, "nacosContextPath", "", "/nacos", "nacos context path")
	onlineCmd.Flags().StringVarP(&NacosUsername, "nacosUsername", "", "nacos", "nacos user")
	onlineCmd.Flags().StringVarP(&NacosPassword, "nacosPassword", "", "nacos", "nacos password")
	onlineCmd.Flags().StringVarP(&NacosNamespaceId, "nacosNamespaceId", "", "", "nacos namespace id")

	onlineCmd.Flags().StringVarP(&ServiceIp, "serviceIp", "i", "", "service ip")
	onlineCmd.Flags().Uint64VarP(&ServicePort, "servicePort", "p", 8080, "service port")
	onlineCmd.Flags().StringVarP(&ServiceName, "serviceName", "n", "", "service name")
	onlineCmd.Flags().StringVarP(&ServiceClusterName, "serviceClusterName", "c", "DEFAULT", "service cluster name")
	onlineCmd.Flags().StringVarP(&ServiceGroupName, "serviceGroupName", "g", "DEFAULT_GROUP", "service group name")

	onlineCmd.Flags().DurationVarP(&OnlineWaitTime, "waitTime", "w", 10*time.Minute, "wait time")

	rootCmd.AddCommand(onlineCmd)
}

func onlineCommand() {
	if NacosNamespaceId == "" || ServiceName == "" {
		logger.Error("NacosNamespaceId 或 ServiceName 空",
			zap.String("NacosNamespaceId", NacosNamespaceId), zap.String("ServiceName", ServiceName),
		)
		os.Exit(110)
	}

	if ServiceIp == "" {
		ip, err := basicutils.GetLocalIp()
		if err != nil {
			logger.Error("获取 ServiceIp 失败: %s", zap.Error(err))
			os.Exit(111)
		}
		ServiceIp = ip
	}

	totalTimeout := OnlineWaitTime + nacosOnlineOperationTimeLimit
	ctx, cancel := context.WithTimeout(context.Background(), totalTimeout)
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			logger.Error("脚本总执行时间超时, 超过上线时间阈值",
				zap.Duration("totalTimeout", totalTimeout),
				zap.Error(ctx.Err()),
			)
			os.Exit(112)
		}
	}()

	logger.Info("等待服务启动", zap.Duration("超时时间", OnlineWaitTime))
	if !waitingForServiceHealth() {
		logger.Info("等待服务启动超时 ")
		os.Exit(113)
	}
	logger.Info("服务启动成功 ")

	logger.Info("开始上线微服务")
	if err := onlineMicroservices(); err != nil {
		logger.Error("上线微服务失败", zap.Error(err))
		os.Exit(114)
	}
	logger.Info("上线微服务成功")
}

func onlineMicroservices() error {
	nacosConfig := nacos.NewConfig().SetScheme(NacosScheme).SetIpAddr(NacosIPAddr).SetPort(NacosPort).SetContextPath(NacosContextPath).SetUsername(NacosUsername).SetPassword(NacosPassword).SetNamespaceId(NacosNamespaceId)
	nacosClient, err := nacos.New(nacosConfig)
	if err != nil {
		return err
	}

	instance, err := nacosClient.RetryGetInstance(
		getInstanceRetryCount, getInstanceRetryInterval,
		nacos.GetInstanceParam{
			ServiceName: ServiceName,
			Ip:          ServiceIp,
			Port:        ServicePort,
			NamespaceId: NacosNamespaceId,
			GroupName:   ServiceGroupName,
			ClusterName: ServiceClusterName,
			HealthyOnly: true,
			Ephemeral:   true,
		})
	if err != nil {
		return err
	}

	metadataType, err := json.Marshal(instance.Metadata)
	if err != nil {
		return err
	}
	metadataString := string(metadataType)
	if err := nacosClient.UpdateInstance(nacos.UpdateInstanceParam{
		Enable: true,

		ServiceName: ServiceName,
		Ip:          instance.IP,
		Port:        instance.Port,
		NamespaceId: NacosNamespaceId,

		Ephemeral:   true,
		Weight:      instance.Weight,
		Healthy:     instance.Healthy,
		Metadata:    metadataString,
		GroupName:   ServiceGroupName,
		ClusterName: ServiceClusterName,
	}); err != nil {
		return err
	}

	return nil
}

func waitingForServiceHealth() bool {
	uri := "http://" + ServiceIp + ":" + strconv.FormatUint(ServicePort, 10) + "/health"

	deadline := time.After(OnlineWaitTime)
	ticker := time.Tick(1 * time.Second)

	for {
		select {
		case <-deadline:
			return false
		case <-ticker:
			if ok, _ := basicutils.CheckURL(uri); ok {
				return true
			}
		}
	}
}
