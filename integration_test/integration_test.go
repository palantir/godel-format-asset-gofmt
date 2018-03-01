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
	"io/ioutil"
	"os/exec"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/pkg/products"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoFmtAsset(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "")
	require.NoError(t, err)
	defer cleanup()

	origSrcContent := `package foo

import (
	_ "os"
	_ "fmt"
)

func Foo() {}
`
	specs := []gofiles.GoFileSpec{
		{
			RelPath: "foo.go",
			Src:     origSrcContent,
		},
	}
	files, err := gofiles.Write(tmpDir, specs)
	require.NoError(t, err)

	cli, err := products.Bin("gofmt-1.10.0-asset")
	require.NoError(t, err)

	srcFilePath := files["foo.go"].Path
	cmd := exec.Command(cli, srcFilePath)
	err = cmd.Run()
	require.NoError(t, err)

	content, err := ioutil.ReadFile(srcFilePath)
	require.NoError(t, err)

	assert.Equal(t, `package foo

import (
	_ "fmt"
	_ "os"
)

func Foo() {}
`, string(content))
}

func TestSimplify(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir(".", "")
	require.NoError(t, err)
	defer cleanup()

	origSrcContent := `package foo

func Foo() {
	for _ := range []string{} {
		_ = "foo"
	}
}
`
	specs := []gofiles.GoFileSpec{
		{
			RelPath: "foo.go",
			Src:     origSrcContent,
		},
	}
	files, err := gofiles.Write(tmpDir, specs)
	require.NoError(t, err)

	cli, err := products.Bin("gofmt-1.10.0-asset")
	require.NoError(t, err)

	srcFilePath := files["foo.go"].Path
	cmd := exec.Command(cli, "-s", srcFilePath)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, string(output))

	content, err := ioutil.ReadFile(srcFilePath)
	require.NoError(t, err)

	assert.Equal(t, `package foo

func Foo() {
	for range []string{} {
		_ = "foo"
	}
}
`, string(content))
}
