/*
Copyright 2022 The Kubernetes Authors.

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

var _ machinery.Template = &RbacCR{}

// RbacCR scaffolds a file that defines the service account the manager is deployed in.
type RbacCR struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin
	machinery.ResourceMixin
	Force bool
}

// SetTemplateDefaults implements file.Template
func (f *RbacCR) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", f.ProjectName, "templates", "rbac_%[group]_%[kind].yaml")
	}
	f.Path = f.Resource.Replacer().Replace(f.Path)

	f.TemplateBody = rbacCRTemplate
	f.SetDelim("[[", "]]")
	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		// If file exists (ex. because a monitor was already created), skip creation.
		f.IfExistsAction = machinery.SkipFile
	}
	return nil
}

const rbacCRTemplate = `
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{ include "[[ .ProjectName ]].fullname" . }}-[[ lower .Resource.Kind ]]
  labels:
    {{- include "[[ .ProjectName ]].labels" . | nindent 4 }}
rules:
# Add CR roles.
- apiGroups:
  - [[ .Resource.QualifiedGroup ]]
  resources:
  - [[ .Resource.Plural ]]
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - [[ .Resource.QualifiedGroup ]]
  resources:
  - [[ .Resource.Plural ]]/status
  verbs:
  - get
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ include "[[ .ProjectName ]].fullname" . }}-rolebinding
  namespace: {{ .Release.Namespace }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ include "[[ .ProjectName ]].fullname" . }}-role
subjects:
- kind: ServiceAccount
  name: {{ include "[[ .ProjectName ]].fullname" . }}
  namespace: {{ .Release.Namespace }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ include "[[ .ProjectName ]].fullname" . }}-clusterrolebinding
  namespace: {{ .Release.Namespace }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: {{ include "[[ .ProjectName ]].fullname" . }}-[[ lower .Resource.Kind ]]
subjects:
- kind: ServiceAccount
  name: {{ include "[[ .ProjectName ]].fullname" . }}
  namespace: {{ .Release.Namespace }}
`
