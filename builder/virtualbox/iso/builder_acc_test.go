// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iso

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func TestAccBuilder_defaultGfxController(t *testing.T) {
	templatePath := filepath.Join("testdata", "default_gfx.pkr.hcl")
	bytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to load template file %s", templatePath)
	}

	testCase := &acctest.PluginTestCase{
		Name:     "virtualbox-iso_default_gfx_test",
		Template: string(bytes),
		Teardown: func() error {
			testutils.CleanupFiles("output-ubuntu_2404")
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

func TestAccBuilder_incorrectGfxController(t *testing.T) {
	templatePath := filepath.Join("testdata", "incorrect_gfx.pkr.hcl")
	bytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to load template file %s", templatePath)
	}

	testCase := &acctest.PluginTestCase{
		Name:     "virtualbox-iso_incorrect_gfx_test",
		Template: string(bytes),
		Teardown: func() error {
			testutils.CleanupFiles("output-ubuntu_2404")
			return nil
		},
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			content, err := os.ReadFile(logfile)
			if err != nil {
				return fmt.Errorf("failed to read logfile: %v", err)
			}
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() == 0 || !strings.Contains(string(content), "Failed to construct device 'vga'") {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

func TestAccBuilder_headlessMode(t *testing.T) {
	templatePath := filepath.Join("testdata", "headless.pkr.hcl")
	bytes, err := ioutil.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to load template file %s", templatePath)
	}

	testCase := &acctest.PluginTestCase{
		Name:     "virtualbox-iso_headless_test",
		Template: string(bytes),
		Teardown: func() error {
			testutils.CleanupFiles("output-ubuntu_2404")
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
