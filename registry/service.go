package registry

import (
	"git.sqcorp.co/cash/gap/cmd/protoc-gen-grpc-gateway-ts/data"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (r *Registry) analyseService(fileData *data.File, packageName string, fileName string, service *descriptorpb.ServiceDescriptorProto) {
	packageIdentifier := service.GetName()
	fqName := "." + packageName + "." + packageIdentifier

	// register itself in the registry map
	r.Types[fqName] = &TypeInformation{
		FullyQualifiedName: fqName,
		Package:            packageName,
		File:               fileName,
		PackageIdentifier:  packageIdentifier,
		LocalIdentifier:    service.GetName(),
	}

	serviceData := data.NewService()
	serviceData.Name = service.GetName()
	serviceURLPart := packageName + "." + serviceData.Name

	for _, method := range service.Method {
		// don't support client streaming, will ignore the client streaming method
		if method.GetClientStreaming() {
			continue
		}

		inputTypeFQName := *method.InputType
		isInputTypeExternal := r.isExternalDependencies(inputTypeFQName, packageName)

		if isInputTypeExternal {
			fileData.ExternalDependingTypes = append(fileData.ExternalDependingTypes, inputTypeFQName)
		}

		outputTypeFQName := *method.OutputType
		isOutputTypeExternal := r.isExternalDependencies(outputTypeFQName, packageName)

		if isOutputTypeExternal {
			fileData.ExternalDependingTypes = append(fileData.ExternalDependingTypes, outputTypeFQName)
		}

		methodData := &data.Method{
			Name: method.GetName(),
			URL:  "/" + serviceURLPart + "/" + method.GetName(),
			Input: &data.MethodArgument{
				Type:       inputTypeFQName,
				IsExternal: isInputTypeExternal,
			},
			Output: &data.MethodArgument{
				Type:       outputTypeFQName,
				IsExternal: isOutputTypeExternal,
			},
			ServerStreaming: method.GetServerStreaming(),
			ClientStreaming: method.GetClientStreaming(),
		}

		serviceData.Methods = append(serviceData.Methods, methodData)
	}

	fileData.Services = append(fileData.Services, serviceData)
}
