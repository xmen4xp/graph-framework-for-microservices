COMPILER_TAG=letsplay

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: run
run: 
	go run main.go

.PHONY: dm
dm:
	${MAKE}  -C datamodel COMPILER_TAG=${COMPILER_TAG} datamodel_build 

.PHONEY: dm.install
dm.install:
	${MAKE}  -C datamodel dm.install

#bin/eps-edge-agent: pkg/dminit/* cmd/eps-edge-agent/*
#	go build -o bin powerscheduler/cmd/eps-edge-agent 

.PHONEY: build
build: #bin/eps-edge-agent
	go build -o bin ./...


.PHONEY: k8s-clean
k8s-clean:
	 kubectl delete configs.config.intel.com --all
	 kubectl delete desirededgeconfigs.desiredconfig.intel.com --all
	 kubectl delete edges.edge.intel.com --all
	 kubectl delete edgedcs.edgedc.intel.com --all
	 kubectl delete inventories.inventory.intel.com --all
	 kubectl delete jobs.jobscheduler.intel.com --all
	 kubectl delete powerinfos.powermgmt.intel.com --all
	 kubectl delete powerschedulers.root.intel.com --all
	 kubectl delete jobinfos.jobmgmt.intel.com --all
	 kubectl delete schedulerconfigs.jobscheduler.intel.com --all	 

.PHONEY: run-edge-10
run-edge-10: ./bin/eps-edge-agent
	EDGE_NAME="edge-10" ./bin/eps-edge-agent

.PHONEY: run-edge-11
run-edge-11: ./bin/eps-edge-agent
	EDGE_NAME="edge-11" ./bin/eps-edge-agent

.PHONEY: run-scheduler
run-scheduler: ./bin/eps-scheduler
	./bin/eps-scheduler

.PHONEY: run-generator
run-generator: ./bin/eps-job-requester
	MAXJOBS=2 ./bin/eps-job-requester	