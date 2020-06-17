package registry

import (
	"git.sqcorp.co/cash/gap/cmd/protoc-gen-grpc-gateway-ts/data"
	"git.sqcorp.co/cash/gap/errors"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	log "github.com/sirupsen/logrus" // nolint: depguard
	"google.golang.org/protobuf/types/descriptorpb"
	"path/filepath"
	"strings"
)

// Registry analyse generation request, spits out the data the the rendering process
// it also holds the information about all the types
type Registry struct {
	// Types stores the type information keyed by the fully qualified name of a type
	Types map[string]*TypeInformation

	// FilesToGenerate contains a list of actual file to generate, different from all the files from the request, some of which are import files
	FilesToGenerate map[string]bool
}

// NewRegistry initialise the registry and return the instance
func NewRegistry() *Registry {
	return &Registry{
		Types: make(map[string]*TypeInformation),
	}
}

// TypeInformation store the information about a given type
type TypeInformation struct {
	// Fully qualified name of the type, it starts with `.` and followed by packages and the nested structure path.
	FullyQualifiedName string
	// Package is the package of the type it belongs to
	Package string
	// Files is the file of the type belongs to, this is important in Typescript as modules is the namespace for types defined inside
	File string
	// ModuleIdentifier is the identifier of the type inside the package, this will be useful for enum and nested enum.
	PackageIdentifier string
	// LocalIdentifier is the identifier inside the types local scope
	LocalIdentifier string
	// ProtoType is the type inside the proto. This is used to tell whether it's an enum or a message
	ProtoType descriptorpb.FieldDescriptorProto_Type
	// IsMapEntry indicates whether this type is a Map Entry
	IsMapEntry bool
	// KeyType is the type information for the map key
	KeyType *data.MapEntryType
	// Value type is the type information for the map value
	ValueType *data.MapEntryType
}

// IsFileToGenerate contains the file to be generated in the request
func (r *Registry) IsFileToGenerate(name string) bool {
	result, ok := r.FilesToGenerate[name]
	return ok && result
}

// Analyse analyses the the file inputs, stores types information and spits out the rendering data
func (r *Registry) Analyse(req *plugin.CodeGeneratorRequest) (map[string]*data.File, error) {
	r.FilesToGenerate = make(map[string]bool)
	for _, f := range req.GetFileToGenerate() {
		r.FilesToGenerate[f] = true
	}

	files := req.GetProtoFile()
	log.Debugf("about to start anaylyse files, %d in total", len(files))
	data := make(map[string]*data.File)
	// analyse all files in the request first
	for _, f := range files {
		fileData := r.analyseFile(f)
		data[f.GetName()] = fileData
	}

	// when finishes we have a full map of types and where they are located
	// collect all the external dependencies and back fill it to the file data.
	err := r.collectExternalDependenciesFromData(data)
	if err != nil {
		return nil, errors.Wrap(err, "error collecting external dependency information after analysis finished")
	}

	return data, nil
}

// This simply just concats the parents name and the entity name.
func (r *Registry) getNameOfPackageLevelIdentifier(parents []string, name string) string {
	return strings.Join(parents, "") + name
}

func (r *Registry) getParentPrefixes(parents []string) string {
	parentsPrefix := ""
	if len(parents) > 0 {
		parentsPrefix = strings.Join(parents, ".") + "."
	}
	return parentsPrefix
}

func (r *Registry) isExternalDependenciesOutsidePackage(fqTypeName, packageName string) bool {
	return strings.Index(fqTypeName, "."+packageName) != 0 && strings.Index(fqTypeName, ".") == 0
}

func (r *Registry) collectExternalDependenciesFromData(filesData map[string]*data.File) error {
	for _, fileData := range filesData {
		// dependency group up the dependency by package+file
		dependencies := make(map[string]*data.Dependency)
		for _, typeName := range fileData.ExternalDependingTypes {
			typeInfo, ok := r.Types[typeName]
			if !ok {
				return errors.Errorf("cannot find type info for %s, $v", typeName)
			}
			identifier := typeInfo.Package + "|" + typeInfo.File

			if _, ok := dependencies[identifier]; !ok {
				// only fill in if this file has not been mentioned before.
				// the way import in the genrated file works is like
				// import * as [ModuleIdentifier] from '[Source File]'
				// so there only needs to be added once.
				// Referencing types will be [ModuleIdentifier].[PackageIdentifier]
				base := fileData.TSFileName
				target := data.GetTSFileName(typeInfo.File)
				sourceFile := ""
				var err error
				if !r.IsFileToGenerate(typeInfo.File) {
					sourceFile = "gap/protos/" + target

				} else {
					sourceFile, err = filepath.Rel(filepath.Dir(base), target)
					if err != nil {
						return errors.Wrapf(err, "error getting relative path between for %s, %s", base, target)
					}
					slashSourceFile := filepath.ToSlash(sourceFile)
					if strings.Index(slashSourceFile, "../") != 0 {
						slashSourceFile = "./" + slashSourceFile
					}

					sourceFile = filepath.FromSlash(slashSourceFile)
				}

				// remove ts suffix
				suffixIndex := strings.LastIndex(sourceFile, ".ts")
				sourceFile = sourceFile[0:suffixIndex]

				dependencies[identifier] = &data.Dependency{
					ModuleIdentifier: data.GetModuleName(typeInfo.Package, typeInfo.File),
					SourceFile:       sourceFile,
				}
			}
		}

		for _, dependency := range dependencies {
			fileData.Dependencies = append(fileData.Dependencies, dependency)
		}
	}

	return nil
}
