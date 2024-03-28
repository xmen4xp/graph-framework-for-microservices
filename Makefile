DOCKER_GO_PATH ?= /go
DOCKER_BUILD_MOUNT_DIR ?= ${DOCKER_GO_PATH}/src/github.com/vmware-tanzu/graph-framework-for-microservices
API_GW_COMPONENT_NAME ?= api-gw
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
K8S_RUNTIME_PORT = $(shell echo $$(( ${CLUSTER_PORT} + 1 )))
LOG_LEVEL ?= ERROR
CONTAINER_REGISTRY_DOMAIN ?= 822995803632.dkr.ecr.us-west-2.amazonaws.com
OS ?= $(shell go env GOOS)
ARCH ?= $(shell go env GOARCH)
DOCKER_BUILDER_PLATFORM ?= linux/${ARCH}

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
	cd compiler; BUILDER_TAG=${TAG} make docker.builder

.PHONY: compiler.builder.image.exists
compiler.builder.image.exists:
	@docker inspect nexus/compiler-builder:${TAG} >/dev/null 2>&1 || \
		(echo "Image  nexus/compiler-builder:${TAG} does not exist. Run 'make compiler.builder' to generate it" && false)

.PHONY: compiler.build
compiler.build: compiler.builder.image.exists
	cd compiler; BUILDER_TAG=${TAG} TAG=${TAG} make docker

.PHONY: api-gw.docker
api-gw.docker:
	docker run \
		--pull=missing \
		--volume $(realpath .):${DOCKER_BUILD_MOUNT_DIR} \
		-w ${DOCKER_BUILD_MOUNT_DIR}/${API_GW_COMPONENT_NAME} \
		golang:1.19.8 \
		/bin/bash -c "go mod download && GOOS=${OS} GOARCH=${ARCH} go build -buildvcs=false -o bin/${API_GW_COMPONENT_NAME}";
	docker build --platform ${DOCKER_BUILDER_PLATFORM} -t ${API_GW_COMPONENT_NAME}:${TAG} -f api-gw/Dockerfile .

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
	$(realpath .)/nexus-runtime-manifests/k8s/stop_k8s.sh

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
	cd api; COMPILER_TAG=${TAG} VERSION=${TAG} make datamodel_build

.PHONY: api.install
api.install:
	docker run \
		--entrypoint /datamodel_installer.sh \
		--net host \
		--pull=missing \
		--volume $(realpath .)/api/build/crds:/crds \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		--mount type=bind,source=$(realpath .)/nexus-runtime-manifests/datamodel-install/datamodel_installer.sh,target=/datamodel_installer.sh,readonly \
		-e KUBECONFIG=${MOUNTED_KUBECONFIG} \
		-e NAME=admin.nexus.com \
		-e IMAGE=gcr.io/nsx-sm/nexus/nexus-api \
		bitnami/kubectl

.PHONY: api-gw.run
api-gw.run: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k0s/.kubeconfig
api-gw.run: api-gw.stop
ifeq (${UNAME_S}, Linux)
	docker run -d \
		--name=nexus-api-gw \
		--rm \
                --network host \
		--pull=missing \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		--mount type=bind,source=$(realpath .)/${API_GW_COMPONENT_NAME}/deploy/config/api-gw-config.yaml,target=/api-gw-config.yaml,readonly \
		-e APIGWCONFIG=/api-gw-config.yaml \
		-e KUBECONFIG=${MOUNTED_KUBECONFIG} \
		-e KUBEAPI_ENDPOINT="http://127.0.0.1:6443" \
		${API_GW_COMPONENT_NAME}:${TAG}
else
	APIGWCONFIG=$(realpath .)/${API_GW_COMPONENT_NAME}/deploy/config/api-gw-config.yaml KUBECONFIG=${HOST_KUBECONFIG} $(realpath .)/${API_GW_COMPONENT_NAME}/bin/${API_GW_COMPONENT_NAME}
endif

.PHONY: api-gw.stop
api-gw.stop:
	docker rm -f nexus-api-gw > /dev/null || true

.PHONY: api-gw.kind.run
api-gw.kind.run: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/kind/.${CLUSTER_NAME}/kubeconfig
api-gw.kind.run:
ifeq (${UNAME_S}, Linux)
	docker rm -f nexus-api-gw-${CLUSTER_NAME} > /dev/null || true
	docker run -d \
		--name=nexus-api-gw-${CLUSTER_NAME} \
		--rm \
		--network host \
		--pull=missing \
		--mount type=bind,source=${HOST_KUBECONFIG},target=${MOUNTED_KUBECONFIG},readonly \
		--mount type=bind,source=$(realpath .)/nexus-runtime-manifests/kind/.${CLUSTER_NAME}/api-gw-config.yaml,target=/api-gw-config.yaml,readonly \
		-e APIGWCONFIG=/api-gw-config.yaml \
		-e KUBECONFIG=${MOUNTED_KUBECONFIG} \
		-e KUBEAPI_ENDPOINT="http://127.0.0.1:${CLUSTER_PORT}" \
		-e LOG_LEVEL=${LOG_LEVEL} \
		${API_GW_COMPONENT_NAME}:${TAG}
else
	APIGWCONFIG=$(realpath .)/nexus-runtime-manifests/kind/.${CLUSTER_NAME}/api-gw-config.yaml KUBECONFIG=${HOST_KUBECONFIG} $(realpath .)/${API_GW_COMPONENT_NAME}/bin/${API_GW_COMPONENT_NAME}
endif

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
runtime.build: compiler.build api.build api-gw.docker

.PHONY: runtime.clean
runtime.clean:
	rm -rf api/build/*

.PHONY: runtime.install.k0s 
runtime.install.k0s: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k0s/.kubeconfig
runtime.install.k0s: check.repodir k0s.install dm.install_init api.install api-gw.run
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
runtime.uninstall.k0s: k0s.uninstall api-gw.stop
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
runtime.install.kind: check.kind check.repodir kind.install dm.install_init api.install api-gw.kind.run
	$(info )
	$(info ====================================================)
	$(info To access runtime, you can execute kubectl as:)
	$(info     kubectl -s localhost:${K8S_RUNTIME_PORT} ...)
	$(info )
	$(info To access nexus api gateway using kubeconfig, export:)
	$(info     export HOST_KUBECONFIG=${HOST_KUBECONFIG})
	$(info )
	$(info ====================================================)

.PHONY: runtime.install.k8s
runtime.install.k8s: HOST_KUBECONFIG=$(realpath .)/nexus-runtime-manifests/k8s/kubeconfig
runtime.install.k8s: runtime.uninstall.k8s k8s.install dm.install_init api.install api-gw.k8s.run
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

.PHONY: runtime.uninstall.k8s
runtime.uninstall.k8s: k8s.uninstall api-gw.stop
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
	$(info )
endif
	@echo > /dev/null

.PHONY: demo.install.kubeconfig
demo.install.kubeconfig:
ifndef TAG
	$(error TAG is mandatory)
endif
ifndef CONTAINER_REGISTRY_DOMAIN
	$(error CONTAINER_REGISTRY_DOMAIN is mandatory)
endif
	cd nexus-runtime-manifests/helm-charts/core; helm install core -n ${NAMESPACE} --set imageTag=${TAG} --set global.resources.kubeapiserver.cpu=6 --set global.tainted=true --set imageRegistry=${CONTAINER_REGISTRY_DOMAIN} .
	cd api; make dm.install.helm
	cd demo/edge-power-scheduler/datamodel; make dm.install.helm
	cd demo/edge-power-scheduler/helm-chart; helm install epm-app --set imageTag=${TAG} --set global.tainted=true --set imageRegistry=${CONTAINER_REGISTRY_DOMAIN} --namespace ${NAMESPACE}  .
	kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

.PHONY: demo.uninstall.kubeconfig
demo.uninstall.kubeconfig:
	helm uninstall -n ${NAMESPACE} core | true
	helm uninstall -n ${NAMESPACE} api | true
	helm uninstall -n ${NAMESPACE} epm | true
	helm uninstall -n ${NAMESPACE}  epm-app | true
	kubectl delete -n ${NAMESPACE} pvc data-nexus-etcd-0
	kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

.PHONY: trafic.demo.install.kubeconfig
traffic.demo.install.kubeconfig:
ifndef TAG
	$(error TAG is mandatory)
endif
ifndef CONTAINER_REGISTRY_DOMAIN
	$(error CONTAINER_REGISTRY_DOMAIN is mandatory)
endif
	cd nexus-runtime-manifests/helm-charts/core; helm install core -n ${NAMESPACE} --set imageTag=${TAG} --set global.resources.kubeapiserver.cpu=6 --set global.tainted=true --set imageRegistry=${CONTAINER_REGISTRY_DOMAIN} .
	cd api; make dm.install.helm
	cd demo/traffic-lights/datamodel; make dm.install.helm
	cd demo/traffic-lights/helm-chart; helm install trafficlightapp --set imageTag=${TAG} --set global.tainted=true --set imageRegistry=${CONTAINER_REGISTRY_DOMAIN} --namespace ${NAMESPACE}  .
	kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

.PHONY: traffic.demo.uninstall.kubeconfig
traffic.demo.uninstall.kubeconfig:
	helm uninstall -n ${NAMESPACE} core | true
	helm uninstall -n ${NAMESPACE} api | true
	helm uninstall -n ${NAMESPACE} trafficlight | true
	helm uninstall -n ${NAMESPACE}  trafficlightapp | true
	kubectl delete -n ${NAMESPACE} pvc data-nexus-etcd-0
	kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
