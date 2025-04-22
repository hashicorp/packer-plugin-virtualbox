// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iso

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	vboxcommon "github.com/hashicorp/packer-plugin-virtualbox/builder/virtualbox/common"
)

type stepGetVMDefaults struct {
}

func (s *stepGetVMDefaults) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	driver := state.Get("driver").(vboxcommon.Driver)
	ui := state.Get("ui").(packersdk.Ui)

	baseFolder := os.TempDir()
	vmName := fmt.Sprintf("packer_temp_vm_%d", time.Now().Unix())

	// Create temp VM
	command := []string{"createvm", "--name", vmName, "--ostype", config.GuestOSType, "--register", "--default", "--basefolder", baseFolder}
	if err := driver.VBoxManage(command...); err != nil {
		err := fmt.Errorf("Error creating temp VM: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Get the temp VM defaults
	command = []string{"showvminfo", vmName, "--machinereadable"}
	vminfo, err := driver.VBoxManageWithOutput(command...)
	if err != nil {
		err := fmt.Errorf("Error reading VM info: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Delete the temp VM
	command = []string{"unregistervm", vmName, "--delete"}
	if err := driver.VBoxManage(command...); err != nil {
		err := fmt.Errorf("Error deleting temp VM: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	defaults, err := parseVMInfoOutput(vminfo)
	if err != nil {
		ui.Error("Error getting VM defaults: " + err.Error())
		return multistep.ActionHalt
	}

	// Store the defaults in the state bag
	state.Put("vmDefaults", defaults)
	ui.Message("VM defaults retrieved successfully.")

	return multistep.ActionContinue
}

func parseVMInfoOutput(output string) (map[string]string, error) {
	defaults := make(map[string]string)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				key := strings.Trim(strings.TrimSpace(parts[0]), `"`)
				value := strings.Trim(strings.TrimSpace(parts[1]), `"`)
				defaults[key] = value
			}
		}
	}
	if len(defaults) == 0 {
		return nil, fmt.Errorf("No defaults found in VM info output")
	}
	return defaults, nil
}

func (s *stepGetVMDefaults) Cleanup(state multistep.StateBag) {
	// No cleanup needed for this step
}
