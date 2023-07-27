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

package chart

import (
	"path/filepath"
	"strings"

	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
)

var _ machinery.Template = &Values{}

// Values scaffolds a file that defines the kustomization scheme for the webhook folder
type Values struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin
	machinery.RepositoryMixin
	Force            bool
	GithubDockerRepo string
}

// SetTemplateDefaults implements file.Template
func (f *Values) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", f.ProjectName, "values.yaml")
	}

	f.GithubDockerRepo = strings.Join(strings.Split(f.Repo, "/")[:2], "/")

	f.TemplateBody = valuesTemplate

	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		// If file exists (ex. because a webhook was already created), skip creation.
		f.IfExistsAction = machinery.SkipFile
	}

	return nil
}

const valuesTemplate = `# Default values for {{ .ProjectName }}.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.
replicaCount: 1
nameOverride: ""
fullnameOverride: ""
main:
  image:
    repository: {{ .GithubDockerRepo }}/{{ .ProjectName }}
    pullPolicy: IfNotPresent
    tag: "latest"
  resources: 
    limits:
      cpu: 100m
      memory: 128Mi
    requests:
      cpu: 100m
      memory: 128Mi
  securityContext:
    allowPrivilegeEscalation: false
    capabilities:
      drop:
        - "ALL"

proxy:
  image:
    repository: gcr.io/kubebuilder/kube-rbac-proxy
    pullPolicy: IfNotPresent
    tag: "v0.13.0"
  resources: 
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 5m
      memory: 64Mi
  securityContext:
    allowPrivilegeEscalation: false
    capabilities:
      drop:
        - "ALL"

prometheus: false

certManager:
  domain: cert-manager-webhook.cert-manager.svc
  port: 443

logger:
  #Development Mode
  zap: true
  #  Can be one of 'debug', 'info', 'error'
  level: info

rateLimiter:
  minRetryDelay: 5ms
  maxRetryDelay: 1000s
  defaultQPS: 10.0
  defaultBurst: 100
  defaultConcurrent: 5


podAnnotations: {}

nodeSelector: {}

tolerations: []

affinity: {}
`
