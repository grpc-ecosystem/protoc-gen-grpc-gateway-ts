.PHONY: all protos

all: protos

protos:
	cd .. && go install && cd testdata && \
	protoc -I . \
	--grpc-gateway-ts_out=logtostderr=true,loglevel=debug:./ \
    log.proto environment.proto ./datasource/datasource.proto
