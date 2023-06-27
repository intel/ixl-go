// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package iaa

import (
	"sync"

	"github.com/intel/ixl-go/internal/config"
	"github.com/intel/ixl-go/internal/device"
)

var (
	globalCtx *device.Context
	ctxLoad   sync.Once
)

// LoadContext load iaa context
func LoadContext() *device.Context {
	ctxLoad.Do(func() {
		globalCtx = device.CreateContext(config.IAA)
	})
	return globalCtx
}

// DeviceReady checks if the device is ready
func DeviceReady() bool {
	return LoadContext() != nil
}
