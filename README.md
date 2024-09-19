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
version=0.2.0

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
