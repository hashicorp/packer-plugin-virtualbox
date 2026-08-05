// Copyright IBM Corp. 2013, 2026
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

// TestStepAttachISOs_EFIUsesLowSATAPorts pins the fix for #39 and #20:
// VirtualBox's EFI firmware does not find a bootable device on the high SATA
// ports historically used for ISO attachment, so an EFI build with a
// SATA-attached installer ISO never boots — the VM sits at a black screen with
// an empty disk until the build times out waiting for SSH.
//
// BIOS guests must keep the historical ports so existing builds are unaffected.
func TestStepAttachISOs_EFIUsesLowSATAPorts(t *testing.T) {
	cases := []struct {
		firmware     string
		wantBootPort string
	}{
		{"efi", "1"},
		{"EFI", "1"}, // case-insensitive
		{"bios", "13"},
		{"", "13"}, // unset behaves as bios
	}

	// The step resolves symlinks on the ISO path, so it must exist on disk.
	isoPath := filepath.Join(t.TempDir(), "boot.iso")
	if err := os.WriteFile(isoPath, []byte("iso"), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, tc := range cases {
		state := testState(t)
		state.Put("vmName", "test-vm")
		step := new(StepAttachISOs)
		step.AttachBootISO = true
		step.ISOInterface = "sata"
		step.GuestAdditionsMode = GuestAdditionsModeDisable
		step.Firmware = tc.firmware

		state.Put("iso_path", isoPath)
		driver := state.Get("driver").(*DriverMock)

		if action := step.Run(context.Background(), state); action != multistep.ActionContinue {
			t.Fatalf("firmware %q: bad action %#v", tc.firmware, action)
		}

		var got string
		for _, cmd := range driver.VBoxManageCalls {
			if len(cmd) == 0 || cmd[0] != "storageattach" {
				continue
			}
			for i, a := range cmd {
				if a == "--port" && i+1 < len(cmd) {
					got = cmd[i+1]
				}
			}
		}
		if got != tc.wantBootPort {
			t.Fatalf("firmware %q: boot ISO attached to SATA port %q, want %q",
				tc.firmware, got, tc.wantBootPort)
		}
	}
}
