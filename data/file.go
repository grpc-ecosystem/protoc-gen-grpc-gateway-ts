package data

import (
	"path/filepath"
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
	Services []*Service
	// Name is the name of the file
	Name string
	// TSFileName is the name of the output file
	TSFileName string
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
	name := fileName[0 : len(baseName)-len(ext)]
	return strings.ReplaceAll(packageName, ".", "") + name
}

// GetTSFileName gets the typescript filename out of the proto file name
func GetTSFileName(fileName string) string {
	baseName := filepath.Base(fileName)
	ext := filepath.Ext(fileName)
	name := fileName[0 : len(baseName)-len(ext)]
	return name + ".ts"
}

// Type is an interface to get type out of field and method arguments
type Type interface {
	// GetType returns some information of the type to aid the rendering
	GetType() *TypeInfo
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
