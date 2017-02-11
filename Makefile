.PHONY: protos

PROTO_PATH := ${GOPATH}/src

protos:
	protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/scheduler/scheduler.proto
	protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/executor/executor.proto
	protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/mesos/mesos.proto
