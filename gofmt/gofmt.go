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

package gofmt

import (
	"io"
	"os/exec"

	"github.com/kardianos/osext"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/godel-format-plugin/formatter"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const TypeName = "gofmt"

func Creator() formatter.Creator {
	return formatter.NewCreator(
		TypeName,
		func(cfgYML []byte) (formatter.Formatter, error) {
			// translate old configuration into new configuration if needed
			upgradedConfig, err := UpgradeConfig(cfgYML)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to load configuration")
			}
			var formatCfg Config
			if err := yaml.Unmarshal(upgradedConfig, &formatCfg); err != nil {
				return nil, errors.Wrapf(err, "failed to unmarshal YAML: %q", string(cfgYML))
			}
			return &gofmtFormatter{
				skipSimplify: formatCfg.SkipSimplify,
			}, nil
		},
	)
}

type gofmtFormatter struct {
	skipSimplify bool
}

func (f *gofmtFormatter) TypeName() (string, error) {
	return TypeName, nil
}

func (f *gofmtFormatter) Format(files []string, list bool, stdout io.Writer) error {
	self, err := osext.Executable()
	if err != nil {
		return errors.Wrapf(err, "failed to determine executable")
	}
	args := []string{
		amalgomated.ProxyCmdPrefix + TypeName,
	}
	if list {
		args = append(args, "-l")
	} else {
		args = append(args, "-w")
	}
	if !f.skipSimplify {
		args = append(args, "-s")
	}
	args = append(args, files...)

	cmd := exec.Command(self, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stdout
	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return errors.Wrapf(err, "failed to run %v", cmd.Args)
		}
	}
	return nil
}
