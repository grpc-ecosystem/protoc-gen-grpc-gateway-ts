package registry

import (
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	log "github.com/sirupsen/logrus" // nolint: depguard

	"github.com/squareup/gap/cmd/protoc-gen-grpc-gateway-ts/data"
)

func (r *Registry) analyseFile(f *descriptorpb.FileDescriptorProto) *data.File {
	log.Debugf("analysing %s", f.GetName())
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

	r.analyseFilePackageTypeDependencies(fileData)

	return fileData
}

func (r *Registry) analyseFilePackageTypeDependencies(fileData *data.File) {
	for _, t := range fileData.PackageNonScalarType {
		// for each non scalar types try to determine if the type comes from same
		// package but a different file. if yes then will need to add the type to
		// the external dependencies for collection later
		// also need to change the type's IsExternal information for rendering purpose
		typeInfo := t.GetType()
		fqTypeName := typeInfo.Type
		log.Debugf("checking whether non scala type %s in the same message is external to the current file", fqTypeName)

		registryType, foundInRegistry := r.Types[fqTypeName]
		if !foundInRegistry || registryType.File != fileData.Name {
			// this means the type from same package in file has yet to be analysed (means in different file)
			// or the type has appeared in another file different to the current file
			// in this case we will put the type as external in the fileData
			// and also mutate the IsExternal field of the given type:w
			log.Debugf("type %s is external to file %s, mutating the external dependencies information", fqTypeName, fileData.Name)

			fileData.ExternalDependingTypes = append(fileData.ExternalDependingTypes, fqTypeName)
			t.SetExternal(true)
		}
	}
}
