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

var _ machinery.Template = &WebhookService{}

// WebhookService scaffolds a file that defines the webhook service
type WebhookService struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin
	Force bool
}

// SetTemplateDefaults implements file.Template
func (f *WebhookService) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", f.ProjectName, "templates", "webhook-service.yaml")
	}
	f.SetDelim("[[", "]]")
	f.TemplateBody = webhookServiceTemplate

	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		// If file exists (ex. because a webhook was already created), skip creation.
		f.IfExistsAction = machinery.SkipFile
	}
	return nil
}

const webhookServiceTemplate = `{{- if include "[[ .ProjectName ]].webhookEnabled" . -}}
apiVersion: v1
kind: Service
metadata:
  name: {{ include "[[ .ProjectName ]].fullname" . }}-webhook-service
  labels:
    {{- include "[[ .ProjectName ]].labels" . | nindent 4 }}
spec:
  ports:
    - port: 443
      targetPort: 9443
      protocol: TCP
      name: webhook
  selector:
    {{- include "[[ .ProjectName ]].selectorLabels" . | nindent 4 }}
{{- end }}
`
