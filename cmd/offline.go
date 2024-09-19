package cmd

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"nacos-service-discovery-controller/pkg/basic-utils"
	"nacos-service-discovery-controller/pkg/logger"
	"nacos-service-discovery-controller/pkg/nacos"
)

const (
	nacosOfflineOperationTimeLimit = 10 * time.Second
)

var OfflineWaitTime time.Duration

var offlineCmd = &cobra.Command{
	Use:   "offline",
	Short: "下线 nacos 注册的微服务",
	Run: func(cmd *cobra.Command, args []string) {
		offlineCommand()
	},
}

func init() {
	logger.Setup()

	offlineCmd.Flags().StringVarP(&NacosIPAddr, "nacosIpAddr", "", "nacos-cs", "nacos ip address")
	offlineCmd.Flags().StringVarP(&NacosScheme, "nacosScheme", "", "http", "nacos schema, http or https")
	offlineCmd.Flags().Uint64VarP(&NacosPort, "nacosPort", "", 8848, "nacos port")
	offlineCmd.Flags().StringVarP(&NacosContextPath, "nacosContextPath", "", "/nacos", "nacos context path")
	offlineCmd.Flags().StringVarP(&NacosUsername, "nacosUsername", "", "nacos", "nacos user")
	offlineCmd.Flags().StringVarP(&NacosPassword, "nacosPassword", "", "nacos", "nacos password")
	offlineCmd.Flags().StringVarP(&NacosNamespaceId, "nacosNamespaceId", "", "", "nacos namespace id")

	offlineCmd.Flags().StringVarP(&ServiceIp, "serviceIp", "i", "", "service ip")
	offlineCmd.Flags().Uint64VarP(&ServicePort, "servicePort", "p", 8080, "service port")
	offlineCmd.Flags().StringVarP(&ServiceName, "serviceName", "n", "", "service name")
	offlineCmd.Flags().StringVarP(&ServiceClusterName, "serviceClusterName", "c", "DEFAULT", "service cluster name")
	offlineCmd.Flags().StringVarP(&ServiceGroupName, "serviceGroupName", "g", "DEFAULT_GROUP", "service group name")

	offlineCmd.Flags().DurationVarP(&OfflineWaitTime, "waitTime", "w", 45*time.Second, "wait time")

	rootCmd.AddCommand(offlineCmd)
}

func offlineCommand() {
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

	totalTimeout := OfflineWaitTime + nacosOfflineOperationTimeLimit
	ctx, cancel := context.WithTimeout(context.Background(), totalTimeout)
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
			logger.Error("脚本总执行时间超时, 下线时间超过",
				zap.Duration("totalTimeout", totalTimeout),
			)
			os.Exit(112)
		}
	}()

	logger.Info("开始下线微服务")
	if err := offlineMicroservices(); err != nil {
		logger.Error("下线微服务失败", zap.Error(err))
		os.Exit(113)
	}
	logger.Info("下线微服务成功")

	logger.Info("等待其他服务刷新微服务列表", zap.Duration("等待时间", OfflineWaitTime))
	time.Sleep(OfflineWaitTime)
	logger.Info("等待终止, 退出 ...")
}

func offlineMicroservices() error {
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

	if err := nacosClient.RetryUpdateInstance(
		updateInstanceRetryCount, updateInstanceRetryInterval,
		nacos.UpdateInstanceParam{
			Enable: false,

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
