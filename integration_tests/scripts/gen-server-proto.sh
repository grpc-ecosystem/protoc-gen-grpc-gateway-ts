#!/bin/bash

protoc -I . --go_out ./ --go_opt plugins=grpc --go_opt paths=source_relative \
	--grpc-gateway_out ./ --grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	service.proto