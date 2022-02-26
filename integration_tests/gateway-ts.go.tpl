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