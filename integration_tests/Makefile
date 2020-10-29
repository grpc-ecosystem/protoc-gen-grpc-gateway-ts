.PHONY: all protos

all: grpc-server-protos client-protos

grpc-server-protos:
	protoc -I . --go_out ./ --go_opt plugins=grpc --go_opt paths=source_relative \
	--grpc-gateway_out ./ --grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	service.proto

client-protos:
	cd .. && go install && cd integration_tests && \
	protoc -I . \
	--grpc-gateway-ts_out=logtostderr=true,loglevel=debug:./ \
	service.proto
