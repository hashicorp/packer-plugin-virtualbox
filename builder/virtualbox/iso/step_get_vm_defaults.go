// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iso

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	// packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	vboxcommon "github.com/hashicorp/packer-plugin-virtualbox/builder/virtualbox/common"
)

type stepGetVMDefaults struct {
}

func (s *stepGetVMDefaults) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	driver := state.Get("driver").(vboxcommon.Driver)

	//Add a default value for vmDefaults in the state bag
	state.Put("vmDefaults", make(map[string]string))

	baseFolder := os.TempDir()
	vmName := fmt.Sprintf("packer_temp_vm_%d", time.Now().Unix())
	defer os.RemoveAll(filepath.Join(baseFolder, vmName))

	// Create temp VM
	command := []string{"createvm", "--name", vmName, "--ostype", config.GuestOSType, "--register", "--default", "--basefolder", baseFolder}
	if err := driver.VBoxManage(command...); err != nil {
		err := fmt.Errorf("Failed to obtain VM defaults: %s", err)
		log.Println(err)
		return multistep.ActionContinue
	}

	defer func() {
		// Delete the temp VM
		command = []string{"unregistervm", vmName, "--delete"}
		if err := driver.VBoxManage(command...); err != nil {
			err := fmt.Errorf("Error deleting temp VM: %s", err)
			log.Println(err)
		}

	}()

	// Get the temp VM defaults
	command = []string{"showvminfo", vmName, "--machinereadable"}
	vminfo, err := driver.VBoxManageWithOutput(command...)
	if err != nil {
		err := fmt.Errorf("Failed to obtain VM defaults: %s", err)
		log.Println(err)
		return multistep.ActionContinue
	}

	defaults, err := parseVMInfoOutput(vminfo)
	if err != nil {
		err := fmt.Errorf("Error parsing VM defaults: %s", err)
		log.Println(err)
		return multistep.ActionContinue
	}

	// Store the defaults in the state bag
	state.Put("vmDefaults", defaults)
	log.Println("VM defaults retrieved successfully.")

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
