package generator

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/Masterminds/sprig"
	"github.com/iancoleman/strcase"

	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/data"
	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/registry"
)

const tmpl = `
{{define "dependencies"}}
import { Observable } from 'rxjs';
{{ range . -}}
{{- if isNotWellKnownDeps .SourceFile -}}
import * as {{.ModuleIdentifier}} from "{{.SourceFile}}"
{{ end }}
{{- end }}
{{end}}

{{define "enums"}}
{{range .}}export enum {{.Name}} {
{{- range .Values}}
  {{.}} = "{{.}}",
{{- end}}
}

{{end}}{{end}}

{{define "messages"}}{{range .}}
{{- if .HasOneOfFields}}
type Base{{.Name}} = {
{{- range .NonOneOfFields}}
  {{fieldName .Name}}?: {{tsType .}};
{{- end}}
}

export type {{.Name}} = Base{{.Name}}
{{range $groupId, $fields := .OneOfFieldsGroups}}  & OneOf<{ {{range $index, $field := $fields}}{{fieldName $field.Name}}: {{tsType $field}}{{if (lt (add $index 1) (len $fields))}}; {{end}}{{end}} }>
{{end}}
{{- else -}}
export type {{.Name}} = {
{{- range .Fields}}
  {{fieldName .Name}}?: {{tsType .}};
{{- end}}
}
{{end}}
{{end}}{{end}}

{{define "services"}}{{range .}}export class {{.Name}} {
{{- range .Methods}}  
{{- if .ServerStreaming }}
  static {{.Name}}(req: {{tsType .Input}}, entityNotifier?: fm.NotifyStreamEntityArrival<{{tsType .Output}}>, initReq?: fm.InitReq): Promise<void> {
    return fm.fetchStreamingRequest<{{tsType .Input}}, {{tsType .Output}}>(` + "`{{renderURL .}}`" + `, entityNotifier, {...initReq, {{buildInitReq .}}});
  }
{{- else }}
  static {{.Name}}(req: {{tsType .Input}}, initReq?: fm.InitReq): Promise<{{tsType .Output}}> {
    return fm.fetchReq<{{tsType .Input}}, {{tsType .Output}}>(` + "`{{renderURL .}}`" + `, {...initReq, {{buildInitReq .}}});
  }
{{- end}}
{{- end}}
}
{{end}}{{end}}

{{define "observableServices"}}{{range .}}export class Observable{{.Name}} {
{{- range .Methods}}  
{{- if .ServerStreaming }}
  static {{.Name}}(req: {{tsType .Input}}, initReq?: fm.InitReq): Observable<{{tsType .Output}}> {
    return fm.fromFetchStreamingRequest<{{tsType .Input}}, {{tsType .Output}}>(` + "`{{renderURL .}}`" + `, {...initReq, {{buildInitReq .}}});
  }
{{- else }}
  static {{.Name}}(req: {{tsType .Input}}, initReq?: fm.InitReq): Observable<{{tsType .Output}}> {
    return fm.fromFetchReq<{{tsType .Input}}, {{tsType .Output}}>(` + "`{{renderURL .}}`" + `, {...initReq, {{buildInitReq .}}});
  }
{{- end}}
{{- end}}
}
{{end}}{{end}}

{{- if not .EnableStylingCheck}}
/* eslint-disable */
// @ts-nocheck
{{- end}}
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
{{if .Dependencies}}{{- include "dependencies" .StableDependencies -}}{{end}}
{{- if .NeedsOneOfSupport}}
type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };
type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (
    keyof T extends infer K ?
      (K extends string & keyof T ? { [k in K]: T[K] } & Absent<T, K>
        : never)
    : never);
{{end}}
{{- if .Enums}}{{include "enums" .Enums}}{{end}}
{{- if .Messages}}{{include "messages" .Messages}}{{end}}
{{- if .Services -}}
{{ include "services" .Services }}
{{ include "observableServices" .Services }}
{{- end }}
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
		"renderURL":    renderURL(r),
		"buildInitReq": buildInitReq,
		"fieldName":    fieldName(r),
		"isNotWellKnownDeps": func(dep string) bool {
			return !isWellKnownDeps(dep)
		},
	})

	t = template.Must(t.Parse(tmpl))
	return t
}

func fieldName(r *registry.Registry) func(name string) string {
	return func(name string) string {
		if r.UseProtoNames {
			return name
		}

		return strcase.ToLowerCamel(name)
	}
}

type fieldPartMapFn func(part string) string

func splitFieldName(fieldName string) []string {
	return strings.Split(fieldName, ".")
}

func mapFieldName(parts []string, fieldPartMapFns ...fieldPartMapFn) []string {
	if len(fieldPartMapFns) == 0 {
		return parts
	}
	newParts := make([]string, 0, len(parts))
	for _, fieldName := range parts {
		newParts = append(newParts, fieldPartMapFns[0](fieldName))
	}
	return mapFieldName(newParts, fieldPartMapFns[1:]...)
}

func renderIndexRequestField(fieldName string, fieldNameFn fieldPartMapFn) string {
	parts := mapFieldName(splitFieldName(fieldName), fieldNameFn, func(part string) string {
		return fmt.Sprintf(`["%s"]`, part)
	})
	return fmt.Sprintf("req%s", strings.Join(parts, ""))
}

func renderParts(fieldName string, fieldNameFn fieldPartMapFn) string {
	return fmt.Sprintf("${%s}", renderIndexRequestField(fieldName, fieldNameFn))
}

func renderFieldInPath(fieldName string, fieldNameFn fieldPartMapFn) string {
	parts := mapFieldName(splitFieldName(fieldName), fieldNameFn)
	return fmt.Sprintf(`"%s"`, strings.Join(parts, "."))
}

func renderURL(r *registry.Registry) func(method data.Method) string {
	fieldNameFn := fieldName(r)
	return func(method data.Method) string {
		methodURL := method.URL
		reg := regexp.MustCompile("{([^}]+)}")
		matches := reg.FindAllStringSubmatch(methodURL, -1)
		fieldsInPath := make([]string, 0, len(matches))
		if len(matches) > 0 {
			log.Debugf("url matches %v", matches)
			for _, m := range matches {
				expToReplace := m[0]
				// cleanup m[1] if the pattern is {fieldname=resources/*}
				cleanedFieldName := strings.Split(m[1], "=")[0]
				part := renderParts(cleanedFieldName, fieldNameFn)
				methodURL = strings.ReplaceAll(methodURL, expToReplace, part)
				fieldsInPath = append(fieldsInPath, renderFieldInPath(cleanedFieldName, fieldNameFn))
			}
		}
		urlPathParams := fmt.Sprintf("[%s]", strings.Join(fieldsInPath, ", "))

		if !method.ClientStreaming && method.HTTPMethod == "GET" {
			// parse the url to check for query string
			parsedURL, err := url.Parse(methodURL)
			if err != nil {
				return methodURL
			}
			renderURLSearchParamsFn := fmt.Sprintf("${fm.renderURLSearchParams(req, %s)}", urlPathParams)
			// prepend "&" if query string is present otherwise prepend "?"
			// trim leading "&" if present before prepending it
			if parsedURL.RawQuery != "" {
				methodURL = strings.TrimRight(methodURL, "&") + "&" + renderURLSearchParamsFn
			} else {
				methodURL += "?" + renderURLSearchParamsFn
			}
		}

		return methodURL
	}
}

func buildInitReq(method data.Method) string {
	httpMethod := method.HTTPMethod
	m := `method: "` + httpMethod + `"`
	fields := []string{m}
	if method.HTTPRequestBody == nil || *method.HTTPRequestBody == "*" {
		fields = append(fields, "body: JSON.stringify(req, fm.replacer)")
	} else if *method.HTTPRequestBody != "" {
		field := fmt.Sprintf("body: JSON.stringify(%s, fm.replacer)", renderIndexRequestField(*method.HTTPRequestBody, strcase.ToLowerCamel))
		fields = append(fields, field)
	}
	return strings.Join(fields, ", ")
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
	typeInfo, ok := r.Types[info.Type]
	if ok && typeInfo.IsMapEntry {
		keyType := tsType(r, typeInfo.KeyType)
		valueType := tsType(r, typeInfo.ValueType)

		return fmt.Sprintf("{[key: %s]: %s}", keyType, valueType)
	}

	typeStr := mapWellKnownType(info.Type)

	if typeStr == "" {
		if !info.IsExternal {
			typeStr = typeInfo.PackageIdentifier
		} else {
			typeStr = data.GetModuleName(typeInfo.Package, typeInfo.File) + "." + typeInfo.PackageIdentifier
		}
	}

	if info.IsRepeated {
		typeStr += "[]"
	}
	return typeStr
}

func mapWellKnownType(protoType string) string {
	switch protoType {
	case "uint64", "sint64", "int64", "fixed64", "sfixed64", "string":
		return "string"
	case "float", "double", "int32", "sint32", "uint32", "fixed32", "sfixed32":
		return "number"
	case "bool":
		return "boolean"
	case "bytes":
		return "Uint8Array"
	case ".google.protobuf.Timestamp":
		return "string"
	case ".google.protobuf.Duration":
		return "string"
	case ".google.protobuf.Struct":
		return "unknown"
	case ".google.protobuf.Value":
		return "unknown"
	case ".google.protobuf.ListValue":
		return "unknown[]"
	case ".google.protobuf.NullValue":
		return "null"
	case ".google.protobuf.FieldMask":
		return "string[]"
	case ".google.protobuf.Any":
		return "unknown"
	case ".google.protobuf.Empty":
		return "{}"
	}
	return ""
}

var (
	wellKnownDeps = []string{
		"google/protobuf/duration",
		"google/protobuf/timestamp",
		"google/protobuf/struct",
		"google/protobuf/field_mask",
		"google/protobuf/empty",
		"google/protobuf/any",
		"../protobuf/any",
	}
)

func isWellKnownDeps(dep string) bool {
	for _, wellKnownDep := range wellKnownDeps {
		if strings.Contains(dep, wellKnownDep) {
			return true
		}
	}
	return false
}
