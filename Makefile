.PHONY: protos

PROTO_PATH := ${GOPATH}/src

test:
	go test ./...

test-race:
	go test -race ./...

protos:
	protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/scheduler/scheduler.proto
	protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/executor/executor.proto
	protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/mesos/mesos.proto
