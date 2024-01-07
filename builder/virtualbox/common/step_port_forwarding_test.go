// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"reflect"
	"testing"
)

func TestVirtualboxVersionIsAnInValidSemver(t *testing.T) {
	state := testState(t)

	state.Put("vmName", "foo")

	driver := state.Get("driver").(*DriverMock)
	driver.VersionResult = "v7.0.0abcd"
	driver.VersionErr = nil

	err := addAccessToLocalhost(state)

	if err == nil {
		t.Fatalf("We expected a failure but we got a success!")
	}
}

func TestVirtualboxVersionNeedsLocalhostAccessFlag(t *testing.T) {
	const versionRequiringFlag = "v7.0"
	const vmName = "foo"

	state := testState(t)

	state.Put("vmName", vmName)

	driver := state.Get("driver").(*DriverMock)
	driver.VersionResult = versionRequiringFlag
	driver.VersionErr = nil

	err := addAccessToLocalhost(state)

	if err != nil {
		t.Fatalf("Unexpected failure with VBox version '%v': %v", versionRequiringFlag, err)
	}

	if len(driver.VBoxManageCalls) == 0 {
		t.Fatal("VBoxManage wasn't called!")
	} else if len(driver.VBoxManageCalls) > 1 {
		t.Fatalf("Expected VBoxManage to be called once, but it was called %v times!", len(driver.VBoxManageCalls))
	}

	expectedArgs := []string{"modifyvm", vmName, "--nat-localhostreachable1", "on"}
	args := driver.VBoxManageCalls[0]
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Fatalf("Expected VBoxManage to be called with arguments %v, but it was called with arguments %v!", expectedArgs, args)
	}
}

func TestVirtualboxVersionDoesNotNeedLocalhostAccessFlag(t *testing.T) {
	const versionNotRequiringFlag = "v6.9"
	const vmName = "foo"

	state := testState(t)

	state.Put("vmName", vmName)

	driver := state.Get("driver").(*DriverMock)
	driver.VersionResult = versionNotRequiringFlag
	driver.VersionErr = nil

	err := addAccessToLocalhost(state)

	if err != nil {
		t.Fatalf("Unexpected failure with VBox version '%v': %v", versionNotRequiringFlag, err)
	}

	if len(driver.VBoxManageCalls) != 0 {
		t.Fatalf("Expected VBoxManage not to be called, but it was called %v times!", len(driver.VBoxManageCalls))
	}
}
