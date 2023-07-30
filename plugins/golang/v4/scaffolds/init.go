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

	"github.com/labring/kubebuilder4helm/internal/version"
	"github.com/spf13/afero"

	"github.com/labring/kubebuilder4helm/plugins/golang/v4/scaffolds/internal/templates"
	"github.com/labring/kubebuilder4helm/plugins/golang/v4/scaffolds/internal/templates/hack"
	helmv3 "github.com/labring/kubebuilder4helm/plugins/helm/v3"
	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugins"
)

const (
	// ControllerRuntimeVersion is the kubernetes-sigs/controller-runtime version to be used in the project
	ControllerRuntimeVersion = "v0.15.0"
	// ControllerToolsVersion is the kubernetes-sigs/controller-tools version to be used in the project
	ControllerToolsVersion = "v0.12.0"
	// EndpointOperatorLibVersion is the labring/operator-sdk version to be used in the project
	EndpointOperatorLibVersion = "v1.0.1"

	imageName = "controller:latest"
)

var _ plugins.Scaffolder = &initScaffolder{}

var helmVersion string

type initScaffolder struct {
	config          config.Config
	boilerplatePath string
	license         string
	owner           string
	isLegacyLayout  bool
	// fs is the filesystem that will be used by the scaffolder
	fs machinery.Filesystem
}

// NewInitScaffolder returns a new Scaffolder for project initialization operations
func NewInitScaffolder(config config.Config, license, owner string, isLegacyLayout bool) plugins.Scaffolder {
	return &initScaffolder{
		config:          config,
		boilerplatePath: hack.DefaultBoilerplatePath,
		license:         license,
		owner:           owner,
		isLegacyLayout:  isLegacyLayout,
	}
}

// InjectFS implements cmdutil.Scaffolder
func (s *initScaffolder) InjectFS(fs machinery.Filesystem) {
	s.fs = fs
}

// Scaffold implements cmdutil.Scaffolder
func (s *initScaffolder) Scaffold() error {
	fmt.Println("Writing scaffold for you to edit...")

	// Initialize the machinery.Scaffold that will write the boilerplate file to disk
	// The boilerplate file needs to be scaffolded as a separate step as it is going to
	// be used by the rest of the files, even those scaffolded in this command call.
	scaffold := machinery.NewScaffold(s.fs,
		machinery.WithConfig(s.config),
	)

	if s.license != "none" {
		bpFile := &hack.Boilerplate{
			License: s.license,
			Owner:   s.owner,
		}
		bpFile.Path = s.boilerplatePath
		if err := scaffold.Execute(bpFile); err != nil {
			return err
		}

		boilerplate, err := afero.ReadFile(s.fs.FS, s.boilerplatePath)
		if err != nil {
			return err
		}
		// Initialize the machinery.Scaffold that will write the files to disk
		scaffold = machinery.NewScaffold(s.fs,
			machinery.WithConfig(s.config),
			machinery.WithBoilerplate(string(boilerplate)),
		)
	} else {
		s.boilerplatePath = ""
		// Initialize the machinery.Scaffold without boilerplate
		scaffold = machinery.NewScaffold(s.fs,
			machinery.WithConfig(s.config),
		)
	}

	// If the KustomizeV2 was used to do the scaffold then
	// we need to ensure that we use its supported Kustomize Version
	// in order to support it
	helmVersion = helmv3.HelmVersion
	//helm := helmv3.Plugin{}
	//gov4 := "go.kubebuilder.io/v4"
	//pluginKeyForKustomizeV2 := plugin.KeyFor(helm)
	//
	//for _, pluginKey := range s.config.GetPluginChain() {
	//	if pluginKey == pluginKeyForKustomizeV2 || pluginKey == gov4 {
	//		kustomizeVersion = kustomizecommonv2alpha.HelmVersion
	//		break
	//	}
	//}

	return scaffold.Execute(
		&templates.Main{IsLegacyLayout: s.isLegacyLayout},
		&templates.GoMod{
			ControllerRuntimeVersion:   ControllerRuntimeVersion,
			EndpointOperatorLibVersion: EndpointOperatorLibVersion,
		},
		&templates.GitIgnore{},
		&templates.Makefile{
			Image:                       imageName,
			BoilerplatePath:             s.boilerplatePath,
			ControllerToolsVersion:      ControllerToolsVersion,
			ControllerToolsVersion4Helm: version.String(),
			HelmVersion:                 helmVersion,
			ControllerRuntimeVersion:    ControllerRuntimeVersion,
			EndpointOperatorLibVersion:  EndpointOperatorLibVersion,
			IsLegacyLayout:              s.isLegacyLayout,
		},
		&templates.Dockerfile{IsLegacyLayout: s.isLegacyLayout},
		&templates.DockerIgnore{},
		&templates.Readme{},
		&templates.Metadata{IsLegacyLayout: s.isLegacyLayout},
	)
}
