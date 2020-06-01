package generator

import (
	"bytes"
	"fmt"
	"git.sqcorp.co/cash/gap/cmd/protoc-gen-grpc-gateway-ts/data"
	"git.sqcorp.co/cash/gap/cmd/protoc-gen-grpc-gateway-ts/registry"
	"github.com/Masterminds/sprig"
	"strings"
	"text/template"
)

const tmpl = `
{{define "dependencies"}}{{range .}}import * as {{.ModuleIdentifier}} from "{{.SourceFile}}"
{{end}}{{end}}

{{define "enums"}}
{{range .}}export enum {{.Name}} {
{{range .Values}}
  {{.}} = "{{.}}",
{{end}}
}

{{end}}{{end}}

{{define "messages"}}{{range .}}export interface {{.Name}} {
{{range .Fields}}{{if .IsOneOfField}}
  /**
  * {{.Name}} is in the one of field {{index .Message.OneOfFieldsNames .OneOfIndex}}'s fields: {{range (index .Message.OneOfFieldsGroups .OneOfIndex)}}{{.Name}}, {{end}}
  */{{end}}
  {{.Name}}?: {{tsType .}}
{{end}}
}

{{end}}{{end}}

{{define "services"}}{{range .}}export class {{.Name}} {
{{range .Methods}}  
  static {{.Name}}(req: {{tsType .Input}}): Promise<gap.FetchState<{{tsType .Output}}>> {
    return gap.gapFetchGRPC<FetchLogRequest, FetchLogResponse>("{{.URL}}", req)
  }
{{end}}

}

{{end}}{{end}}

/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

{{if gt (len .Dependencies) 0}}{{include "dependencies" .Dependencies}}{{end}}
{{if gt (len .Enums) 0}}{{include "enums" .Enums}}{{end}}
{{if gt (len .Messages) 0}}{{include "messages" .Messages}}{{end}}
{{if gt (len .Services) 0}}{{include "services" .Services}}{{end}}
`

// GetTemplate gets the templates to for the typescript file
func GetTemplate(r *registry.Registry) *template.Template {
	t := template.New("file")
	t = t.Funcs(sprig.TxtFuncMap())

	t = t.Funcs(template.FuncMap{
		"include": include(t),
		"tsType": func(fieldType data.Type) string {
			return tsType(r, fieldType)
		},
	})

	t = template.Must(t.Parse(tmpl))
	return t
}

// include is the include template functions copied from
// copied from: https://github.com/helm/helm/blob/8648ccf5d35d682dcd5f7a9c2082f0aaf071e817/pkg/engine/engine.go#L147-L154
func include(t *template.Template) func(name string, data interface{}) (string, error) {
	return func(name string, data interface{}) (string, error) {
		buf := bytes.NewBufferString("")
		if err := t.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
}

func tsType(r *registry.Registry, fieldType data.Type) string {
	info := fieldType.GetType()
	if strings.Index(info.Type, ".") != 0 {
		scalaType := mapScalaType(info.Type)
		if info.IsRepeated {
			scalaType += "[]"
		}
		return scalaType
	}
	typeInfo := r.Types[info.Type]
	if typeInfo.IsMapEntry {
		keyType := tsType(r, typeInfo.KeyType)
		valueType := tsType(r, typeInfo.ValueType)

		return fmt.Sprintf("{[key: %s]: %s}", keyType, valueType)

	}
	if !info.IsExternal {
		return typeInfo.PackageIdentifier
	}

	return data.GetModuleName(typeInfo.Package, typeInfo.File) + "." + typeInfo.PackageIdentifier
}

func mapScalaType(protoType string) string {
	switch protoType {
	case "uint64", "sint64", "int64", "fixed64", "sfixed64", "string":
		return "string"
	case "float", "double", "int32", "sint32", "uint32", "fixed32", "sfixed32":
		return "number"
	case "bool":
		return "boolean"
	case "bytes":
		return "Uint8Array"
	}

	return ""

}
