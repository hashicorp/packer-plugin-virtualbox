// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"testing"
)

func TestVBox42Driver_impl(t *testing.T) {
	var _ Driver = new(VBox42Driver)
}
