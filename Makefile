.PHONY: test test-race bench protos

PROTO_PATH := ${GOPATH}/src

test:
	@go test -cover $(shell go list ./... | grep -v /vendor/)

test-race:
	@go test -race $(shell go list ./... | grep -v /vendor/)

bench:
	@go test -bench . $(shell go list ./... | grep -v /vendor/)

protos:
	@protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/mesos_v1_scheduler/scheduler.proto
	@protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/mesos_v1_executor/executor.proto
	@protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/mesos_v1/mesos.proto
