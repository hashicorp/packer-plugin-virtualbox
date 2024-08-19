// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/net"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"golang.org/x/mod/semver"
)

// This step adds a NAT port forwarding definition so that SSH or WinRM is available
// on the guest machine.
//
// Uses:
//
//	driver Driver
//	ui packersdk.Ui
//	vmName string
//
// Produces:
type StepPortForwarding struct {
	CommConfig       *communicator.Config
	HostPortMin      int
	HostPortMax      int
	SkipNatMapping   bool
	SSHListenAddress string

	l *net.Listener
}

var vboxVerMinNeedLHReachable string

func init() {
	vboxVerMinNeedLHReachable = semver.Canonical("v7.0")
	if vboxVerMinNeedLHReachable == "" {
		fmt.Fprintf(os.Stderr, "Constant version is invalid; this is a bug with the VirtualBox plugin, please open an issue to report it.\n")
		os.Exit(1)
	}
}

func addAccessToLocalhost(state multistep.StateBag) error {
	driver := state.Get("driver").(Driver)
	vmName := state.Get("vmName").(string)

	vboxVer, err := driver.Version()
	if err != nil {
		return fmt.Errorf("Error getting VirtualBox version: %s", err)
	}
	if !strings.HasPrefix(vboxVer, "v") {
		vboxVer = "v" + vboxVer
	}
	if !semver.IsValid(vboxVer) {
		return fmt.Errorf("The VirtualBox version isn't a valid SemVer: %s", vboxVer)
	}

	vboxVer = semver.Canonical(vboxVer)

	if semver.Compare(vboxVer, vboxVerMinNeedLHReachable) >= 0 {
		command := []string{
			"modifyvm", vmName,
			"--nat-localhostreachable1",
			"on",
		}
		if err := driver.VBoxManage(command...); err != nil {
			return fmt.Errorf("Failed to configure host's local network as reachable for NAT interface: %s", err)
		}
		log.Printf("[TRACE] VirtualBox's version (%q) is >= %q, setting --nat-localhostreachable1 on", vboxVer, vboxVerMinNeedLHReachable)
	}

	return nil
}

func (s *StepPortForwarding) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packersdk.Ui)
	vmName := state.Get("vmName").(string)

	if s.CommConfig.Type == "none" {
		log.Printf("Not using a communicator, skipping setting up port forwarding...")
		state.Put("commHostPort", 0)
		return multistep.ActionContinue
	}

	guestPort := s.CommConfig.Port()
	commHostPort := guestPort
	if !s.SkipNatMapping {
		log.Printf("Looking for available communicator (SSH, WinRM, etc) port between %d and %d",
			s.HostPortMin, s.HostPortMax)

		var err error
		s.l, err = net.ListenRangeConfig{
			Addr:    "127.0.0.1",
			Min:     s.HostPortMin,
			Max:     s.HostPortMax,
			Network: "tcp",
		}.Listen(ctx)
		if err != nil {
			err := fmt.Errorf("Error creating port forwarding rule: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		s.l.Listener.Close() // free port, but don't unlock lock file
		commHostPort = s.l.Port

		// Make sure to configure the network interface to NAT
		command := []string{
			"modifyvm", vmName,
			"--nic1",
			"nat",
		}
		if err := driver.VBoxManage(command...); err != nil {
			err := fmt.Errorf("Failed to configure NAT interface: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		// Add the `--nat-localhostreachableN=on` option if necessary
		if err := addAccessToLocalhost(state); err != nil {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		// Create a forwarded port mapping to the VM
		ui.Say(fmt.Sprintf("Creating forwarded port mapping for communicator (SSH, WinRM, etc) (host port %d)", commHostPort))

		command = []string{
			"modifyvm", vmName,
			"--natpf1",
			fmt.Sprintf("packercomm,tcp,%s,%d,,%d", s.SSHListenAddress, commHostPort, guestPort),
		}
		retried := false
	retry:
		if err := driver.VBoxManage(command...); err != nil {
			if !strings.Contains(err.Error(), "A NAT rule of this name already exists") || retried {
				err := fmt.Errorf("Error creating port forwarding rule: %s", err)
				state.Put("error", err)
				ui.Error(err.Error())
				return multistep.ActionHalt
			} else {
				log.Printf("A packer NAT rule already exists. Trying to delete ...")
				delcommand := []string{
					"modifyvm", vmName,
					"--natpf1",
					"delete", "packercomm",
				}
				if err := driver.VBoxManage(delcommand...); err != nil {
					err := fmt.Errorf("Error deleting packer NAT forwarding rule: %s", err)
					state.Put("error", err)
					ui.Error(err.Error())
					return multistep.ActionHalt
				}
				goto retry
			}
		}
	}

	// Save the port we're using so that future steps can use it
	state.Put("commHostPort", commHostPort)

	return multistep.ActionContinue
}

func (s *StepPortForwarding) Cleanup(state multistep.StateBag) {
	if s.l != nil {
		err := s.l.Close()
		if err != nil {
			log.Printf("failed to unlock port lockfile: %v", err)
		}
	}
}
