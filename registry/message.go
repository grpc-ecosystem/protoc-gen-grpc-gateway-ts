package registry

import (
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/data"
)

func (r *Registry) analyseMessage(fileData *data.File, packageName, fileName string, parents []string, message *descriptorpb.DescriptorProto) {
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

	newParents := append(parents, message.GetName())

	// handle enums, by pulling the enums out to the top level
	for _, enum := range message.EnumType {
		r.analyseEnumType(fileData, packageName, fileName, newParents, enum)
	}

	// nested type also got pull out to the top level of the file
	for _, msg := range message.NestedType {
		r.analyseMessage(fileData, packageName, fileName, newParents, msg)
	}

	// store a map of one of names
	for idx, oneOf := range message.GetOneofDecl() {
		data.OneOfFieldsNames[int32(idx)] = oneOf.GetName()
	}

	// analyse fields in the messages
	for _, f := range message.Field {
		r.analyseField(fileData, data, packageName, f)
	}

	fileData.Messages = append(fileData.Messages, data)
}
