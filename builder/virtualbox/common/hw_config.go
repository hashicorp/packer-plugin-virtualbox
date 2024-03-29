// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate packer-sdc struct-markdown

package common

import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type HWConfig struct {
	// The number of cpus to use for building the VM.
	// Defaults to 1.
	CpuCount int `mapstructure:"cpus" required:"false"`
	// The amount of memory to use for building the VM
	// in megabytes. Defaults to 512 megabytes.
	MemorySize int `mapstructure:"memory" required:"false"`
	// Defaults to none. The type of audio device to use for
	// sound when building the VM. Some of the options that are available are
	// dsound, oss, alsa, pulse, coreaudio, null.
	Sound string `mapstructure:"sound" required:"false"`
	// Specifies whether or not to enable the USB bus when
	// building the VM. Defaults to false.
	USB bool `mapstructure:"usb" required:"false"`
}

func (c *HWConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error

	// Hardware and cpu options
	if c.CpuCount < 0 {
		errs = append(errs, fmt.Errorf("An invalid number of cpus was specified (cpus < 0): %d", c.CpuCount))
	}
	if c.CpuCount == 0 {
		c.CpuCount = 1
	}

	if c.MemorySize < 0 {
		errs = append(errs, fmt.Errorf("An invalid memory size was specified (memory < 0): %d", c.MemorySize))
	}
	if c.MemorySize == 0 {
		c.MemorySize = 512
	}

	// devices
	if c.Sound == "" {
		c.Sound = "none"
	}

	return errs
}
