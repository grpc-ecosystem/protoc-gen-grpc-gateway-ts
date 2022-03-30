#!/bin/bash

# remove binaries to ensure that binaries present in tools.go are installed
rm -f $GOBIN/protoc-gen-go $GOBIN/protoc-gen-go-grpc $GOBIN/protoc-gen-grpc-gateway $GOBIN/protoc-gen-swagger

go install \
	google.golang.org/protobuf/cmd/protoc-gen-go \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

protoc -I . -I ../.. --go_out ./ --go-grpc_out ./ --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./ --grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	service.proto msg.proto
