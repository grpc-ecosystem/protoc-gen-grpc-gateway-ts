package registry

import (
	"fmt"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"

	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/data"
)

func getHTTPAnnotation(m *descriptorpb.MethodDescriptorProto) *annotations.HttpRule {
	option := proto.GetExtension(m.GetOptions(), annotations.E_Http)
	return option.(*annotations.HttpRule)
}

func hasHTTPAnnotation(m *descriptorpb.MethodDescriptorProto) bool {
	return getHTTPAnnotation(m) != nil
}

func getHTTPMethodPath(m *descriptorpb.MethodDescriptorProto) (method, path string) {
	if !hasHTTPAnnotation(m) {
		return "", ""
	}

	rule := getHTTPAnnotation(m)
	pattern := rule.Pattern
	switch pattern.(type) {
	case *annotations.HttpRule_Get:
		return "GET", rule.GetGet()
	case *annotations.HttpRule_Post:
		return "POST", rule.GetPost()
	case *annotations.HttpRule_Put:
		return "PUT", rule.GetPut()
	case *annotations.HttpRule_Patch:
		return "PATCH", rule.GetPatch()
	case *annotations.HttpRule_Delete:
		return "DELETE", rule.GetDelete()
	default:
		panic(fmt.Sprintf("unsupported HTTP method %T", pattern))
	}
}

func getHTTPBody(m *descriptorpb.MethodDescriptorProto) *string {
	if !hasHTTPAnnotation(m) {
		return nil
	}
	empty := ""
	rule := getHTTPAnnotation(m)
	pattern := rule.Pattern
	switch pattern.(type) {
	case *annotations.HttpRule_Get:
		return &empty
	default:
		body := rule.GetBody()
		return &body
	}
}

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

		httpMethod := "POST"
		url := "/" + serviceURLPart + "/" + method.GetName()
		if hasHTTPAnnotation(method) {
			hm, u := getHTTPMethodPath(method)
			if hm != "" && u != "" {
				httpMethod = hm
				url = u
			}
		}
		body := getHTTPBody(method)

		methodData := &data.Method{
			Name: method.GetName(),
			URL:  url,
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
			HTTPMethod:      httpMethod,
			HTTPRequestBody: body,
		}

		fileData.TrackPackageNonScalarType(methodData.Input)
		fileData.TrackPackageNonScalarType(methodData.Output)

		serviceData.Methods = append(serviceData.Methods, methodData)
	}

	fileData.Services = append(fileData.Services, serviceData)
}
