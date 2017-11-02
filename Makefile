# Copyright 2017 Verizon
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: test test-race bench protos

PROTO_PATH := ${GOPATH}/src

test:
	@go test -timeout 5m -cover $(shell go list ./... | grep -v /vendor/)

test-race:
	@go test -timeout 5m -race $(shell go list ./... | grep -v /vendor/)

bench:
	@go test -timeout 5m -bench . $(shell go list ./... | grep -v /vendor/)

protos:
	@protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/mesos_v1_scheduler/scheduler.proto
	@protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/mesos_v1_executor/executor.proto
	@protoc --go_out=. --proto_path=.:${PROTO_PATH} ./include/mesos_v1/mesos.proto
