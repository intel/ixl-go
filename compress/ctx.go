// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"github.com/intel/ixl-go/internal/iaa"
)

// Ready checks if the hardware is usable.
func Ready() bool {
	return iaa.LoadContext().Ready()
}
