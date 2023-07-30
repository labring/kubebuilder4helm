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

package main

import (
	"log"

	"github.com/labring/kubebuilder4helm/internal/version"
	"github.com/labring/kubebuilder4helm/plugins"
	golangv4 "github.com/labring/kubebuilder4helm/plugins/golang/v4"
	helmv1 "github.com/labring/kubebuilder4helm/plugins/helm/v3"
	"sigs.k8s.io/kubebuilder/v3/pkg/cli"
	cfgv3 "sigs.k8s.io/kubebuilder/v3/pkg/config/v3"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin"
)

func main() {
	gov4Bundle, _ := plugin.NewBundleWithOptions(plugin.WithName(plugins.DefaultNameQualifier),
		plugin.WithVersion(plugin.Version{Number: 4}),
		plugin.WithPlugins(helmv1.Plugin{}, golangv4.Plugin{}),
	)
	c, err := cli.New(
		cli.WithCommandName("kubebuilder4helm"),
		cli.WithVersion(version.String()),
		cli.WithExtraCommands(), // 如果有额外的命令
		cli.WithPlugins(
			gov4Bundle,
			// 可以添加其他插件
		),
		cli.WithDefaultPlugins(cfgv3.Version, gov4Bundle),
		cli.WithCompletion(),
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
