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

package scaffolds

import (
	"fmt"

	"github.com/labring/kubebuilder4helm/plugins/helm/v3/scaffolds/internal/templates/config/chart"
	templates2 "github.com/labring/kubebuilder4helm/plugins/helm/v3/scaffolds/internal/templates/config/chart/templates"
	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugins"
)

const (
	imageName = "controller:latest"
)

var _ plugins.Scaffolder = &initScaffolder{}

type initScaffolder struct {
	config config.Config
	// fs is the filesystem that will be used by the scaffolder
	fs machinery.Filesystem
}

// NewInitScaffolder returns a new Scaffolder for project initialization operations
func NewInitScaffolder(config config.Config) plugins.Scaffolder {
	return &initScaffolder{
		config: config,
	}
}

// InjectFS implements cmdutil.Scaffolder
func (s *initScaffolder) InjectFS(fs machinery.Filesystem) {
	s.fs = fs
}

// Scaffold implements cmdutil.Scaffolder
func (s *initScaffolder) Scaffold() error {
	fmt.Println("Writing helm manifests for you to edit...")

	// Initialize the machinery.Scaffold that will write the files to disk
	scaffold := machinery.NewScaffold(s.fs,
		machinery.WithConfig(s.config),
	)

	templates := []machinery.Builder{
		//&rbac2.Kustomization{},
		//&rbac2.AuthProxyRole{},
		//&rbac2.AuthProxyRoleBinding{},
		//&rbac2.AuthProxyService{},
		//&rbac2.AuthProxyClientRole{},
		//&rbac2.RoleBinding{},
		//&rbac2.LeaderElectionRole{},
		//&rbac2.LeaderElectionRoleBinding{},
		//&rbac2.ServiceAccount{},
		//&manager2.Kustomization{},
		//&manager2.Config{Image: imageName},
		//&kdefault2.Kustomization{},
		//&kdefault2.ManagerAuthProxyPatch{},
		//&kdefault2.ManagerConfigPatch{},
		//&prometheus2.Kustomization{},
		//&prometheus2.Monitor{},
		&chart.Chart{},
		&chart.HelmIgnore{},
		&chart.Values{},
		&templates2.Helpers{},
		&templates2.MonitorService{Force: true},
		&templates2.Monitor{Force: true},
		&templates2.Rbac{Force: true},
		&templates2.Deployment{Force: true},
	}

	return scaffold.Execute(templates...)
}
