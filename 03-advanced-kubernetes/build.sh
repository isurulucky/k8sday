#!/usr/bin/env bash

export CGO_ENABLED=0
DOCKER_REPO=<your_dockerhub_id>
DOCKER_TAG=v1.1.0

GOOS=linux GOARCH=amd64 go build -o hello-controller -x ./cmd/controller

docker build -t ${DOCKER_REPO}/k8s-hello-controller:${DOCKER_TAG} .
docker push ${DOCKER_REPO}/k8s-hello-controller:${DOCKER_TAG}
