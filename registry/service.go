package registry

import (
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"git.sqcorp.co/cash/gap/cmd/protoc-gen-grpc-gateway-ts/data"
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
		isInputTypeExternal := r.isExternalDependenciesOutsidePackage(inputTypeFQName, packageName)

		if isInputTypeExternal {
			fileData.ExternalDependingTypes = append(fileData.ExternalDependingTypes, inputTypeFQName)
		}

		outputTypeFQName := *method.OutputType
		isOutputTypeExternal := r.isExternalDependenciesOutsidePackage(outputTypeFQName, packageName)

		if isOutputTypeExternal {
			fileData.ExternalDependingTypes = append(fileData.ExternalDependingTypes, outputTypeFQName)
		}

		methodData := &data.Method{
			Name: method.GetName(),
			URL:  "/api/" + serviceURLPart + "/" + method.GetName(),
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

		fileData.TrackPackageNonScalarType(methodData.Input)
		fileData.TrackPackageNonScalarType(methodData.Output)

		serviceData.Methods = append(serviceData.Methods, methodData)
	}

	fileData.Services = append(fileData.Services, serviceData)
}
