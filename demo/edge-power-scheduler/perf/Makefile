#TAG ?= $(shell cat ../../../TAG | awk '{ print $1 }')
TAG ?= dbakiaraj-perf
CONTAINER_REGISTRY_DOMAIN ?= 822995803632.dkr.ecr.us-west-2.amazonaws.com
COMPONENT_NAME ?= perf

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run: 
	go run main.go

.PHONY: build build.amd64 build.darwin
build: build.amd64 build.darwin

build.amd64:
	mkdir -p bin/amd64
	GOARCH=amd64 GOOS=linux go build -o bin/amd64 ./...

build.darwin:
	mkdir -p bin/darwin
	GOARCH=arm64 GOOS=darwin go build -o bin/darwin ./...

.PHONY: docker
docker:
	docker build -t ${CONTAINER_REGISTRY_DOMAIN}/${COMPONENT_NAME}:${TAG} --progress plain .

.PHONY: docker.push
docker.push:
	docker push ${CONTAINER_REGISTRY_DOMAIN}/${COMPONENT_NAME}:${TAG}

.PHONY: docker.login
docker.login:
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin ${CONTAINER_REGISTRY_DOMAIN}
