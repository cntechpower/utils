#!/bin/bash
VERSION=v1.20.5
docker pull k8s.gcr.io/kube-apiserver-arm64:$VERSION
docker pull k8s.gcr.io/kube-controller-manager-arm64:$VERSION
docker pull k8s.gcr.io/kube-scheduler-arm64:$VERSION
docker pull k8s.gcr.io/kube-proxy-arm64:$VERSION
docker pull k8s.gcr.io/pause-arm64:3.2


docker image tag k8s.gcr.io/kube-apiserver-arm64:$VERSION 10.0.0.2:5000/kube-apiserver-arm64:$VERSION
docker image tag k8s.gcr.io/kube-controller-manager-arm64:$VERSION 10.0.0.2:5000/kube-controller-manager-arm64:$VERSION
docker image tag k8s.gcr.io/kube-scheduler-arm64:$VERSION 10.0.0.2:5000/kube-scheduler-arm64:$VERSION
docker image tag k8s.gcr.io/kube-proxy-arm64:$VERSION 10.0.0.2:5000/kube-proxy-arm64:$VERSION
docker image tag k8s.gcr.io/pause-arm64:3.2 10.0.0.2:5000/pause-arm64:3.2

docker push 10.0.0.2:5000/kube-apiserver-arm64:$VERSION
docker push 10.0.0.2:5000/kube-controller-manager-arm64:$VERSION
docker push 10.0.0.2:5000/kube-scheduler-arm64:$VERSION
docker push 10.0.0.2:5000/kube-proxy-arm64:$VERSION
docker push 10.0.0.2:5000/pause-arm64:3.2



