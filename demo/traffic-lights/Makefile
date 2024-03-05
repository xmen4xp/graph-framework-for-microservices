COMPILER_TAG ?= $(shell cat ../../TAG | awk '{ print $1 }')
TAG ?= $(shell cat ../../TAG | awk '{ print $1 }')
#TAG ?= dbakiaraj-perf
NAMESPACE ?= default
CONTAINER_REGISTRY_DOMAIN ?= 822995803632.dkr.ecr.us-west-2.amazonaws.com

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run: 
	go run main.go

.PHONY: dm
dm:
	${MAKE} -C datamodel COMPILER_TAG=${COMPILER_TAG} datamodel_build

.PHONY: dm.install
dm.install:
	${MAKE} -C datamodel dm.install

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
	docker build -t ${CONTAINER_REGISTRY_DOMAIN}/traffic-light:${TAG} .

.PHONY: docker.push
docker.push:
	docker push ${CONTAINER_REGISTRY_DOMAIN}/traffic-light:${TAG}

.PHONY: docker.login
docker.login:
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin ${CONTAINER_REGISTRY_DOMAIN}
	
.PHONY: test
test:
	go test -v ./...

.PHONY: k8s-clean
k8s-clean:
	kubectl ${K8S_ARGS} delete -f ./datamodel/build/crds
	kubectl ${K8S_ARGS} apply -f ./datamodel/build/crds

.PHONY: apis
apis:
	kubectl -s localhost:8001 get crds --no-headers | grep trafficlight

.PHONY: system-status
system-status:
	watch -n3 "curl -s -m 60 'localhost:8080/systemstatus'"

.PHONY: register-10k-lights
register-10k-lights:
	kubectl scale --replicas=50 statefulset/light-connector

.PHONY: broadcast
broadcast:
	curl -s -m 120 localhost:8080/broadcast

.PHONY: lightcoldboot
lightcoldboot:
	curl -s -m 120 localhost:8080/lightcoldboot

.PHONY: show-light-config
show-light-config:
	kubectl -s localhost:8001 get lightconfigs.desiredconfiglight.trafficlight.com --no-headers | wc -l
