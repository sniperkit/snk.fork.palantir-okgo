/*
Sniperkit-Bot
- Status: analyzed
*/

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

package checker

import (
	"encoding/base64"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/sniperkit/snk.fork.palantir-okgo/okgo"
)

type assetConfigUpgrader struct {
	typeName  okgo.CheckerType
	assetPath string
}

func (u *assetConfigUpgrader) TypeName() okgo.CheckerType {
	return u.typeName
}

func (u *assetConfigUpgrader) UpgradeConfig(config []byte) ([]byte, error) {
	upgradeConfigCmd := exec.Command(u.assetPath, "upgrade-config", base64.StdEncoding.EncodeToString(config))
	output, err := upgradeConfigCmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// if upgrade fails due to execution error, don't wrap it (will probably just be "exit 1")
			output := strings.TrimSuffix(strings.TrimPrefix(string(output), "Error: "), "\n")
			return nil, errors.Errorf("failed to upgrade asset configuration: %s", output)
		}
		return nil, errors.Wrapf(err, "failed to upgrade asset configuration: %s", output)
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(string(output))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode base64")
	}
	return decodedBytes, nil
}
