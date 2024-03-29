// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

func TestStepRemoveDevices_impl(t *testing.T) {
	var _ multistep.Step = new(StepRemoveDevices)
}

func TestStepRemoveDevices(t *testing.T) {
	state := testState(t)
	step := new(StepRemoveDevices)

	state.Put("vmName", "foo")

	driver := state.Get("driver").(*DriverMock)

	// Test the run
	if action := step.Run(context.Background(), state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %#v", action)
	}
	if _, ok := state.GetOk("error"); ok {
		t.Fatal("should NOT have error")
	}

	// Test that ISO was removed
	if len(driver.VBoxManageCalls) != 0 {
		t.Fatalf("bad: %#v", driver.VBoxManageCalls)
	}
}

func TestStepRemoveDevices_attachedIso(t *testing.T) {
	state := testState(t)
	step := new(StepRemoveDevices)

	diskUnmountCommands := map[string][]string{
		"boot_iso": []string{
			"storageattach", "myvm",
			"--storagectl", "IDE Controller",
			"--port", "0",
			"--device", "1",
			"--type", "dvddrive",
			"--medium", "none",
		},
	}
	state.Put("disk_unmount_commands", diskUnmountCommands)
	state.Put("vmName", "foo")

	driver := state.Get("driver").(*DriverMock)

	// Test the run
	if action := step.Run(context.Background(), state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %#v", action)
	}
	if _, ok := state.GetOk("error"); ok {
		t.Fatal("should NOT have error")
	}

	// Test that ISO was removed
	if len(driver.VBoxManageCalls) != 1 {
		t.Fatalf("bad: %#v", driver.VBoxManageCalls)
	}
	if driver.VBoxManageCalls[0][3] != "IDE Controller" {
		t.Fatalf("bad: %#v", driver.VBoxManageCalls)
	}
}

func TestStepRemoveDevices_attachedIsoOnSata(t *testing.T) {
	state := testState(t)
	step := new(StepRemoveDevices)

	diskUnmountCommands := map[string][]string{
		"boot_iso": []string{
			"storageattach", "myvm",
			"--storagectl", "SATA Controller",
			"--port", "0",
			"--device", "1",
			"--type", "dvddrive",
			"--medium", "none",
		},
	}
	state.Put("disk_unmount_commands", diskUnmountCommands)
	state.Put("vmName", "foo")

	driver := state.Get("driver").(*DriverMock)

	// Test the run
	if action := step.Run(context.Background(), state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %#v", action)
	}
	if _, ok := state.GetOk("error"); ok {
		t.Fatal("should NOT have error")
	}

	// Test that ISO was removed
	if len(driver.VBoxManageCalls) != 1 {
		t.Fatalf("bad: %#v", driver.VBoxManageCalls)
	}
	if driver.VBoxManageCalls[0][3] != "SATA Controller" {
		t.Fatalf("bad: %#v", driver.VBoxManageCalls)
	}
}

func TestStepRemoveDevices_floppyPath(t *testing.T) {
	state := testState(t)
	step := new(StepRemoveDevices)

	state.Put("floppy_path", "foo")
	state.Put("vmName", "foo")

	driver := state.Get("driver").(*DriverMock)

	// Test the run
	if action := step.Run(context.Background(), state); action != multistep.ActionContinue {
		t.Fatalf("bad action: %#v", action)
	}
	if _, ok := state.GetOk("error"); ok {
		t.Fatal("should NOT have error")
	}

	// Test that both were removed
	if len(driver.VBoxManageCalls) != 2 {
		t.Fatalf("bad: %#v", driver.VBoxManageCalls)
	}
	if driver.VBoxManageCalls[0][3] != "Floppy Controller" {
		t.Fatalf("bad: %#v", driver.VBoxManageCalls)
	}
	if driver.VBoxManageCalls[1][3] != "Floppy Controller" {
		t.Fatalf("bad: %#v", driver.VBoxManageCalls)
	}
}
