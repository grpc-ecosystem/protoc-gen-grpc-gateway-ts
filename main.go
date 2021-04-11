package main

import (
	"io/ioutil"
	"os"
	"strings"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	log "github.com/sirupsen/logrus" // nolint: depguard
	"google.golang.org/protobuf/proto"

	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/generator"
	"github.com/pkg/errors"
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
	paramsMap := getParamsMap(req)
	err := configureLogging(paramsMap)
	if err != nil {
		panic(err)
	}

	g, err := generator.New(paramsMap)
	if err != nil {
		panic(err)
	}

	log.Debug("Starts generating file request")
	resp, err := g.Generate(req)
	if err != nil {
		panic(err)
	}

	encodeResponse(resp)
	log.Debug("generation finished")
}

func configureLogging(paramsMap map[string]string) error {
	if paramsMap["logtostderr"] == "true" { // configure logging when it's in the options
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
		log.SetOutput(os.Stderr)
		log.Debugf("Logging configured completed, logging has been enabled")
		levelStr := paramsMap["loglevel"]
		if levelStr != "" {
			level, err := log.ParseLevel(levelStr)
			if err != nil {
				return errors.Wrapf(err, "error parsing log level %s", levelStr)
			}

			log.SetLevel(level)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	}

	return nil
}

func getParamsMap(req *plugin.CodeGeneratorRequest) map[string]string {
	paramsMap := make(map[string]string)
	params := req.GetParameter()

	for _, p := range strings.Split(params, ",") {
		if i := strings.Index(p, "="); i < 0 {
			paramsMap[p] = ""
		} else {
			paramsMap[p[0:i]] = p[i+1:]
		}
	}

	return paramsMap
}
