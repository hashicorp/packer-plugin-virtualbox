// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// This step cleans up forwarded ports and exports the VM to an OVF.
// When DiskFormat is set to a non-VMDK format (e.g. VDI), the exported
// VMDK is converted to the requested format after export, and the OVF
// references are updated accordingly.
type StepExport struct {
	Format         string
	OutputDir      string
	OutputFilename string
	ExportOpts     []string
	Bundling       VBoxBundleConfig
	SkipNatMapping bool
	SkipExport     bool
	DiskFormat     string
}

func (s *StepExport) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	// If ISO export is configured, ensure this option is propagated to VBoxManage.
	for _, option := range s.ExportOpts {
		if option == "--iso" || option == "-I" {
			s.ExportOpts = append(s.ExportOpts, "--iso")
			break
		}
	}

	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packersdk.Ui)
	vmName := state.Get("vmName").(string)
	if s.OutputFilename == "" {
		s.OutputFilename = vmName
	}

	// Skip export if requested
	if s.SkipExport {
		ui.Say("Skipping export of virtual machine...")
		return multistep.ActionContinue
	}

	ui.Say("Preparing to export machine...")

	// Clear out the Packer-created forwarding rule
	commPort := state.Get("commHostPort")
	if !s.SkipNatMapping && commPort != 0 {
		ui.Message(fmt.Sprintf(
			"Deleting forwarded port mapping for the communicator (SSH, WinRM, etc) (host port %d)", commPort))
		command := []string{"modifyvm", vmName, "--natpf1", "delete", "packercomm"}
		if err := driver.VBoxManage(command...); err != nil {
			err := fmt.Errorf("Error deleting port forwarding rule: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	// Export the VM to an OVF
	outputPath := filepath.Join(s.OutputDir, s.OutputFilename+"."+s.Format)

	command := []string{
		"export",
		vmName,
		"--output",
		outputPath,
	}
	command = append(command, s.ExportOpts...)

	ui.Say("Exporting virtual machine...")
	ui.Message(fmt.Sprintf("Executing: %s", strings.Join(command, " ")))
	err := driver.VBoxManage(command...)
	if err != nil {
		err := fmt.Errorf("Error exporting virtual machine: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// If a non-VMDK disk format is requested, convert the exported VMDK
	// to the requested format and update the OVF references.
	diskFormat := strings.ToUpper(s.DiskFormat)
	if diskFormat != "" && diskFormat != "VMDK" {
		if err := s.convertExportedDisk(driver, ui, diskFormat); err != nil {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	state.Put("exportPath", outputPath)

	return multistep.ActionContinue
}

// convertExportedDisk converts the VMDK produced by VBoxManage export
// to the requested disk format (VDI, VHD) and updates the OVF file.
func (s *StepExport) convertExportedDisk(driver Driver, ui packersdk.Ui, diskFormat string) error {
	ext := strings.ToLower(diskFormat)

	// Step A: Find the VMDK file in the output directory
	pattern := filepath.Join(s.OutputDir, "*.vmdk")
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 0 {
		return fmt.Errorf("No VMDK file found in %s after export", s.OutputDir)
	}
	vmdkPath := matches[0]
	vmdkBasename := filepath.Base(vmdkPath)

	// Step B: Build the target disk path
	vdiBasename := strings.TrimSuffix(vmdkBasename, ".vmdk") + "." + ext
	vdiPath := filepath.Join(s.OutputDir, vdiBasename)

	ui.Say(fmt.Sprintf("Converting exported disk from VMDK to %s...", diskFormat))
	ui.Message(fmt.Sprintf("  %s -> %s", vmdkBasename, vdiBasename))

	// Step C: Convert VMDK to requested format using VBoxManage clonemedium
	err := driver.VBoxManage("clonemedium", "disk", vmdkPath, vdiPath, "--format", diskFormat)
	if err != nil {
		return fmt.Errorf("Error converting disk to %s: %s", diskFormat, err)
	}

	// Step D: Delete the VMDK
	if err := os.Remove(vmdkPath); err != nil {
		ui.Message(fmt.Sprintf("Warning: could not delete VMDK: %s", err))
	}

	// Step E: Read the OVF file
	ovfPath := filepath.Join(s.OutputDir, s.OutputFilename+"."+s.Format)
	ovfData, err := os.ReadFile(ovfPath)
	if err != nil {
		return fmt.Errorf("Error reading OVF file: %s", err)
	}

	// Step F: Replace VMDK references in OVF
	// There are exactly 2 references:
	//   1. File href: ovf:href="...disk001.vmdk"
	//   2. Disk format URL: ovf:format="http://www.vmware.com/interfaces/specifications/vmdk.html#streamOptimized"
	ovfContent := string(ovfData)
	ovfContent = strings.ReplaceAll(ovfContent, vmdkBasename, vdiBasename)
	ovfContent = strings.ReplaceAll(ovfContent,
		"http://www.vmware.com/interfaces/specifications/vmdk.html#streamOptimized",
		"http://www.virtualbox.org/VirtualBox/ExtPack/"+diskFormat)
	// Also handle the OVF 0.9 sparse format URL
	ovfContent = strings.ReplaceAll(ovfContent,
		"http://www.vmware.com/specifications/vmdk.html#sparse",
		"http://www.virtualbox.org/VirtualBox/ExtPack/"+diskFormat)

	// Step G: Write the OVF file back
	if err := os.WriteFile(ovfPath, []byte(ovfContent), 0644); err != nil {
		return fmt.Errorf("Error writing updated OVF file: %s", err)
	}

	ui.Say(fmt.Sprintf("Disk converted to %s successfully.", diskFormat))
	return nil
}

func (s *StepExport) Cleanup(state multistep.StateBag) {}
