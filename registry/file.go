package registry

import (
	"path/filepath"
	"strings"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/data"
	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/options"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus" // nolint: depguard
	"google.golang.org/protobuf/proto"
)

func (r *Registry) analyseFile(f *descriptorpb.FileDescriptorProto) (*data.File, error) {
	log.Debugf("analysing %s", f.GetName())
	fileData := data.NewFile()
	fileName := f.GetName()
	packageName := f.GetPackage()
	parents := make([]string, 0)
	fileData.Name = fileName
	fileData.TSFileName = data.GetTSFileName(fileName)
	if proto.HasExtension(f.Options, options.E_TsPackage) {
		r.TSPackages[fileData.TSFileName] = proto.GetExtension(f.Options, options.E_TsPackage).(string)
	}

	// analyse enums
	for _, enum := range f.EnumType {
		r.analyseEnumType(fileData, packageName, fileName, parents, enum)
	}

	// analyse messages, each message will go recursively
	for _, message := range f.MessageType {
		r.analyseMessage(fileData, packageName, fileName, parents, message)
	}

	// analyse services
	for _, service := range f.Service {
		r.analyseService(fileData, packageName, fileName, service)
	}

	// add fetch module after analysed all services in the file. will add dependencies if there is any
	err := r.addFetchModuleDependencies(fileData)
	if err != nil {
		return nil, errors.Wrapf(err, "error adding fetch module for file %s", fileData.Name)
	}

	r.analyseFilePackageTypeDependencies(fileData)

	return fileData, nil
}

func (r *Registry) addFetchModuleDependencies(fileData *data.File) error {
	if !fileData.Services.NeedsFetchModule() {
		log.Debugf("no services found for %s, skipping fetch module", fileData.Name)
		return nil
	}

	absDir, err := filepath.Abs(r.FetchModuleDirectory)
	if err != nil {
		return errors.Wrapf(err, "error looking up absolute path for fetch module directory %s", r.FetchModuleDirectory)
	}

	foundAtRoot, alias, err := r.findRootAliasForPath(func(absRoot string) (bool, error) {
		return strings.HasPrefix(absDir, absRoot), nil

	})
	if err != nil {
		return errors.Wrapf(err, "error looking up root alias for fetch module directory %s", r.FetchModuleDirectory)
	}

	fileName := filepath.Join(r.FetchModuleDirectory, r.FetchModuleFilename)

	sourceFile, err := r.getSourceFileForImport(fileData.TSFileName, fileName, foundAtRoot, alias)
	if err != nil {
		return errors.Wrapf(err, "error replacing source file with alias for %s", fileName)
	}

	log.Debugf("added fetch dependency %s for %s", sourceFile, fileData.TSFileName)
	fileData.Dependencies = append(fileData.Dependencies, &data.Dependency{
		ModuleIdentifier: "fm",
		SourceFile:       sourceFile,
	})

	return nil
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
