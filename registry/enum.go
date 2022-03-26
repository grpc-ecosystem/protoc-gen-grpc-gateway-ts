package registry

import (
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/data"
)

const (
	// enumDescriptorValueFieldNumber is the field number of the field EnumDescriptorProro.value
	enumDescriptorValueFieldNumber = 2
)

func (r *Registry) analyseEnumType(fileData *data.File, packageName, fileName string, parents []string, enum *descriptorpb.EnumDescriptorProto, commentInfo *CommentInfo) {
	packageIdentifier := r.getNameOfPackageLevelIdentifier(parents, enum.GetName())
	fqName := r.getFullQualifiedName(packageName, parents, enum.GetName())
	protoType := descriptorpb.FieldDescriptorProto_TYPE_ENUM
	r.Types[fqName] = &TypeInformation{
		FullyQualifiedName: fqName,
		Package:            packageName,
		File:               fileName,
		PackageIdentifier:  packageIdentifier,
		LocalIdentifier:    enum.GetName(),
		ProtoType:          protoType,
	}

	enumData := data.NewEnum()
	enumData.Name = packageIdentifier
	enumData.Comment = commentInfo.GetText()

	for i, e := range enum.GetValue() {
		value := &data.EnumValue{
			Name:    e.GetName(),
			Comment: commentInfo.GetSubComment(enumDescriptorValueFieldNumber, i).GetText(),
		}
		enumData.Values = append(enumData.Values, value)
	}

	fileData.Enums = append(fileData.Enums, enumData)
}
