package main

import (
	"git.sqcorp.co/cash/gap/cmd/protoc-gen-grpc-gateway-ts/generator"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"io/ioutil"
	"os"
)

func decodeReq() *plugin.CodeGeneratorRequest {
	req := &plugin.CodeGeneratorRequest{}
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	err = proto.Unmarshal(data, req)
	if err != nil {
		panic(err)
	}
	return req
}

func encodeResponse(resp proto.Message) {
	data, err := proto.Marshal(resp)
	if err != nil {
		panic(err)
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		panic(err)
	}
}

func main() {
	req := decodeReq()
	f, err := os.Create("./protoc-gen-grpc-gateway-ts.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	g := generator.New()

	resp, err := g.Generate(req)
	if err != nil {
		panic(err)
	}

	encodeResponse(resp)
}
