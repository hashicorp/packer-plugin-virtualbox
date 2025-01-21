// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iso

import (
	"testing"

	versionUtil "github.com/hashicorp/go-version"
)

func TestAudioDriverConfigurationArgForVersionsUsingAudioArg(t *testing.T) {
	versions := []string{
		"4.3",
		"4.3.1",
		"6.0.0",
		"6.0.2",
		"6.0.4",
		"6.0.6",
		"6.0.8",
		"6.0.10",
		"6.0.20",
		"6.0.22",
		"6.0.24",
		"6.0.28",
		"6.1.0",
		"6.1.2",
		"6.1.8",
		"6.1.10",
		"6.1.12",
		"6.1.14",
		"6.1.20",
		"6.1.22",
		"6.1.30",
		"6.1.32",
		"6.1.34",
		"6.1.38",
		"6.1.40",
		"6.1.42",
	}

	for _, in := range versions {
		in := in
		t.Run(in, func(t *testing.T) {
			got := audioDriverConfigurationArg(in)
			if got != "--audio" {
				t.Errorf("audioDriverConfigurationArg for with version %q returned %s but expected --audio", in, got)
			}

		})
	}
}
func TestAudioDriverConfigurationArgForVersionsUsingAudioDriverArg(t *testing.T) {
	versions := []string{
		"7.0.0",
		"7.0.2",
		"7.0.4",
		"7.0.6",
		"7.0.8",
		"7.0.10",
		"7.0.12",
		"7.0.14",
		"7.0.16",
	}

	for _, in := range versions {
		in := in
		t.Run(in, func(t *testing.T) {
			got := audioDriverConfigurationArg(in)
			if got != "--audio-driver" {
				t.Errorf("audioDriverConfigurationArg for with version %q returned %s but expected --audio-driver", in, got)
			}
		})
	}
}
func FuzzAudioDriverConfigurationArg(f *testing.F) {
	inputs := []string{
		"4.0",
		"4.1",
		"4.2",
		"4.3",
		"5.0",
		"5.1",
		"5.2",
		"6.0",
		"6.1",
		"7.0",
	}
	for _, in := range inputs {
		f.Add(in)
	}
	f.Fuzz(func(t *testing.T, input string) {
		expected := "--audio"

		got := audioDriverConfigurationArg(input)
		versionIn, _ := versionUtil.NewVersion(input)
		constraints, _ := versionUtil.NewConstraint(">= 7.0")
		if versionIn != nil && constraints.Check(versionIn) {
			expected = "--audio-driver"
		}
		if got != expected {
			t.Errorf("audioDriverConfigurationArg for with version %q returned %s but expected %s", versionIn, got, expected)
		}
	})
}
