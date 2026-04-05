// Copyright IBM Corp. 2013, 2025
// SPDX-License-Identifier: MPL-2.0

package iso

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	versionUtil "github.com/hashicorp/go-version"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	vboxcommon "github.com/hashicorp/packer-plugin-virtualbox/builder/virtualbox/common"
)

// This step creates the actual virtual machine.
//
// Produces:
//
//	vmName string - The name of the VM
type stepCreateVM struct {
	vmName string
}

func (s *stepCreateVM) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	driver := state.Get("driver").(vboxcommon.Driver)
	ui := state.Get("ui").(packersdk.Ui)

	name := config.VMName

	commands := [][]string{}
	commands = append(commands, []string{
		"createvm", "--name", name,
		"--ostype", config.GuestOSType, "--register",
	})
	commands = append(commands, []string{
		"modifyvm", name,
		"--boot1", "disk", "--boot2", "dvd", "--boot3", "none", "--boot4", "none",
	})
	commands = append(commands, []string{"modifyvm", name, "--cpus", strconv.Itoa(config.HWConfig.CpuCount)})
	if config.HWConfig.CpuCount > 1 {
		commands = append(commands, []string{"modifyvm", name, "--ioapic", "on"})
	}
	commands = append(commands, []string{"modifyvm", name, "--memory", strconv.Itoa(config.HWConfig.MemorySize)})
	// Configure USB controller
	if config.HWConfig.USB && strings.ToLower(config.USBController) != "none" {
		commands = append(commands, []string{"modifyvm", name, "--usb", "on"})
		commands = append(commands, []string{"modifyvm", name, "--usbohci", "off", "--usbehci", "off", "--usbxhci", "off"})
		switch strings.ToLower(config.USBController) {
		case "ohci":
			commands = append(commands, []string{"modifyvm", name, "--usbohci", "on"})
		case "ehci":
			commands = append(commands, []string{"modifyvm", name, "--usbehci", "on"})
		case "xhci":
			commands = append(commands, []string{"modifyvm", name, "--usbxhci", "on"})
		}
	} else {
		commands = append(commands, []string{"modifyvm", name, "--usb", "off"})
	}

	vboxVersion, _ := driver.Version()
	audioDriverArg := audioDriverConfigurationArg(vboxVersion)
	// Only configure audio if the audio controller is not set to "none"
	if strings.ToLower(config.AudioController) == "none" {
		commands = append(commands, []string{"modifyvm", name, "--audio-enabled", "off"})
	} else {
		if strings.ToLower(config.HWConfig.Sound) == "none" {
			commands = append(commands, []string{"modifyvm", name, audioDriverArg, config.HWConfig.Sound,
				"--audiocontroller", config.AudioController})
		} else {
			commands = append(commands, []string{"modifyvm", name, audioDriverArg, config.HWConfig.Sound, "--audioin", "on", "--audioout", "on",
				"--audiocontroller", config.AudioController})
		}
	}

	// Configure mouse device
	switch strings.ToLower(config.Mouse) {
	case "ps2":
		commands = append(commands, []string{"modifyvm", name, "--mouse", "ps2"})
	case "usb":
		commands = append(commands, []string{"modifyvm", name, "--mouse", "usb"})
	case "usbtablet":
		commands = append(commands, []string{"modifyvm", name, "--mouse", "usbtablet"})
	case "usbmultitouch":
		commands = append(commands, []string{"modifyvm", name, "--mouse", "usbmultitouch"})
	}

	// Configure keyboard device
	switch strings.ToLower(config.Keyboard) {
	case "ps2":
		commands = append(commands, []string{"modifyvm", name, "--keyboard", "ps2"})
	case "usb":
		commands = append(commands, []string{"modifyvm", name, "--keyboard", "usb"})
	}

	commands = append(commands, []string{"modifyvm", name, "--chipset", config.Chipset})
	commands = append(commands, []string{"modifyvm", name, "--firmware", config.Firmware})
	// Set the configured NIC type for all 8 possible NICs
	commands = append(commands, []string{"modifyvm", name,
		"--nictype1", config.NICType,
		"--nictype2", config.NICType,
		"--nictype3", config.NICType,
		"--nictype4", config.NICType,
		"--nictype5", config.NICType,
		"--nictype6", config.NICType,
		"--nictype7", config.NICType,
		"--nictype8", config.NICType})

	// Set the graphics controller, defaulting to "vboxsvga" unless overridden by "vmDefaults" or config.
	if config.GfxController == "" {
		config.GfxController = "vboxsvga"
		vmDefaultConfigs, defaultConfigsOk := state.GetOk("vmDefaults")
		if defaultConfigsOk {
			vmDefaultConfigs := vmDefaultConfigs.(map[string]string)
			if _, ok := vmDefaultConfigs["graphicscontroller"]; ok {
				config.GfxController = vmDefaultConfigs["graphicscontroller"]
			}
		}
	}

	commands = append(commands, []string{"modifyvm", name, "--graphicscontroller", config.GfxController, "--vram", strconv.FormatUint(uint64(config.GfxVramSize), 10)})
	if config.RTCTimeBase == "UTC" {
		commands = append(commands, []string{"modifyvm", name, "--rtcuseutc", "on"})
	} else {
		commands = append(commands, []string{"modifyvm", name, "--rtcuseutc", "off"})
	}

	if config.NestedVirt.True() {
		commands = append(commands, []string{"modifyvm", name, "--nested-hw-virt", "on"})
	} else if config.NestedVirt.False() {
		commands = append(commands, []string{"modifyvm", name, "--nested-hw-virt", "off"})
	}

	if config.GfxAccelerate3D {
		commands = append(commands, []string{"modifyvm", name, "--accelerate3d", "on"})
	} else {
		commands = append(commands, []string{"modifyvm", name, "--accelerate3d", "off"})
	}
	if config.GfxEFIResolution != "" {
		commands = append(commands, []string{"setextradata", name, "VBoxInternal2/EfiGraphicsResolution", config.GfxEFIResolution})
	}

	ui.Say("Creating virtual machine...")
	for _, command := range commands {
		err := driver.VBoxManage(command...)
		if err != nil {
			err := fmt.Errorf("Error creating VM: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		// Set the VM name property on the first command
		if s.vmName == "" {
			s.vmName = name
		}
	}

	// Set the final name in the state bag so others can use it
	state.Put("vmName", s.vmName)

	return multistep.ActionContinue
}

func (s *stepCreateVM) Cleanup(state multistep.StateBag) {
	if s.vmName == "" {
		return
	}

	driver := state.Get("driver").(vboxcommon.Driver)
	ui := state.Get("ui").(packersdk.Ui)
	config := state.Get("config").(*Config)

	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if (config.KeepRegistered) && (!cancelled && !halted) {
		ui.Say("Keeping virtual machine registered with VirtualBox host (keep_registered = true)")
		return
	}

	ui.Say("Deregistering and deleting VM...")
	if err := driver.Delete(s.vmName); err != nil {
		ui.Error(fmt.Sprintf("Error deleting VM: %s", err))
	}
}

func audioDriverConfigurationArg(vboxVersion string) string {
	// The '--audio' argument was deprecated in v7.0.x giving it
	//  the highest level of compatibility.
	compatibleAudioArg := "--audio"
	currentVersion, err := versionUtil.NewVersion(vboxVersion)
	if err != nil {
		log.Printf("[TRACE] attempt to read VBox version %q resulted in an error; using deprecated --audio argument: %s", vboxVersion, err)
		return compatibleAudioArg
	}

	constraints, _ := versionUtil.NewConstraint(">= 7.0")
	if currentVersion != nil && constraints.Check(currentVersion) {
		compatibleAudioArg = "--audio-driver"
	}

	return compatibleAudioArg
}
