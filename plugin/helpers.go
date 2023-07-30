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

package plugin

import (
	"fmt"

	"github.com/spf13/afero"
	"sigs.k8s.io/kubebuilder/v3/pkg/machinery"
	"sigs.k8s.io/yaml"
)

type ConfigExtension struct {
	// If set to true, the plugin will use the legacy layout.
	// This is only used for testing purposes.
	IsLegacyLayout bool `yaml:"isLegacyLayout,omitempty"`
}

func GetConfigExtension() ConfigExtension {
	filename := "METADATA"
	fs := machinery.Filesystem{FS: afero.NewOsFs()}
	bs, err := afero.ReadFile(fs.FS, filename)
	if err != nil {
		fmt.Println("Using default config extension", err.Error())
		return ConfigExtension{}
	}
	config := &ConfigExtension{}
	err = yaml.Unmarshal(bs, config)
	if err != nil {
		return ConfigExtension{}
	}
	fmt.Println("Using config extension: isLegacyLayout", config.IsLegacyLayout)
	return *config
}

func SetConfigExtension(config *ConfigExtension) error {
	filename := "METADATA"
	bs, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	fs := machinery.Filesystem{FS: afero.NewOsFs()}
	_ = afero.WriteFile(fs.FS, filename, bs, 0644)
	return nil
}
