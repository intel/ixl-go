// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package crc

import (
	"github.com/intel/ixl-go/internal/iaa"
)

// Ready returns true if the device is ready for use.
func Ready() bool {
	return iaa.LoadContext() != nil
}
