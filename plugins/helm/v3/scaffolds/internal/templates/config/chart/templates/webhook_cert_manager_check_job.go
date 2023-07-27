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
	"path/filepath"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
)

var _ machinery.Template = &WebhookCertManagerCheck{}

// WebhookCertManagerCheck scaffolds a file that defines the helm scheme for the helpers.
type WebhookCertManagerCheck struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin
	machinery.RepositoryMixin
	Force bool
}

// SetTemplateDefaults implements file.Template
func (f *WebhookCertManagerCheck) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", f.ProjectName, "templates", "webhook-cert-manager-check.yaml")
	}
	f.SetDelim("[[", "]]")
	f.TemplateBody = certManagerCheckTemplate

	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		// If file exists (ex. because a webhook was already created), skip creation.
		f.IfExistsAction = machinery.SkipFile
	}

	return nil
}

const certManagerCheckTemplate = `{{- if include "[[ .ProjectName ]].webhookEnabled" . -}}
apiVersion: batch/v1
kind: Job
metadata:
  name: "{{ .Release.Name }}-cert-manager-check"
  annotations:
    "helm.sh/hook": pre-install,pre-upgrade
    "helm.sh/hook-weight": "-5"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    metadata:
      name: "{{ .Release.Name }}-cert-manager-check"
      labels:
        app: "{{ .Chart.Name }}"
        release: "{{ .Release.Name }}"
    spec:
      restartPolicy: Never
      containers:
        - name: cert-manager-check
          image: busybox:latest
          command: ["sh", "-c", "until echo exit | telnet {{ .Values.certManager.domain }} {{ .Values.certManager.port }}; do echo waiting for cert-manager; sleep 10; done;"]
{{- end }}
`
