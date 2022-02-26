package generator

import (
	"bytes"
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"github.com/Masterminds/sprig"
	"github.com/iancoleman/strcase"

	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/data"
	"github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/registry"
)

// GetTemplate gets the templates to for the typescript file
func GetTemplate(r *registry.Registry, tmplPath string) *template.Template {
	t := template.New(path.Base(tmplPath))
	t = t.Funcs(sprig.TxtFuncMap())

	t = t.Funcs(template.FuncMap{
		"include": include(t),
		"tsType": func(fieldType data.Type) string {
			return tsType(r, fieldType)
		},
		"renderURL":     renderURL(r),
		"buildInitReq":  buildInitReq,
		"fieldName":     fieldName(r),
		"jsonFieldName": jsonFieldName(r),
		"isNotWellKnownDeps": func(dep string) bool {
			return !isWellKnownDeps(dep)
		},
	})

	return template.Must(t.ParseFiles(tmplPath))
}

func jsonFieldName(r *registry.Registry) func(name data.Field) string {
	re := regexp.MustCompile("^[a-zA-Z_$][a-zA-Z_$0-9]*$")

	return func(name data.Field) string {
		if re.Match([]byte(name.JSONName)) {
			return name.JSONName
		}
		return fmt.Sprintf(`["%s"]`, name.JSONName)
	}
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
