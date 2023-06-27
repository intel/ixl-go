// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package filter

import (
	"github.com/intel/ixl-go/errors"
	"github.com/intel/ixl-go/internal/device"
	"github.com/intel/ixl-go/internal/iaa"
	"github.com/intel/ixl-go/util/mem"
)

// Context is used to store filters state.
type Context struct {
	desc *iaa.Descriptor
	cr   *iaa.CompletionRecord
	aecs *iaa.FilterAECS
	ctx  *device.Context
}

// Ready returns true if the device is ready
func Ready() bool {
	return iaa.LoadContext().Ready()
}

// NewContext returns a new context
func NewContext() (*Context, error) {
	ctx := iaa.LoadContext()
	if ctx == nil {
		return nil, errors.NoHardwareDeviceDetected
	}
	return &Context{
		ctx:  iaa.LoadContext(),
		desc: mem.Alloc64Align[iaa.Descriptor](),
		cr:   mem.Alloc64Align[iaa.CompletionRecord](),
		aecs: mem.Alloc64Align[iaa.FilterAECS](),
	}, nil
}
