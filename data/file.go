package data

import (
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// File store the information about rendering a file
type File struct {
	// Dependencies is a list of dependencies for the file, which will be rendered at the top of the file as import statements
	Dependencies []*Dependency
	// Enums is a list of enums to render, due to the fact that there cannot be any enum defined nested in the class in Typescript.
	// All Enums will be rendered at the top level
	Enums []*Enum
	// Messages represents top level messages inside the file.
	Messages []*Message
	// ExternalDependingTypes stores the external dependenciees fully qualified name,
	ExternalDependingTypes []string
	// Services stores the information to render service
	Services Services
	// Name is the name of the file
	Name string
	// TSFileName is the name of the output file
	TSFileName string
	// PackageNonScalarType stores the type inside the same packages within the file, which will be used to figure out external dependencies inside the same package (different files)
	PackageNonScalarType []Type
	// EnableStylingCheck enables the styling check for the given file
	EnableStylingCheck bool
}

// StableDependencies are dependencies in a stable order.
func (f *File) StableDependencies() []*Dependency {
	out := make([]*Dependency, len(f.Dependencies))
	copy(out, f.Dependencies)
	sort.Slice(out, func(i, j int) bool {
		return out[i].SourceFile < out[j].SourceFile
	})
	return out
}

// NeedsOneOfSupport indicates the file needs one of support type utilities
func (f *File) NeedsOneOfSupport() bool {
	for _, m := range f.Messages {
		if m.HasOneOfFields() {
			return true
		}
	}

	return false
}

// TrackPackageNonScalarType tracks the supplied non scala type in the same package
func (f *File) TrackPackageNonScalarType(t Type) {
	isNonScalarType := strings.Index(t.GetType().Type, ".") == 0
	if isNonScalarType {
		f.PackageNonScalarType = append(f.PackageNonScalarType, t)
	}
}

func (f *File) IsEmpty() bool {
	return len(f.Enums) == 0 && len(f.Messages) == 0 && len(f.Services) == 0
}

// NewFile returns an initialised new file
func NewFile() *File {
	return &File{
		Dependencies:           make([]*Dependency, 0),
		Enums:                  make([]*Enum, 0),
		Messages:               make([]*Message, 0),
		Services:               make([]*Service, 0),
		ExternalDependingTypes: make([]string, 0),
	}

}

// Dependency stores the information about dependencies.
type Dependency struct {
	// ModuleIdentifier will be a concanation of package + file base name to make it
	// unnique inside the file. This will act as a name space for other file and
	// types inside other file can be referred to using .
	ModuleIdentifier string
	// Source file will be the file at the end of the import statement.
	SourceFile string
}

// GetModuleName returns module name = package name + file name to be the unique identifier for source file in a ts file
func GetModuleName(packageName, fileName string) string {
	baseName := filepath.Base(fileName)
	ext := filepath.Ext(fileName)
	name := baseName[0 : len(baseName)-len(ext)]
	packageParts := strings.Split(packageName, ".")

	if packageName != "" {
		for i, p := range packageParts {
			packageParts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}

	return strings.Join(packageParts, "") + strings.ToUpper(name[:1]) + name[1:]
}

// GetTSFileName gets the typescript filename out of the proto file name
func GetTSFileName(fileName string) string {
	baseName := filepath.Base(fileName)
	ext := filepath.Ext(fileName)
	name := baseName[0 : len(baseName)-len(ext)]
	return path.Join(filepath.Dir(fileName), name+".pb.ts")
}

// Type is an interface to get type out of field and method arguments
type Type interface {
	// GetType returns some information of the type to aid the rendering
	GetType() *TypeInfo
	// SetExternal changes the external field inside the data structure
	SetExternal(bool)
}

// TypeInfo stores some common type information for rendering
type TypeInfo struct {
	// Type
	Type string
	// IsRepeated indicates whether this field is a repeated field
	IsRepeated bool
	// IsExternal indicates whether this type is external
	IsExternal bool
}
