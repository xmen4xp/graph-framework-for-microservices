COMPILER_TAG ?= $(shell cat ../../TAG | awk '{ print $1 }')
TAG ?= $(shell cat ../../TAG | awk '{ print $1 }')
#TAG ?= dbakiaraj-perf
NAMESPACE ?= default
APIGW_PORT=8000
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

#bin/eps-edge-agent: pkg/dminit/* cmd/eps-edge-agent/*
#	go build -o bin powerscheduler/cmd/eps-edge-agent 

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
	docker build -t ${CONTAINER_REGISTRY_DOMAIN}/edge-power-manager:${TAG} .

.PHONY: docker.push
docker.push:
	docker push ${CONTAINER_REGISTRY_DOMAIN}/edge-power-manager:${TAG}

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
	# kubectl ${K8S_ARGS} delete configs.config.intel.com --all
	# kubectl ${K8S_ARGS} delete desirededgeconfigs.runtimedesiredconfig.intel.com --all
	# kubectl ${K8S_ARGS} delete edges.edge.intel.com --all
	# kubectl ${K8S_ARGS} delete edgedcs.edgedc.intel.com --all
	# kubectl ${K8S_ARGS} delete inventories.inventory.intel.com --all
	# kubectl ${K8S_ARGS} delete jobs.jobscheduler.intel.com --all
	# kubectl ${K8S_ARGS} delete jobinfos.jobmgmt.intel.com --all
	# kubectl ${K8S_ARGS} delete schedulerconfigs.jobgroup.intel.com --all
	# kubectl ${K8S_ARGS} delete powerschedulers.root.intel.com --all
	# kubectl ${K8S_ARGS} delete jobgroups.jobgroup.intel.com --all

.PHONY: eps-0
eps-0:
	kubectl patch statefulsets eps-edge-agent  -p '{"spec":{"replicas":0}}'
	kubectl patch statefulsets eps-generator -p '{"spec":{"replicas":0}}'
	kubectl patch statefulsets eps-scheduler  -p '{"spec":{"replicas":0}}'
# To install runtime 
# Make sure to stop k0s
# sudo k0s stop
# DM_APIGW_PORT=${APIGW_PORT} EDGE_NAME="edge-10" ./bin/eps-edge-agent


.PHONY: run-edge-10
run-edge-10: ./bin/eps-edge-agent
	DM_APIGW_PORT=${APIGW_PORT} EDGE_NAME="edge-10" ./bin/eps-edge-agent

.PHONY: run-edge-11
run-edge-11: ./bin/eps-edge-agent
	DM_APIGW_PORT=${APIGW_PORT} EDGE_NAME="edge-11" ./bin/eps-edge-agent

.PHONY: run-edge-agent
run-edge-agent: ./bin/eps-edge-agent
	$(info     Starting eps edge agent app ...)
ifndef HOST_KUBECONFIG
	DM_APIGW_PORT=${APIGW_PORT} EDGE_NAME="edge-agent" ./bin/eps-edge-agent 2>&1 | tee eps-edge-agent.log&
else
	EDGE_NAME="edge-agent" ./bin/eps-edge-agent -k ${HOST_KUBECONFIG} 2>&1 | tee eps-edge-agent.log&
endif
	sleep 2

.PHONY: run-scheduler
run-scheduler: ./bin/eps-scheduler
	$(info     Starting eps scheduler app ...)
ifndef HOST_KUBECONFIG
	DM_APIGW_PORT=8000 ./bin/eps-scheduler 2>&1 | tee eps-scheduler.log& 
else
	./bin/eps-scheduler -k ${HOST_KUBECONFIG} 2>&1 | tee eps-scheduler.log&
endif
	sleep 2

.PHONY: run-generator
run-generator: ./bin/eps-job-requester
	$(info     Starting eps job requestor app ...)
ifndef HOST_KUBECONFIG
	DM_APIGW_PORT=8000 ./bin/eps-job-requester 2>&1 | tee eps-job-requestor.log& 
else
	./bin/eps-job-requester -k ${HOST_KUBECONFIG} 2>&1 | tee eps-job-requestor.log&
endif
	sleep 2

.PHONY: run-eps-app
run-eps-app: run-scheduler run-edge-agent run-generator
	sleep 3

.PHONY: stop-eps-app
stop-eps-app:
	- pkill -f eps-job-requester
	- pkill -f eps-edge-agent
	- pkill -f eps-scheduler

.PHONY: port-forward
port-forward:
	pkill -f kubectl | true
	kubectl port-forward deploy/nexus-api-gw 8080:8082 &
	kubectl port-forward svc/nexus-k8s-proxy 8082:80 &

.PHONY: watch-job-config
watch-job-config:
ifndef HOST_KUBECONFIG
	kubectl -s localhost:8082 get jobs.jobscheduler.intel.com -w
else
	kubectl --kubeconfig ${HOST_KUBECONFIG} get jobinfos.jobmgmt.intel.com -w
endif

.PHONY: watch-job-dc
watch-job-dc:
ifndef HOST_KUBECONFIG
	kubectl -s localhost:8082 get jobinfos.jobmgmt.intel.com -w
else
	kubectl --kubeconfig ${HOST_KUBECONFIG} get jobinfos.jobmgmt.intel.com -w
endif

.PHONY: watch-edges
watch-edges:
ifndef HOST_KUBECONFIG
	kubectl -s localhost:8082 get edges.edge.intel.com -w
else
	kubectl --kubeconfig ${HOST_KUBECONFIG} get edges.edge.intel.com -w
endif

.PHONY: watch-jobs
watch-jobs:
ifndef HOST_KUBECONFIG
	watch -n 4 'kubectl -s localhost:8082 get schedulerconfigs.jobgroup.intel.com -o json | jq -r ".items[].status.state"'
else
	watch -n 4 'kubectl --kubeconfig ${HOST_KUBECONFIG} get schedulerconfigs.jobscheduler.intel.com -o yaml  | egrep "dges|job|total"'
endif

.PHONY: demo.init
demo.init:
ifndef HOST_KUBECONFIG
	kubectl -s localhost:8080 create -f demo-yamls/jobgroup.yaml
else
	kubectl --kubeconfig ${HOST_KUBECONFIG} create -f demo-yamls/jobgroup.yaml
endif

.PHONY: stop 
stop: stop-eps-app
	make -C ../.. runtime.uninstall.k8s

.PHONY: start 
start:
	make -C ../.. runtime.install.k8s
	HOST_KUBECONFIG=~/applications.development.framework.nexus/nexus-runtime-manifests/k8s/kubeconfig make dm.install

.PHONY: run-edges
run-edges: 
	./runsim.sh 0 10 ${APIGW_PORT}

.PHONY: run-prom
run-prom:
	docker run -d --rm --name prom --net host  --mount type=bind,source="${HOME}/work/prometheus/prom.yaml",target=/etc/prometheus/prometheus.yml,readonly prom/prometheus

.PHONY: run-grafana
run-grafana:
	docker run -d --rm --name grafana --network host  --publish 3000:3000  grafana/grafana-oss:latest

.PHONY: docker.perf
docker.perf:
	docker build -t ${CONTAINER_REGISTRY_DOMAIN}/perf:${TAG} -f Dockerfile.perf --progress plain .

.PHONY: docker.push.perf
docker.push.perf:
	docker push ${CONTAINER_REGISTRY_DOMAIN}/perf:${TAG}

.PHONY: scale
scale:
	kubectl scale --replicas=${SCALE} statefulset/eps-scheduler
	kubectl scale --replicas=${SCALE} statefulset/eps-generator
	kubectl scale --replicas=${SCALE} statefulset/eps-edge-agent
