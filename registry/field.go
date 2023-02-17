package registry

import (
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/data"
)

// getFieldType generates an intermediate type and leave the rendering logic to choose what to render
func (r *Registry) getFieldType(f *descriptorpb.FieldDescriptorProto) string {
	typeName := ""
	if f.Type != nil {
		switch *f.Type {
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			typeName = f.GetTypeName()
		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			typeName = "string"
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			typeName = "bool"
		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			typeName = "bytes"
		case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
			typeName = "float"
		case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
			typeName = "double"
		case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
			typeName = "fixed32"
		case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
			typeName = "sfixed32"
		case descriptorpb.FieldDescriptorProto_TYPE_INT32:
			typeName = "int32"
		case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
			typeName = "sint32"
		case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
			typeName = "uint32"
		case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
			typeName = "fixed64"
		case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
			typeName = "sfixed64"
		case descriptorpb.FieldDescriptorProto_TYPE_INT64:
			typeName = "int64"
		case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
			typeName = "sint64"
		case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
			typeName = "uint64"
		}
	}

	return typeName
}

func (r *Registry) analyseField(fileData *data.File, msgData *data.Message, packageName string, f *descriptorpb.FieldDescriptorProto) {
	fqTypeName := r.getFieldType(f)

	isExternal := r.isExternalDependenciesOutsidePackage(fqTypeName, packageName)

	fieldData := &data.Field{
		Name:         f.GetName(),
		Type:         fqTypeName,
		IsExternal:   isExternal,
		IsOneOfField: f.OneofIndex != nil,
		Message:      msgData,
	}

	if f.Label != nil {
		if f.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			fieldData.IsRepeated = true
		}
	}

	msgData.Fields = append(msgData.Fields, fieldData)

	if !fieldData.IsOneOfField {
		msgData.NonOneOfFields = append(msgData.NonOneOfFields, fieldData)
	}

	// if it's an external dependencies. store in the file data so that they can be collected when every file's finished
	if isExternal {
		fileData.ExternalDependingTypes = append(fileData.ExternalDependingTypes, fqTypeName)
	}

	// if it's a one of field. register the field data in the group of the same one of index.
	if fieldData.IsOneOfField { // one of field
		index := f.GetOneofIndex()
		fieldData.OneOfIndex = index
		_, ok := msgData.OneOfFieldsGroups[index]
		if !ok {
			msgData.OneOfFieldsGroups[index] = make([]*data.Field, 0)
		}
		msgData.OneOfFieldsGroups[index] = append(msgData.OneOfFieldsGroups[index], fieldData)
	}

	fileData.TrackPackageNonScalarType(fieldData)
}
