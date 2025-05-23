// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// A driver is able to talk to VirtualBox and perform certain
// operations with it. Some of the operations on here may seem overly
// specific, but they were built specifically in mind to handle features
// of the VirtualBox builder for Packer, and to abstract differences in
// versions out of the builder steps, so sometimes the methods are
// extremely specific.
type Driver interface {
	// Create a SATA controller.
	CreateSATAController(vm string, controller string, portcount int) error

	// Create a SCSI controller.
	CreateSCSIController(vm string, controller string) error

	// Create a VirtIO controller.
	CreateVirtIOController(vm string, controller string) error

	// Create an NVME controller
	CreateNVMeController(vm string, controller string, portcount int) error

	// Delete all floppy controllers
	RemoveFloppyControllers(vm string) error

	// Delete a VM by name
	Delete(string) error

	// Import a VM
	Import(string, string, []string) error

	// The complete path to the Guest Additions ISO
	Iso() (string, error)

	// Checks if the VM with the given name is running.
	IsRunning(string) (bool, error)

	// Stop stops a running machine, forcefully.
	Stop(string) error

	// ACPIStop stops a running machine via ACPI power button.
	StopViaACPI(string) error

	// SuppressMessages should do what needs to be done in order to
	// suppress any annoying popups from VirtualBox.
	SuppressMessages() error

	// VBoxManage executes the given VBoxManage command
	// and returns an error
	VBoxManage(...string) error

	// VBoxManage executes the given VBoxManage command
	// and returns the stdout channel as string
	VBoxManageWithOutput(args ...string) (string, error)

	// Verify checks to make sure that this driver should function
	// properly. If there is any indication the driver can't function,
	// this will return an error.
	Verify() error

	// Version reads the version of VirtualBox that is installed.
	Version() (string, error)

	// LoadSnapshots Loads all defined snapshots for a vm.
	// if no snapshots are defined nil will be returned
	LoadSnapshots(string) (*VBoxSnapshot, error)

	// CreateSnapshot Creates a snapshot for a vm with a given name
	CreateSnapshot(string, string) error

	// HasSnapshots tests if a vm has snapshots
	HasSnapshots(string) (bool, error)

	// GetCurrentSnapshot Returns the current snapshot for a vm
	GetCurrentSnapshot(string) (*VBoxSnapshot, error)

	// SetSnapshot sets the for a vm
	SetSnapshot(string, *VBoxSnapshot) error

	// DeleteSnapshot deletes the specified snapshot from a vm
	DeleteSnapshot(string, *VBoxSnapshot) error
}

func NewDriver() (Driver, error) {
	var vboxmanagePath string

	// On Windows, we check VBOX_INSTALL_PATH env var for the path
	if runtime.GOOS == "windows" {
		vars := []string{"VBOX_INSTALL_PATH", "VBOX_MSI_INSTALL_PATH"}
		for _, key := range vars {
			value := os.Getenv(key)
			if value != "" {
				log.Printf(
					"[DEBUG] builder/virtualbox: %s = %s", key, value)
				vboxmanagePath = findVBoxManageWindows(value)
			}
			if vboxmanagePath != "" {
				break
			}
		}
	}

	if vboxmanagePath == "" {
		var err error
		vboxmanagePath, err = exec.LookPath("VBoxManage")
		if err != nil {
			return nil, err
		}
	}

	log.Printf("VBoxManage path: %s", vboxmanagePath)
	driver := &VBox42Driver{vboxmanagePath}
	if err := driver.Verify(); err != nil {
		return nil, err
	}

	return driver, nil
}

func findVBoxManageWindows(paths string) string {
	for _, path := range strings.Split(paths, ";") {
		path = filepath.Join(path, "VBoxManage.exe")
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}
