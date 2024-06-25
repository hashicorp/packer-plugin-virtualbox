// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/tmp"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"

	"strconv"
)

// This step attaches the boot ISO, cd_files iso, and guest additions to the
// virtual machine, if present.
type StepAttachISOs struct {
	AttachBootISO           bool
	ISOInterface            string
	GuestAdditionsMode      string
	GuestAdditionsInterface string
	diskUnmountCommands     map[string][]string
}

func (s *StepAttachISOs) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	// Check whether there is anything to attach
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Mounting ISOs...")
	diskMountMap := map[string]string{}
	s.diskUnmountCommands = map[string][]string{}
	// Track the bootable iso (only used in virtualbox-iso builder. )
	if s.AttachBootISO {
		isoPath := state.Get("iso_path").(string)
		diskMountMap["boot_iso"] = isoPath
	}

	// Determine if we even have a cd_files disk to attach
	if cdPathRaw, ok := state.GetOk("cd_path"); ok {
		cdFilesPath := cdPathRaw.(string)
		diskMountMap["cd_files"] = cdFilesPath
	}

	// Determine if we have guest additions to attach
	if s.GuestAdditionsMode != GuestAdditionsModeAttach {
		log.Println("Not attaching guest additions since we're uploading.")
	} else {
		// Get the guest additions path since we're doing it
		guestAdditionsPath := state.Get("guest_additions_path").(string)
		diskMountMap["guest_additions"] = guestAdditionsPath
	}

	if len(diskMountMap) == 0 {
		ui.Message("No ISOs to mount; continuing...")
		return multistep.ActionContinue
	}

	driver := state.Get("driver").(Driver)
	vmName := state.Get("vmName").(string)

	for diskCategory, isoPath := range diskMountMap {
		// If it's a symlink, resolve it to its target.
		resolvedIsoPath, err := filepath.EvalSymlinks(isoPath)
		if err != nil {
			err := fmt.Errorf("Error resolving symlink for ISO: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		isoPath = resolvedIsoPath

		// We have three different potential isos we can attach, so let's
		// assign each one its own spot so they don't conflict.
		var controllerName, diskType string
		var device, port int
		switch diskCategory {
		case "boot_iso":
			// figure out controller path
			controllerName = "IDE Controller"
			port = 0
			device = 1
			if s.ISOInterface == "sata" {
				controllerName = "SATA Controller"
				port = 13
				device = 0
			} else if s.ISOInterface == "virtio" {
				controllerName = "VirtIO Controller"
				port = 13
				device = 0
			}
			// If iso is a vhd, copy it to a temp dir and attach it as a hdd, else attach it as a dvddrive
			if filepath.Ext(isoPath) == ".vhd" {
				ui.Say("Copying boot VHD...")
				isoPath, err = s.CopyVHD(isoPath)
				if err != nil {
					err := fmt.Errorf("error copying VHD file: %s", err)
					state.Put("error", err)
					ui.Error(err.Error())
					return multistep.ActionHalt
				}
				diskType = "hdd"
			} else {
				diskType = "dvddrive"
			}
			ui.Message("Mounting boot ISO...")
		case "guest_additions":
			controllerName = "IDE Controller"
			port = 1
			device = 0
			diskType = "dvddrive"
			if s.GuestAdditionsInterface == "sata" {
				controllerName = "SATA Controller"
				port = 14
				device = 0
			} else if s.GuestAdditionsInterface == "virtio" {
				controllerName = "VirtIO Controller"
				port = 14
				device = 0
			}
			ui.Message("Mounting guest additions ISO...")
		case "cd_files":
			controllerName = "IDE Controller"
			port = 1
			device = 1
			diskType = "dvddrive"
			if s.ISOInterface == "sata" {
				controllerName = "SATA Controller"
				port = 15
				device = 0
			} else if s.ISOInterface == "virtio" {
				controllerName = "VirtIO Controller"
				port = 15
				device = 0
			}
			ui.Message("Mounting cd_files ISO...")
		}

		// Attach the disk to the controller
		command := []string{
			"storageattach", vmName,
			"--storagectl", controllerName,
			"--port", strconv.Itoa(port),
			"--device", strconv.Itoa(device),
			"--type", diskType,
			"--medium", isoPath,
		}
		if err := driver.VBoxManage(command...); err != nil {
			err := fmt.Errorf("Error attaching ISO: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		if diskType != "hdd" {
			// Track the disks we've mounted so we can remove them without having
			// to re-derive what was mounted where
			unmountCommand := []string{
				"storageattach", vmName,
				"--storagectl", controllerName,
				"--port", strconv.Itoa(port),
				"--device", strconv.Itoa(device),
				"--type", diskType,
				"--medium", "none",
			}

			s.diskUnmountCommands[diskCategory] = unmountCommand
		}
	}

	state.Put("disk_unmount_commands", s.diskUnmountCommands)

	return multistep.ActionContinue
}

func (s *StepAttachISOs) Cleanup(state multistep.StateBag) {
	if len(s.diskUnmountCommands) == 0 {
		return
	}

	driver := state.Get("driver").(Driver)
	_, ok := state.GetOk("detached_isos")

	if !ok {
		for _, command := range s.diskUnmountCommands {
			err := driver.VBoxManage(command...)
			if err != nil {
				log.Printf("error detaching iso: %s", err)
			}
		}
	}
}

func (s *StepAttachISOs) CopyVHD(path string) (string, error) {
	tempDir, err := tmp.Dir("virtualbox")
	if err != nil {
		return "", err
	}

	sourceFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(tempDir + "/" + filepath.Base(path))
	if err != nil {
		return "", err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return "", err
	}

	err = destinationFile.Chmod(0666)
	if err != nil {
		return "", err
	}
	newVHDPath := destinationFile.Name()

	return newVHDPath, nil
}
