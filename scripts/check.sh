#!/bin/bash

./nacos-service-discovery-controller \
    online \
    --nacosIpAddr 10.251.100.175 \
    --nacosNamespaceId 103c333a-7141-4927-b47d-679128f8a130 \
    --serviceName flm-flc-cg-mservice \
    --serviceIp 10.110.6.247 \
    --waitTime 12m


./nacos-service-discovery-controller \
    offline \
    --nacosIpAddr 10.251.100.175 \
    --nacosNamespaceId 103c333a-7141-4927-b47d-679128f8a130 \
    --serviceName flm-flc-cg-mservice \
    --serviceIp 10.110.6.247 \
    --waitTime 15s

./nacos-service-discovery-controller \
    exporter \
    --nacosIpAddr 10.251.100.175 \
    --nacosUsername nacos \
    --nacosPassword nacos

## k8s 环境中执行命令
# ./nacos-service-discovery-controller offline --nacosNamespaceId 225203d6-f6c9-4da0-ae07-c44bfeb9be8c --serviceName flm-flc-cg-mservice --waitTime 60s
