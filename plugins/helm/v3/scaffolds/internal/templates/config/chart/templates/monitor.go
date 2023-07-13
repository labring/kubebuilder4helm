/*
Copyright 2020 The Kubernetes Authors.

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
	"path/filepath"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
)

var _ machinery.Template = &Monitor{}

// Monitor scaffolds a file that defines the prometheus service monitor
type Monitor struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin
	Force bool
}

// SetTemplateDefaults implements file.Template
func (f *Monitor) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", f.ProjectName, "templates", "monitor.yaml")
	}
	f.SetDelim("[[", "]]")

	f.TemplateBody = monitorTemplate

	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		// If file exists (ex. because a monitor was already created), skip creation.
		f.IfExistsAction = machinery.SkipFile
	}
	return nil
}

const monitorTemplate = `{{- if .Values.prometheus -}}
# Prometheus Monitor Service (Metrics)
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: {{ include "[[ .ProjectName ]].fullname" . }}-metrics
  labels:
    {{- include "[[ .ProjectName ]].labels" . | nindent 4 }}
spec:
  endpoints:
    - path: /metrics
      port: https
      scheme: https
      bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
      tlsConfig:
        insecureSkipVerify: true
  selector:
    matchLabels:
      {{- include "[[ .ProjectName ]].selectorLabels" . | nindent 4 }}
{{- end }}
`
