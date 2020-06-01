package registry

import (
	"git.sqcorp.co/cash/gap/cmd/protoc-gen-grpc-gateway-ts/data"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (r *Registry) analyseFile(f *descriptorpb.FileDescriptorProto) *data.File {
	fileData := data.NewFile()
	fileName := f.GetName()
	packageName := f.GetPackage()
	parents := make([]string, 0)
	fileData.Name = fileName
	fileData.TSFileName = data.GetTSFileName(fileName)

	// analyse enums
	for _, enum := range f.EnumType {
		r.analyseEnumType(fileData, packageName, fileName, parents, enum)
	}

	// analyse messages, each message will go recursively
	for _, message := range f.MessageType {
		r.analyseMessage(fileData, packageName, fileName, parents, message)
	}

	// when we have a service we will need to pull functions from gap admin to make the call
	if len(f.Service) > 0 {
		fileData.Dependencies = append(fileData.Dependencies, &data.Dependency{
			ModuleIdentifier: "gap",
			SourceFile:       "gap/admin/lib/useGapFetch",
		})
	}

	// analyse services
	for _, service := range f.Service {
		r.analyseService(fileData, packageName, fileName, service)
	}

	return fileData
}
