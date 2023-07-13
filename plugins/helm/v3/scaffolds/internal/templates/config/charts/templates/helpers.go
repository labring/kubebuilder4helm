/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package templates

import (
	"fmt"
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"text/template"
)

var _ machinery.Template = &Helpers{}

// Helpers scaffolds a file that defines the helm scheme for the helpers.
type Helpers struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin
	machinery.RepositoryMixin
	Force          bool
	WebhookEnabled bool
}

// SetTemplateDefaults implements file.Template
func (f *Helpers) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", f.ProjectName, "templates", "_helpers.tpl")
	}
	f.SetDelim("[[", "]]")

	f.TemplateBody = helpersTemplate
	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		// If file exists (ex. because a webhook was already created), skip creation.
		f.IfExistsAction = machinery.SkipFile
	}

	return nil
}

// GetFuncMap implements file.UseCustomFuncMap
func (f *Helpers) GetFuncMap() template.FuncMap {
	funcMap := machinery.DefaultFuncMap()
	funcMap["JSONTag"] = func(tag string) string {
		return fmt.Sprintf("`json:%q`", tag)
	}
	return funcMap
}

const helpersTemplate = `{{/*
Expand the name of the chart.
*/}}
{{- define "[[ .ProjectName ]].name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "[[ .ProjectName ]].fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "[[ .ProjectName ]].chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "[[ .ProjectName ]].labels" -}}
helm.sh/chart: {{ include "[[ .ProjectName ]].chart" . }}
{{ include "[[ .ProjectName ]].selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "[[ .ProjectName ]].selectorLabels" -}}
app.kubernetes.io/name: {{ include "[[ .ProjectName ]].name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "[[ .ProjectName ]].webhookEnabled" }}
{{- "[[ .WebhookEnabled ]]" }}
{{- end }}
`
