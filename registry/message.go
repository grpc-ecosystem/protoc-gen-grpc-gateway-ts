package registry

import (
	"git.sqcorp.co/cash/gap/cmd/protoc-gen-grpc-gateway-ts/data"
	"google.golang.org/protobuf/types/descriptorpb"
)

func (r *Registry) analyseMessage(fileData *data.File, packageName, fileName string, parents []string, message *descriptorpb.DescriptorProto) {
	parentsPrefix := r.getParentPrefixes(parents)
	packageIdentifier := r.getNameOfPackageLevelIdentifier(parents, message.GetName())
	fqName := "." + packageName + "." + parentsPrefix + message.GetName()
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
					typeInfo.KeyType = &MapEntryType{
						Type:       r.getFieldType(f),
						IsExternal: r.isExternalDependencies(f.GetTypeName(), packageName),
					}
				case "value":
					typeInfo.ValueType = &MapEntryType{
						Type:       r.getFieldType(f),
						IsExternal: r.isExternalDependencies(f.GetTypeName(), packageName),
					}
				}
			}
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
