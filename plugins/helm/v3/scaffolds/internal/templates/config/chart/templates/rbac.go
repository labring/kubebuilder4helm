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

var _ machinery.Template = &Rbac{}

// Rbac scaffolds a file that defines the service account the manager is deployed in.
type Rbac struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin

	Force bool
}

// SetTemplateDefaults implements file.Template
func (f *Rbac) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", f.ProjectName, "templates", "rbac.yaml")
	}

	f.TemplateBody = rbacTemplate
	f.SetDelim("[[", "]]")
	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		// If file exists (ex. because a monitor was already created), skip creation.
		f.IfExistsAction = machinery.SkipFile
	}
	return nil
}

const rbacTemplate = `apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "[[ .ProjectName ]].fullname" . }}
  labels:
    {{- include "[[ .ProjectName ]].labels" . | nindent 4 }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ include "[[ .ProjectName ]].fullname" . }}-role
  labels:
    {{- include "[[ .ProjectName ]].labels" . | nindent 4 }}
rules:
# Add leases roles.
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - coordination.k8s.io
  resources:
  - leases
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{ include "[[ .ProjectName ]].fullname" . }}-cluster-role
  labels:
    {{- include "[[ .ProjectName ]].labels" . | nindent 4 }}
rules:
# Add proxy roles.
- nonResourceURLs:
  - /metrics
  verbs:
  - get
- apiGroups:
  - authentication.k8s.io
  resources:
  - tokenreviews
  verbs:
  - create
- apiGroups:
  - authorization.k8s.io
  resources:
  - subjectaccessreviews
  verbs:
  - create
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
  name: {{ include "[[ .ProjectName ]].fullname" . }}-cluster-role
subjects:
- kind: ServiceAccount
  name: {{ include "[[ .ProjectName ]].fullname" . }}
  namespace: {{ .Release.Namespace }}
`
