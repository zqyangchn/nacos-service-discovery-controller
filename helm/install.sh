#!/bin/bash

helm -n kube-monitor upgrade --install --wait --timeout 500s \
    nacos-service-discovery-controller ./nacos-service-discovery-controller \
    -f ./nacos-service-discovery-controller/values.yaml
