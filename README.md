# nacos-service-discovery-controller

------


## 编译
``` shell
Usage:
        ./build mac | linux
        Compile For Mac or Linux
```

## 容器
``` shell
version=0.2.1

## company
docker buildx build \
    --platform linux/amd64 \
    --push -t reg.flashexpress.com/library/nacos-service-discovery-controller:$version \
    -f Dockerfile .

## cloud
docker buildx build \
    --platform linux/amd64 \
    --push -t acr-pro-registry.ap-southeast-1.cr.aliyuncs.com/flashexpressopen/nacos-service-discovery-controller:$version \
    -f Dockerfile .

```

## kubernetes 容器生命周期回调
``` shell
需要把 nacos-service-discovery-controller 打到业务容器镜像里

k8s deployment 配置示例
lifecycle:
  #postStart:
  #  exec:
  #    command: ["/mnt/scripts/nacos-service-discovery-controller", "online",
  #      "--nacosIpAddr", "nacos-cs",
  #      "--nacosNamespaceId", "nacosNamespaceId",
  #      "--serviceName", "service name",
  #      "--waitTime", "12m"
  #    ]
  preStop:
    exec:
      command: ["/mnt/scripts/nacos-service-discovery-controller", "offline",
        "--nacosIpAddr", "nacos-cs",
        "--nacosNamespaceId", "nacosNamespaceId",
        "--serviceName", "service name",
        "--waitTime", "45s"
      ]
```

## nacos 服务监控
``` shell
部署参考 helm 目录

启动示例
./nacos-service-discovery-controller \
    exporter \
    --nacosIpAddr 10.251.100.175 \
    --nacosUsername nacos \
    --nacosPassword nacos
```
