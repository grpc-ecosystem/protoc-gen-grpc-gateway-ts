#!/bin/bash

# remove binaries to ensure that binaries present in tools.go are installed
rm -f $GOBIN/protoc-gen-go $GOBIN/protoc-gen-grpc-gateway $GOBIN/protoc-gen-swagger

go install \
	github.com/golang/protobuf/protoc-gen-go \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

protoc -I . -I ../.. --go_out ./ --go_opt plugins=grpc --go_opt paths=source_relative \
	--grpc-gateway_out ./ --grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	service.proto msg.proto
