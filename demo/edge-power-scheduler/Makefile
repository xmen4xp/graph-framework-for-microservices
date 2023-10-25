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