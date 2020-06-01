package generator

import (
	"bytes"
	"git.sqcorp.co/cash/gap/cmd/protoc-gen-grpc-gateway-ts/data"
	"git.sqcorp.co/cash/gap/cmd/protoc-gen-grpc-gateway-ts/registry"
	"git.sqcorp.co/cash/gap/errors"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"strings"
	"text/template"
)

// TypeScriptGRPCGatewayGenerator is the protobuf generator for typescript
type TypeScriptGRPCGatewayGenerator struct {
	Registry *registry.Registry
}

// New returns an initialised generator
func New() *TypeScriptGRPCGatewayGenerator {
	return &TypeScriptGRPCGatewayGenerator{
		Registry: registry.NewRegistry(),
	}
}

// Generate take a code generator request and returns a response. it analyse request with registry and use the generated data to render ts files
func (t *TypeScriptGRPCGatewayGenerator) Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	resp := &plugin.CodeGeneratorResponse{}

	filesData, err := t.Registry.Analyse(req.ProtoFile)
	if err != nil {
		return nil, errors.Wrap(err, "error analysing proto files")
	}
	tmpl := GetTemplate(t.Registry)

	// feed fileData into rendering process
	for _, fileData := range filesData {
		generated, err := t.generateFile(fileData, tmpl)
		if err != nil {
			return nil, errors.Wrap(err, "error generating file")
		}
		resp.File = append(resp.File, generated)
	}

	return resp, nil
}

func (t *TypeScriptGRPCGatewayGenerator) generateFile(fileData *data.File, tmpl *template.Template) (*plugin.CodeGeneratorResponse_File, error) {
	w := bytes.NewBufferString("")

	err := tmpl.Execute(w, fileData)
	if err != nil {
		return nil, errors.Wrapf(err, "error generating ts file for %s", fileData.Name)
	}

	fileName := fileData.TSFileName
	content := strings.TrimSpace(w.String())

	return &plugin.CodeGeneratorResponse_File{
		Name:           &fileName,
		InsertionPoint: nil,
		Content:        &content,
	}, nil
}
