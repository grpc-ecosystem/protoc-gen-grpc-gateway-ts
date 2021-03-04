#!/bin/bash
USE_PROTO_NAMES=${1:-"false"}
cd .. && go install && cd integration_tests && \
	protoc -I . \
	--grpc-gateway-ts_out=logtostderr=true,use_proto_names=$USE_PROTO_NAMES,loglevel=debug:./ \
	service.proto msg.proto empty.proto