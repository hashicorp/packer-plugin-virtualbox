// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iso

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/hashicorp/packer-plugin-sdk/acctest/testutils"
)

func TestAccBuilder_basic(t *testing.T) {
	templatePath := filepath.Join("testdata", "minimal.json")
	bytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to load template file %s", templatePath)
	}

	testCase := &acctest.PluginTestCase{
		Name:     "virtualbox-iso_basic_test",
		Template: string(bytes),
		Teardown: func() error {
			testutils.CleanupFiles("output-virtualbox-iso")
			return nil
		},
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}
