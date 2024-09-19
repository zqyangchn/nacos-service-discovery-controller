package cmd

import "time"

var (
	NacosIPAddr      string
	NacosScheme      string
	NacosPort        uint64
	NacosUsername    string
	NacosPassword    string
	NacosContextPath string
	NacosNamespaceId string

	ServiceIp          string
	ServicePort        uint64
	ServiceName        string
	ServiceClusterName string
	ServiceGroupName   string
)

const (
	getInstanceRetryCount    = 3
	getInstanceRetryInterval = 3 * time.Second

	updateInstanceRetryCount    = 3
	updateInstanceRetryInterval = 3 * time.Second
)
