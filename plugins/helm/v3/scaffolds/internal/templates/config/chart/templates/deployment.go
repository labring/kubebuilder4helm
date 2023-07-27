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

var _ machinery.Template = &Deployment{}

// Deployment scaffolds a file that defines the service account the manager is deployed in.
type Deployment struct {
	machinery.TemplateMixin
	machinery.ProjectNameMixin

	Force bool
}

// SetTemplateDefaults implements file.Template
func (f *Deployment) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", f.ProjectName, "templates", "deployment.yaml")
	}

	f.TemplateBody = deploymentTemplate
	f.SetDelim("[[", "]]")
	if f.Force {
		f.IfExistsAction = machinery.OverwriteFile
	} else {
		// If file exists (ex. because a monitor was already created), skip creation.
		f.IfExistsAction = machinery.SkipFile
	}
	return nil
}

const deploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "[[ .ProjectName ]].fullname" . }}
  labels:
    {{- include "[[ .ProjectName ]].labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "[[ .ProjectName ]].selectorLabels" . | nindent 6 }}
  template:
    metadata:
      {{- with .Values.podAnnotations }}
      annotations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      labels:
        {{- include "[[ .ProjectName ]].selectorLabels" . | nindent 8 }}
    spec:
      serviceAccountName: {{ include "[[ .ProjectName ]].fullname" . }}
      containers:
        - name: {{ .Chart.Name }}
          command:
          - /manager
          args:
            - --health-probe-bind-address=:8081
            - --metrics-bind-address=127.0.0.1:8080
            - --leader-elect
            - --zap-devel={{ .Values.logger.zap }}
            - --zap-log-level={{ .Values.logger.level }}
            - --default-burst={{ .Values.rateLimiter.defaultBurst }}
            - --default-concurrent={{ .Values.rateLimiter.defaultConcurrent }}
            - --default-qps={{ .Values.rateLimiter.defaultQPS }}
            - --max-retry-delay={{ .Values.rateLimiter.maxRetryDelay }}
            - --min-retry-delay={{ .Values.rateLimiter.minRetryDelay }}
          securityContext:
            {{- toYaml .Values.main.securityContext | nindent 12 }}
          image: "{{ .Values.main.image.repository }}:{{ .Values.main.image.tag | default .Chart.AppVersion }}"
          imagePullPolicy: {{ .Values.main.image.pullPolicy }}
          ports:
          - containerPort: 8081
            name: health
            protocol: TCP
          livenessProbe:
            httpGet:
              path: /healthz
              port: health
          readinessProbe:
            httpGet:
              path: /readyz
              port: health
          resources:
            {{- toYaml .Values.main.resources | nindent 12 }}
          {{- if include "[[ .ProjectName ]].webhookEnabled" . -}}
          volumeMounts:
            - mountPath: /tmp/k8s-webhook-server/serving-certs
              name: cert
              readOnly: true
          {{- end }}
        - name: kube-rbac-proxy
          args:
            - --secure-listen-address=0.0.0.0:8443
            - --upstream=http://127.0.0.1:8080/
            - --logtostderr=true
            - --v=0
          securityContext:
            {{- toYaml .Values.proxy.securityContext | nindent 12 }}
          image: "{{ .Values.proxy.image.repository }}:{{ .Values.proxy.image.tag }}"
          imagePullPolicy: {{ .Values.proxy.image.pullPolicy }}
          ports:
            - containerPort: 8443
              name: https
              protocol: TCP
          resources:
            {{- toYaml .Values.proxy.resources | nindent 12 }}
      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- if include "[[ .ProjectName ]].webhookEnabled" . -}}
      volumes:
        - name: cert
          secret:
            defaultMode: 420
            secretName: {{ include "[[ .ProjectName ]].fullname" . }}-webhook-server-cert
      {{- end }}
`
