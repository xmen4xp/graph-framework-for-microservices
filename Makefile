DOCKER_GO_PATH ?= /go
DOCKER_BUILD_MOUNT_DIR ?= ${DOCKER_GO_PATH}/src/github.com/vmware-tanzu/graph-framework-for-microservices
NAMESPACE ?= default
CLI_DIR ?= cli
HOST_KUBECONFIG ?= ${HOME}/.kube/config
MOUNTED_KUBECONFIG ?= /etc/config/kubeconfig
RUNTIME_NAMESPACE ?= default
UNAME_S ?= $(shell uname -s)
CWD ?= $(shell pwd)
DATAMODEL_NAME ?= $(shell grep groupName ${DATAMODEL_DIR}/nexus.yaml | cut -f 2 -d" " | tr -d '"')
DATAMODEL_IMAGE_NAME ?= $(shell grep dockerRepo ${DATAMODEL_DIR}/nexus.yaml | cut -f 2 -d" ")
TAG ?= $(shell cat TAG | awk '{ print $1 }')
COMPILER_TAG ?= latest
CLUSTER_PORT ?= 9000
K8S_RUNTIME_PORT = $(shell echo $$(( ${CLUSTER_PORT} + 1 )))
LOG_LEVEL ?= ERROR
DOCKER_REGISTRY ?= amr-registry-pre.caas.intel.com
OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
DOCKER_BUILDER_PLATFORM ?= linux/${ARCH}
ADMIN_DATAMODEL_DEFAULT_RUN_TAG ?= latest
API_GW_DEFAULT_RUN_TAG ?= latest

KIND_API_GW_PORT=$(shell expr ${CLUSTER_PORT} + 1)
API_GW_COMPONENT_NAME ?= api-gw
API_GW_DOCKER_IMAGE ?= ${DOCKER_REGISTRY}/nexus/${API_GW_COMPONENT_NAME}:${TAG}
API_GW_RUN_DOCKER_IMAGE ?= ${DOCKER_REGISTRY}/nexus/${API_GW_COMPONENT_NAME}:${API_GW_DEFAULT_RUN_TAG}

.PHONY: check.prereq.git
check.prereq.git:
	$(info Checking access to required git repositories ...)
	@git ls-remote git@github.com:kubernetes/code-generator.git > /dev/null

.PHONY: check.prereq.go
check.prereq.go:
	$(info Checking if you golang requirements are met...)
	@go version > /dev/null

.PHONY: check.prereq.docker
check.prereq.docker:
	$(info Checking if you docker requirements are met...)
	@docker ps > /dev/null

.PHONY: check.prereq
check.prereq: check.prereq.git check.prereq.go check.prereq.docker

.PHONY: check.env
check.env:
ifndef NEXUS_REPO_DIR
	$(error NEXUS_REPO_DIR is mandatory and should be set to the Nexus repo directory)
endif

.PHONY: check.repodir
check.repodir: check.env
ifneq (${NEXUS_REPO_DIR}, ${CWD})
	$(error NEXUS_REPO_DIR and current working directory does not match)
endif

.PHONY: create.nexus.docker.network
create.nexus.docker.network:
	docker network create nexus | true

.PHONY: rm.nexus.docker.network
rm.nexus.docker.network:
	docker network rm nexus | true

.PHONY: cli.build.darwin
cli.build.darwin:
	docker run \
		--pull=missing \
		--volume $(realpath .):${DOCKER_BUILD_MOUNT_DIR} \
		--volume ${HOME}/.cache/go-build:/root/.cache/go-build \
		-e GOCACHE=/root/.cache/go-build \
		-w ${DOCKER_BUILD_MOUNT_DIR}/${CLI_DIR} \
		golang:1.19.8 \
		/bin/bash -c "git config --global --add safe.directory '*' && go mod download && make build.darwin";

.PHONY: cli.build.linux
cli.build.linux:
	docker run \
		--pull=missing \
		--volume $(realpath .):${DOCKER_BUILD_MOUNT_DIR} \
		-w ${DOCKER_BUILD_MOUNT_DIR}/${CLI_DIR} \
		--volume ${HOME}/.cache/go-build:/root/.cache/go-build \
		-e GOCACHE=/root/.cache/go-build \
		golang:1.19.8 \
		/bin/bash -c "git config --global --add safe.directory '*' && go mod download && make build.linux";

.PHONY: cli.install.darwin
cli.install.darwin:
	cd cli; make install.darwin

.PHONY: cli.install.linux
cli.install.linux:
	cd cli; make install.linux

.PHONY: compiler.builder
compiler.builder:
	cd compiler; DOCKER_REGISTRY=${DOCKER_REGISTRY} BUILDER_TAG=${COMPILER_TAG} make docker.builder

.PHONY: compiler.builder.publish
compiler.builder.publish:
	cd compiler; DOCKER_REGISTRY=${DOCKER_REGISTRY} BUILDER_TAG=${COMPILER_TAG} make docker.builder.publish

.PHONY: compiler.build
compiler.build:
	cd compiler; DOCKER_REGISTRY=${DOCKER_REGISTRY} BUILDER_TAG=${TAG} TAG=${TAG} make docker

.PHONY: compiler.build.publish
compiler.build.publish:
	cd compiler; DOCKER_REGISTRY=${DOCKER_REGISTRY} BUILDER_TAG=${TAG} TAG=${TAG} make docker.publish

#
# Usage: DOCKER_REGISTRY=<registry> TAG=<tag-value> make api-gw.docker
# Tag and publish nexus api gw.
#
# Example: DOCKER_REGISTRY=822995803632.dkr.ecr.us-west-2.amazonaws.com TAG=letstest make api-gw.docker
#
.PHONY: api-gw.docker
api-gw.docker:
	docker run \
		--pull=missing \
		--volume $(realpath .):${DOCKER_BUILD_MOUNT_DIR} \
		-w ${DOCKER_BUILD_MOUNT_DIR}/${API_GW_COMPONENT_NAME} \
		golang:1.19.8 \
		/bin/bash -c "go mod download && GOOS=${OS} GOARCH=${ARCH} go build -buildvcs=false -o bin/${API_GW_COMPONENT_NAME}";
	docker build --platform ${DOCKER_BUILDER_PLATFORM} -t ${API_GW_DOCKER_IMAGE} -f api-gw/Dockerfile .

#
# Usage: DOCKER_REGISTRY=<registry> TAG=<tag-value> make api-gw.docker.kubeconfig
# Tag and publish nexus api gw for execution on K8s cluster.
# This forces the the api-gw to be built for linux/amd64 arch.
#
# Example: DOCKER_REGISTRY=822995803632.dkr.ecr.us-west-2.amazonaws.com TAG=letstest make api-gw.docker.kubeconfig
#
.PHONY: api-gw.docker.kubeconfig
api-gw.docker.kubeconfig:
	docker run \
		--pull=missing \
		--volume $(realpath .):${DOCKER_BUILD_MOUNT_DIR} \
		-w ${DOCKER_BUILD_MOUNT_DIR}/${API_GW_COMPONENT_NAME} \
		golang:1.19.8 \
		/bin/bash -c "go mod download && GOOS=linux GOARCH=amd64 go build -buildvcs=false -o bin/${API_GW_COMPONENT_NAME}";
	docker build --platform ${DOCKER_BUILDER_PLATFORM} -t ${API_GW_DOCKER_IMAGE} -f api-gw/Dockerfile .

.PHONY: api-gw.docker.publish
api-gw.docker.publish:
	docker push ${API_GW_DOCKER_IMAGE}

.PHONY: k0s.install
k0s.install:
	$(realpath .)/nexus-runtime-manifests/k0s/run_k0s.sh

.PHONY: k0s.uninstall
k0s.uninstall:
	$(realpath .)/nexus-runtime-manifests/k0s/stop_k0s.sh

.PHONY: k8s.install
k8s.install:
	$(realpath .)/nexus-runtime-manifests/k8s/run_k8s.sh -p ${CLUSTER_PORT}

.PHONY: k8s.uninstall
k8s.uninstall:
	$(realpath .)/nexus-runtime-manifests/k8s/run_k8s.sh -p ${CLUSTER_PORT} -d

.PHONY: dm.install_init
dm.install_init:
	docker run \
		--net host \
		--pull=missing \
		--volume $(realpath .):${DOCKER_BUILD_MOUNT_DIR} \
		-w ${DOCKER_BUILD_MOUNT_DIR}/nexus-runtime-manifests/datamodel-install \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		alpine/helm \
		upgrade --install datamodel-install-scripts . --set global.namespace=${RUNTIME_NAMESPACE} --kubeconfig ${MOUNTED_KUBECONFIG}

.PHONY: api.build
api.build:
	cd api; DATAMODEL_DOCKER_REGISTRY=${DOCKER_REGISTRY} COMPILER_TAG=${COMPILER_TAG} TAG=${TAG} make docker.build

.PHONY: api.install
api.install:
	cd api; DATAMODEL_DOCKER_REGISTRY=${DOCKER_REGISTRY} TAG=${ADMIN_DATAMODEL_DEFAULT_RUN_TAG} MOUNTED_KUBECONFIG=${MOUNTED_KUBECONFIG} HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k8s/kubeconfig DOCKER_NETWORK=nexus make dm.install

.PHONY: api.install.k0s
api.install.k0s:
	cd api; DATAMODEL_DOCKER_REGISTRY=${DOCKER_REGISTRY} TAG=${ADMIN_DATAMODEL_DEFAULT_RUN_TAG} HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k0s/.kubeconfig MOUNTED_KUBECONFIG=${MOUNTED_KUBECONFIG} DOCKER_NETWORK=nexus make dm.install

.PHONY: api.install.kind
api.install.kind:
	cd api; DATAMODEL_DOCKER_REGISTRY=${DOCKER_REGISTRY} TAG=${ADMIN_DATAMODEL_DEFAULT_RUN_TAG} HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/kind/.${CLUSTER_NAME}/kubeconfig MOUNTED_KUBECONFIG=${MOUNTED_KUBECONFIG} DOCKER_NETWORK=kind make dm.install

.PHONY: api-gw.run
api-gw.run: api-gw.stop
	docker run -d \
		--name=nexus-api-gw \
		--restart=always \
                --network nexus \
		--pull=missing \
		-p 8082:8082 \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		--mount type=bind,source=$(realpath .)/${API_GW_COMPONENT_NAME}/deploy/config/api-gw-config.yaml,target=/api-gw-config.yaml,readonly \
		-e APIGWCONFIG=/api-gw-config.yaml \
		-e KUBECONFIG=${MOUNTED_KUBECONFIG} \
		-e KUBEAPI_ENDPOINT="http://k8s-proxy:8001" \
		${API_GW_RUN_DOCKER_IMAGE}

.PHONY: api-gw.stop
api-gw.stop:
	docker rm -f nexus-api-gw > /dev/null || true

.PHONY: api-gw.kind.run
api-gw.kind.run:
ifndef CLUSTER_NAME
	$(error CLUSTER_NAME is mandatory)
endif
ifndef CLUSTER_PORT
	$(error CLUSTER_PORT is mandatory)
endif
	docker rm -f nexus-api-gw-${CLUSTER_NAME} > /dev/null || true
	docker run -d \
		--name=nexus-api-gw-${CLUSTER_NAME} \
		--rm \
		--network kind \
		--pull=missing \
		-p ${KIND_API_GW_PORT}:${KIND_API_GW_PORT} \
		--mount type=bind,source=$(realpath .)/nexus-runtime-manifests/kind/.${CLUSTER_NAME}/kubeconfig,target=${MOUNTED_KUBECONFIG},readonly \
		--mount type=bind,source=$(realpath .)/nexus-runtime-manifests/kind/.${CLUSTER_NAME}/api-gw-config.yaml,target=/api-gw-config.yaml,readonly \
		-e APIGWCONFIG=/api-gw-config.yaml \
		-e KUBECONFIG=${MOUNTED_KUBECONFIG} \
		-e KUBEAPI_ENDPOINT="http://k8s-proxy-${CLUSTER_NAME}:${CLUSTER_PORT}" \
		-e LOG_LEVEL=${LOG_LEVEL} \
		${API_GW_RUN_DOCKER_IMAGE}

.PHONY: api-gw.k8s.run
api-gw.k8s.run: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k8s/kubeconfig
api-gw.k8s.run:
ifeq (${UNAME_S}, Linux)
	docker rm -f nexus-api-gw-${CLUSTER_NAME} > /dev/null || true
	docker run -d \
		--name=nexus-api-gw \
		--rm \
		--entrypoint /bin/bash \
		--network host \
		--pull=missing \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		--mount type=bind,source=$(realpath .)/${API_GW_COMPONENT_NAME}/deploy/config/api-gw-config.yaml,target=/api-gw-config.yaml,readonly \
		-e APIGWCONFIG=/api-gw-config.yaml \
		-e KUBECONFIG=${MOUNTED_KUBECONFIG} \
		-e KUBEAPI_ENDPOINT="http://127.0.0.1:${CLUSTER_PORT}" \
		-e LOG_LEVEL=${LOG_LEVEL} \
		${API_GW_COMPONENT_NAME}:${TAG} -c "api-gw --metrics-bind-address :8083"
else
	APIGWCONFIG=$(realpath .)/${API_GW_COMPONENT_NAME}/deploy/config/api-gw-config.yaml KUBECONFIG=${HOST_KUBECONFIG} $(realpath .)/${API_GW_COMPONENT_NAME}/bin/${API_GW_COMPONENT_NAME}
endif

.PHONY: api-gw.kind.stop
api-gw.kind.stop:
	docker rm -f nexus-api-gw-${CLUSTER_NAME} > /dev/null || true


.PHONY: runtime.build
runtime.build: api.build api-gw.docker

.PHONY: runtime.clean
runtime.clean:
	rm -rf api/build/*

.PHONY: runtime.install.k0s 
runtime.install.k0s: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k0s/.kubeconfig
runtime.install.k0s: check.repodir create.nexus.docker.network k0s.install api.install.k0s api-gw.run
	$(info )
	$(info ====================================================)
	$(info To access runtime, you can execute kubectl as:)
	$(info     kubectl -s localhost:8082 ...)
	$(info )
	$(info )
	$(info To access nexus api gateway using kubeconfig, export:)
	$(info     export HOST_KUBECONFIG=${HOST_KUBECONFIG})
	$(info )
	$(info ====================================================)

.PHONY: runtime.uninstall.k0s
runtime.uninstall.k0s: k0s.uninstall api-gw.stop rm.nexus.docker.network
	$(info )
	$(info ====================================================)
	$(info Runtime is now uninstalled)
	$(info ====================================================)

.PHONY: kind.install
kind.install:
	$(realpath .)/nexus-runtime-manifests/kind/run_kind.sh -n ${CLUSTER_NAME} -p ${CLUSTER_PORT}

.PHONY: kind.uninstall
kind.uninstall:
	$(realpath .)/nexus-runtime-manifests/kind/run_kind.sh -n ${CLUSTER_NAME} -d

.PHONY: check.kind
check.kind:
	$(info Checking if kind requirements are met...)
	@kind version > /dev/null

.PHONY: runtime.uninstall.kind
runtime.uninstall.kind: check.kind kind.uninstall api-gw.kind.stop
	$(info )
	$(info ====================================================)
	$(info Runtime is now uninstalled)
	$(info ====================================================)

.PHONY: runtime.install.kind
runtime.install.kind: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/kind/.${CLUSTER_NAME}/kubeconfig
runtime.install.kind: check.kind check.repodir kind.install api.install.kind api-gw.kind.run
	$(info )
	$(info ====================================================)
	$(info To access runtime, you can execute kubectl as:)
	$(info     kubectl -s localhost:${K8S_RUNTIME_PORT} ...)
	$(info )
	$(info To access nexus api gateway using kubeconfig, export:)
	$(info     export HOST_KUBECONFIG=${HOST_KUBECONFIG})
	$(info     export DOCKER_NETWORK=kind
	$(info )
	$(info ====================================================)

.PHONY: runtime.install
runtime.install: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k8s/kubeconfig
runtime.install: runtime.uninstall create.nexus.docker.network k8s.install api.install api-gw.run
	$(info )
	$(info ====================================================)
	$(info To access runtime, you can execute kubectl as:)
	$(info     kubectl -s localhost:${CLUSTER_PORT} ...)
	$(info )
	$(info )
	$(info To access nexus api gateway using kubeconfig, export:)
	$(info     export HOST_KUBECONFIG=${HOST_KUBECONFIG})
	$(info )
	$(info ====================================================)

.PHONY: runtime.uninstall
runtime.uninstall: k8s.uninstall api-gw.stop
	$(info )
	$(info ====================================================)
	$(info Runtime is now uninstalled)
	$(info ====================================================)

.PHONY: dm.check-datamodel-dir
dm.check-datamodel-dir:
ifndef DATAMODEL_DIR
	$(error DATAMODEL_DIR is mandatory)
endif

.PHONY: dm.install
dm.install: dm.check-datamodel-dir
	docker run \
		--net host \
		--pull=missing \
		--volume ${DATAMODEL_DIR}/build/crds:/crds \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		--mount type=bind,source=$(realpath .)/nexus-runtime-manifests/datamodel-install/datamodel_installer.sh,target=/datamodel_installer.sh,readonly \
		--entrypoint /datamodel_installer.sh \
		-e KUBECONFIG=${MOUNTED_KUBECONFIG} \
		-e NAME=${DATAMODEL_NAME} \
		-e IMAGE=${DATAMODEL_IMAGE_NAME} \
		bitnami/kubectl

.PHONY: runtime.k0s.kubeconfig.export
runtime.k0s.kubeconfig.export: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k0s/.kubeconfig
runtime.k0s.kubeconfig.export:
	$(info )
	$(info Execute the below export statement on your shell:)
	$(info     export HOST_KUBECONFIG=${HOST_KUBECONFIG})
	$(info )
	@echo > /dev/null

.PHONY: runtime.kind.kubeconfig.export
runtime.kind.kubeconfig.export: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/kind/.${CLUSTER_NAME}/kubeconfig
runtime.kind.kubeconfig.export:
ifndef CLUSTER_NAME
	$(error CLUSTER_NAME is mandatory)
else
	$(info )
	$(info Execute the below export statement on your shell:)
	$(info     export HOST_KUBECONFIG=${HOST_KUBECONFIG})
	$(info     export DOCKER_NETWORK=kind
	$(info )
endif
	@echo > /dev/null

#
# Usage: TAG=<tag-value> make core.runtime.build
# Build artifacts that mke up nexus runtime core.
#
# Example: TAG=latest make core.runtime.build
#
.PHONY: core.runtime.build
core.runtime.build: api.build api-gw.docker.kubeconfig

#
# Usage: DOCKER_REGISTRY=<registry> TAG=<tag-value> make core.runtime.release
# Publish core artifiacats to container registry.
#
# Example: DOCKER_REGISTRY=amr-registry-pre.caas.intel.com TAG=latest make core.runtime.release
#
.PHONY: core.runtime.release
core.runtime.release:
	docker push ${API_GW_DOCKER_IMAGE}
	cd api; DATAMODEL_DOCKER_REGISTRY=${DOCKER_REGISTRY} TAG=${TAG} make docker.publish
#
# Usage: TAG=<tag-value> make core.runtime.download
# Download all artifacts that make up the nexus runtime.
# TAG provides a way to download custom tagged artifacts.
#
# Example: TAG=latest make core.runtime.download
#
.PHONY: core.runtime.download
core.runtime.download:
	docker pull amr-registry-pre.caas.intel.com/nexus/api-gw:${TAG}
	docker pull amr-registry-pre.caas.intel.com/nexus/admin.nexus.com:${TAG}
	docker pull amr-registry-pre.caas.intel.com/nexus/debugtools:latest
	docker pull amr-registry-pre.caas.intel.com/nexus/nexus-kube-apiserver:latest
	docker pull amr-registry-pre.caas.intel.com/nexus/nexus-etcd-kubectl:latest

#
# Usage: DOCKER_REGISTRY=<registry> TAG=<tag-value> make core.runtime.artifacts.push
# Tag and publish all the images needed by nexus rutnime.
#
# Example: DOCKER_REGISTRY=822995803632.dkr.ecr.us-west-2.amazonaws.com TAG=letstest make core.runtime.artifacts.push
#
.PHONY: core.runtime.artifacts.push
core.runtime.artifacts.push:
	docker tag amr-registry-pre.caas.intel.com/nexus/api-gw:${TAG} ${API_GW_DOCKER_IMAGE}
	docker push ${API_GW_DOCKER_IMAGE}
	docker tag amr-registry-pre.caas.intel.com/nexus/admin.nexus.com:${TAG} ${DOCKER_REGISTRY}/nexus/admin.nexus.com:${TAG}
	docker push ${DOCKER_REGISTRY}/nexus/admin.nexus.com:${TAG}
	docker tag amr-registry-pre.caas.intel.com/nexus/debugtools:latest ${DOCKER_REGISTRY}/nexus/debugtools:${TAG}
	docker push ${DOCKER_REGISTRY}/nexus/debugtools:${TAG}
	docker tag amr-registry-pre.caas.intel.com/nexus/nexus-kube-apiserver:latest ${DOCKER_REGISTRY}/nexus/nexus-kube-apiserver:${TAG}
	docker push ${DOCKER_REGISTRY}/nexus/nexus-kube-apiserver:${TAG}
	docker tag amr-registry-pre.caas.intel.com/nexus/nexus-etcd-kubectl:latest ${DOCKER_REGISTRY}/nexus/nexus-etcd-kubectl:${TAG}
	docker push ${DOCKER_REGISTRY}/nexus/nexus-etcd-kubectl:${TAG}

#
# Usage: DOCKER_REGISTRY=<registry> TAG=<tag-value> make core.install.kubeconfig
# Tag and publish all the images needed by nexus rutnime.
#
# Example: DOCKER_REGISTRY=822995803632.dkr.ecr.us-west-2.amazonaws.com TAG=letstest make core.install.kubeconfig
#
.PHONY: core.install.kubeconfig
core.install.kubeconfig:
ifndef TAG
	$(error TAG is mandatory)
endif
ifndef DOCKER_REGISTRY
	$(error DOCKER_REGISTRY is mandatory)
endif
	cd nexus-runtime-manifests/helm-charts/core; helm install core -n ${NAMESPACE} --set imageTag=${TAG} --set global.resources.kubeapiserver.cpu=6 --set global.tainted=true --set imageRegistry=${DOCKER_REGISTRY} .
	cd api; DATAMODEL_DOCKER_REGISTRY=${DOCKER_REGISTRY} TAG=${TAG} make dm.install.helm

.PHONY: core.uninstall.kubeconfig
core.uninstall.kubeconfig:
	helm uninstall -n ${NAMESPACE} core | true
	helm uninstall -n ${NAMESPACE} api | true

.PHONY: demo.install.kubeconfig
demo.install.kubeconfig: core.install.kubeconfig
	cd demo/edge-power-scheduler/datamodel; DATAMODEL_DOCKER_REGISTRY=${DOCKER_REGISTRY} TAG=${TAG} make dm.install.helm
	cd demo/edge-power-scheduler/helm-chart; helm install epm-app --set imageTag=${TAG} --set global.tainted=true --set imageRegistry=${DOCKER_REGISTRY} --namespace ${NAMESPACE}  .
	kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

.PHONY: demo.uninstall.kubeconfig
demo.uninstall.kubeconfig: core.uninstall.kubeconfig
	helm uninstall -n ${NAMESPACE} epm | true
	helm uninstall -n ${NAMESPACE}  epm-app | true
	kubectl delete -n ${NAMESPACE} pvc data-nexus-etcd-0
	kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

.PHONY: trafic.demo.install.kubeconfig
traffic.demo.install.kubeconfig: core.install.kubeconfig
	cd demo/traffic-lights/datamodel; DATAMODEL_DOCKER_REGISTRY=${DOCKER_REGISTRY} TAG=${TAG} make dm.install.helm
	cd demo/traffic-lights/helm-chart; helm install trafficlightapp --set imageTag=${TAG} --set global.tainted=true --set imageRegistry=${DOCKER_REGISTRY} --namespace ${NAMESPACE}  .
	kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

.PHONY: traffic.demo.uninstall.kubeconfig
traffic.demo.uninstall.kubeconfig: core.uninstall.kubeconfig
	helm uninstall -n ${NAMESPACE} trafficlight | true
	helm uninstall -n ${NAMESPACE}  trafficlightapp | true
	kubectl delete -n ${NAMESPACE} pvc data-nexus-etcd-0
	kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

