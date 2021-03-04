package generator

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	log "github.com/sirupsen/logrus" // nolint: depguard

	"github.com/pkg/errors"
	"github.com/squareup/protoc-gen-grpc-gateway-ts/data"
	"github.com/squareup/protoc-gen-grpc-gateway-ts/registry"
)

// TypeScriptGRPCGatewayGenerator is the protobuf generator for typescript
type TypeScriptGRPCGatewayGenerator struct {
	Registry *registry.Registry
}

// New returns an initialised generator
func New(paramsMap map[string]string) (*TypeScriptGRPCGatewayGenerator, error) {
	registry, err := registry.NewRegistry(paramsMap)
	if err != nil {
		return nil, errors.Wrap(err, "error instantiating a new registry")
	}

	return &TypeScriptGRPCGatewayGenerator{
		Registry: registry,
	}, nil
}

// Generate take a code generator request and returns a response. it analyse request with registry and use the generated data to render ts files
func (t *TypeScriptGRPCGatewayGenerator) Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	resp := &plugin.CodeGeneratorResponse{}

	filesData, err := t.Registry.Analyse(req)
	if err != nil {
		return nil, errors.Wrap(err, "error analysing proto files")
	}
	tmpl := GetTemplate(t.Registry)
	log.Debugf("files to generate %v", req.GetFileToGenerate())

	needToGenerateFetchModule := false
	// feed fileData into rendering process
	for _, fileData := range filesData {
		if !t.Registry.IsFileToGenerate(fileData.Name) {
			log.Debugf("file %s is not the file to generate, skipping", fileData.Name)
			continue
		}

		log.Debugf("generating file for %s", fileData.TSFileName)
		generated, err := t.generateFile(fileData, tmpl)
		if err != nil {
			return nil, errors.Wrap(err, "error generating file")
		}
		resp.File = append(resp.File, generated)
		needToGenerateFetchModule = needToGenerateFetchModule || fileData.Services.NeedsFetchModule()
	}

	if needToGenerateFetchModule {
		// generate fetch module
		fetchTmpl := GetFetchModuleTemplate()
		log.Debugf("generate fetch template")
		generatedFetch, err := t.generateFetchModule(fetchTmpl)
		if err != nil {
			return nil, errors.Wrap(err, "error generating fetch module")
		}

		resp.File = append(resp.File, generatedFetch)
	}

	return resp, nil
}

func (t *TypeScriptGRPCGatewayGenerator) generateFile(fileData *data.File, tmpl *template.Template) (*plugin.CodeGeneratorResponse_File, error) {
	w := bytes.NewBufferString("")

	if fileData.IsEmpty() {
		w.Write([]byte(fmt.Sprintln("export default {}")))
	} else {
		err := tmpl.Execute(w, fileData)
		if err != nil {
			return nil, errors.Wrapf(err, "error generating ts file for %s", fileData.Name)
		}
	}

	fileName := fileData.TSFileName
	content := strings.TrimSpace(w.String())

	return &plugin.CodeGeneratorResponse_File{
		Name:           &fileName,
		InsertionPoint: nil,
		Content:        &content,
	}, nil
}

func (t *TypeScriptGRPCGatewayGenerator) generateFetchModule(tmpl *template.Template) (*plugin.CodeGeneratorResponse_File, error) {
	w := bytes.NewBufferString("")
	fileName := filepath.Join(t.Registry.FetchModuleDirectory, t.Registry.FetchModuleFilename)
	err := tmpl.Execute(w, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error generating fetch module at %s", fileName)
	}

	content := strings.TrimSpace(w.String())
	return &plugin.CodeGeneratorResponse_File{
		Name:           &fileName,
		InsertionPoint: nil,
		Content:        &content,
	}, nil
}
