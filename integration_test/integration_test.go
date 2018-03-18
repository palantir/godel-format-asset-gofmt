// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration_test

import (
	"fmt"
	"testing"

	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel-format-plugin/formattester"
	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/products"
	"github.com/stretchr/testify/require"
)

const (
	formatPluginLocator  = "com.palantir.godel-format-plugin:format-plugin:1.0.0-rc3"
	formatPluginResolver = "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"

	godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`
)

func TestGofmt(t *testing.T) {
	pluginProvider, err := pluginapitester.NewPluginProviderFromLocator(formatPluginLocator, formatPluginResolver)
	require.NoError(t, err)

	assetPath, err := products.Bin("gofmt-asset")
	require.NoError(t, err)

	configFiles := map[string]string{
		"godel/config/godel.yml":         godelYML,
		"godel/config/format-plugin.yml": "",
	}
	formattester.RunAssetFormatTest(t,
		pluginProvider,
		pluginapitester.NewAssetProvider(assetPath),
		[]formattester.AssetTestCase{
			{
				Name: "formats file",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo.go",
						Src: `package foo

import (
	_ "os"
	_ "fmt"
)

func Foo() {}
`,
					},
				},
				ConfigFiles: configFiles,
				WantFiles: map[string]string{
					"foo.go": `package foo

import (
	_ "fmt"
	_ "os"
)

func Foo() {}
`,
				},
			},
			{
				Name: "simplifies files",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo.go",
						Src: `package foo

func Foo() {
	for _ := range []string{} {
		_ = "foo"
	}
}
`,
					},
				},
				ConfigFiles: configFiles,
				WantFiles: map[string]string{
					"foo.go": `package foo

func Foo() {
	for range []string{} {
		_ = "foo"
	}
}
`,
				},
			},
			{
				Name: "does not simplify files if skip-simplify is true",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo.go",
						Src: `package foo

func Foo() {
	for _ := range []string{} {
		_ = "foo"
	}
}
`,
					},
				},
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/format-plugin.yml": `
formatters:
  gofmt:
    config:
      skip-simplify: true
`,
				},
				WantFiles: map[string]string{
					"foo.go": `package foo

func Foo() {
	for _ := range []string{} {
		_ = "foo"
	}
}
`,
				},
			},
			{
				Name: "verify does not modify files and prints unformatted files",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo.go",
						Src: `package foo

import (
	_ "os"
	_ "fmt"
)

func Foo() {}
`,
					},
				},
				ConfigFiles: configFiles,
				Verify:      true,
				WantError:   true,
				WantOutput: func(projectDir string) string {
					return fmt.Sprintf(`%s/foo.go
`, projectDir)
				},
				WantFiles: map[string]string{
					"foo.go": `package foo

import (
	_ "os"
	_ "fmt"
)

func Foo() {}
`,
				},
			},
		},
	)
}

func TestUpgradeConfig(t *testing.T) {
	pluginProvider, err := pluginapitester.NewPluginProviderFromLocator(formatPluginLocator, formatPluginResolver)
	require.NoError(t, err)

	assetPath, err := products.Bin("gofmt-asset")
	require.NoError(t, err)
	assetProvider := pluginapitester.NewAssetProvider(assetPath)

	pluginapitester.RunUpgradeConfigTest(t,
		pluginProvider,
		[]pluginapitester.AssetProvider{assetProvider},
		[]pluginapitester.UpgradeConfigTestCase{
			{
				Name: "current configuration is not upgraded",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/format-plugin.yml": `
# comment
formatters:
  gofmt:
    config:
      # inner comment
      skip-simplify: true
`,
				},
				WantOutput: "",
				WantFiles: map[string]string{
					"godel/config/format-plugin.yml": `
# comment
formatters:
  gofmt:
    config:
      # inner comment
      skip-simplify: true
`,
				},
			},
		},
	)
}
