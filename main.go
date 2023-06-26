// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"

	"github.com/hashicorp/packer-plugin-virtualbox/builder/virtualbox/iso"
	"github.com/hashicorp/packer-plugin-virtualbox/builder/virtualbox/ovf"
	"github.com/hashicorp/packer-plugin-virtualbox/builder/virtualbox/vm"
	"github.com/hashicorp/packer-plugin-virtualbox/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("iso", new(iso.Builder))
	pps.RegisterBuilder("ovf", new(ovf.Builder))
	pps.RegisterBuilder("vm", new(vm.Builder))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
