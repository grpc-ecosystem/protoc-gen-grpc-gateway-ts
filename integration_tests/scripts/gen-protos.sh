#!/bin/bash
USE_PROTO_NAMES=${1:-"false"}
ENABLE_STYLING_CHECK=${2:-"false"}
cd .. && go install && cd integration_tests && \
	protoc -I .  -I ../.. \
	--grpc-gateway-ts_out=logtostderr=true,use_proto_names=$USE_PROTO_NAMES,enable_styling_check=$ENABLE_STYLING_CHECK,loglevel=debug:./ \
	service.proto msg.proto empty.proto