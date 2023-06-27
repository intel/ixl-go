// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package avx512 provides avx512 version PrepareForCodeLenCode
package avx512

import (
	"runtime"

	"github.com/intel/ixl-go/internal/iaa"
)

//go:noescape
func _prepareForCodeLenCode(h *iaa.Histogram, dest *byte, l, d *uint16)

// PrepareForCodeLenCode prepares code length code
func PrepareForCodeLenCode(h *iaa.Histogram, dest []byte) (litNum, disNum uint16) {
	_prepareForCodeLenCode(h, &dest[0], &litNum, &disNum)
	runtime.KeepAlive(dest)
	runtime.KeepAlive(h)
	return litNum, disNum
}
