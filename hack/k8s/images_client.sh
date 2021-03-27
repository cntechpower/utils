#!/bin/bash
VERSION=v1.20.5
docker pull 10.0.0.2:5000/kube-apiserver-arm64:$VERSION
docker pull 10.0.0.2:5000/kube-controller-manager-arm64:$VERSION
docker pull 10.0.0.2:5000/kube-scheduler-arm64:$VERSION
docker pull 10.0.0.2:5000/kube-proxy-arm64:$VERSION
docker pull 10.0.0.2:5000/pause-arm64:3.2

docker image tag 10.0.0.2:5000/kube-apiserver-arm64:$VERSION registry.cn-hangzhou.aliyuncs.com/google_containers/kube-apiserver-arm64:$VERSION
docker image tag 10.0.0.2:5000/kube-controller-manager-arm64:$VERSION registry.cn-hangzhou.aliyuncs.com/google_containers/kube-controller-manager-arm64:$VERSION
docker image tag 10.0.0.2:5000/kube-scheduler-arm64:$VERSION registry.cn-hangzhou.aliyuncs.com/google_containers/kube-scheduler-arm64:$VERSION
docker image tag 10.0.0.2:5000/kube-proxy-arm64:$VERSION registry.cn-hangzhou.aliyuncs.com/google_containers/kube-proxy-arm64:$VERSION
docker image tag 10.0.0.2:5000/pause-arm64:3.2 registry.cn-hangzhou.aliyuncs.com/google_containers/pause-arm64:3.2


docker rmi 10.0.0.2:5000/kube-apiserver-arm64:$VERSION
docker rmi 10.0.0.2:5000/kube-controller-manager-arm64:$VERSION
docker rmi 10.0.0.2:5000/kube-scheduler-arm64:$VERSION
docker rmi 10.0.0.2:5000/kube-proxy-arm64:$VERSION
docker rmi 10.0.0.2:5000/pause-arm64:3.2
