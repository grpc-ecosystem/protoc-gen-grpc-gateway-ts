package registry

import (
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/data"
)

const (
	// messageDescriptorFieldFieldNumber is the field number of the field Descriptor.field
	messageDescriptorFieldFieldNumber = 2
	// messageDescriptorNestedTypeFieldNumber is the field number of the field Descriptor.nested_type
	messageDescriptorNestedTypeFieldNumber = 3
	// messageDescriptorEnumTypeFieldNumber is the field number of the field Descriptor.enum_type
	messageDescriptorEnumTypeFieldNumber = 4
)

func (r *Registry) analyseMessage(fileData *data.File, packageName, fileName string, parents []string, message *descriptorpb.DescriptorProto, commentInfo *CommentInfo) {
	packageIdentifier := r.getNameOfPackageLevelIdentifier(parents, message.GetName())

	fqName := r.getFullQualifiedName(packageName, parents, message.GetName()) // "." + packageName + "." + parentsPrefix + message.GetName()
	protoType := descriptorpb.FieldDescriptorProto_TYPE_MESSAGE

	typeInfo := &TypeInformation{
		FullyQualifiedName: fqName,
		Package:            packageName,
		File:               fileName,
		PackageIdentifier:  packageIdentifier,
		LocalIdentifier:    message.GetName(),
		ProtoType:          protoType,
	}

	// register itself in the registry map
	r.Types[fqName] = typeInfo

	if message.Options != nil {
		if message.GetOptions().GetMapEntry() {
			// is a map entry, need to find out the type for key and value
			typeInfo.IsMapEntry = true

			for _, f := range message.Field {
				switch f.GetName() {
				case "key":
					typeInfo.KeyType = &data.MapEntryType{
						Type:       r.getFieldType(f),
						IsExternal: r.isExternalDependenciesOutsidePackage(f.GetTypeName(), packageName),
					}
				case "value":
					typeInfo.ValueType = &data.MapEntryType{
						Type:       r.getFieldType(f),
						IsExternal: r.isExternalDependenciesOutsidePackage(f.GetTypeName(), packageName),
					}
				}

			}
			fileData.TrackPackageNonScalarType(typeInfo.KeyType)
			fileData.TrackPackageNonScalarType(typeInfo.ValueType)
			// no need to add a map type into
			return

		}
	}

	data := data.NewMessage()
	data.Name = packageIdentifier
	data.FQType = fqName
	data.Comment = commentInfo.GetText()

	newParents := append(parents, message.GetName())

	// handle enums, by pulling the enums out to the top level
	for idx, enum := range message.EnumType {
		r.analyseEnumType(fileData, packageName, fileName, newParents, enum, commentInfo.GetSubComment(messageDescriptorEnumTypeFieldNumber, idx))
	}

	// nested type also got pull out to the top level of the file
	for idx, msg := range message.NestedType {
		r.analyseMessage(fileData, packageName, fileName, newParents, msg, commentInfo.GetSubComment(messageDescriptorNestedTypeFieldNumber, idx))
	}

	// store a map of one of names
	for idx, oneOf := range message.GetOneofDecl() {
		data.OneOfFieldsNames[int32(idx)] = oneOf.GetName()
	}

	// analyse fields in the messages
	for idx, f := range message.Field {
		r.analyseField(fileData, data, packageName, f, commentInfo.GetSubComment(messageDescriptorFieldFieldNumber, idx))
	}

	fileData.Messages = append(fileData.Messages, data)
}
