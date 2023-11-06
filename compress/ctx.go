// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package compress

import (
	"github.com/intel/ixl-go/internal/config"
	"github.com/intel/ixl-go/internal/iaa"
)

type iaadevice struct {
	ID            uint64 // Device ID.
	NumaNode      int    // NUMA node.
	NumWorkQueues int    // Number of work queues.
}

// Ready checks if the hardware is usable.
func Ready() bool {
	return iaa.LoadContext().Ready()
}

func GetDevices() []iaadevice {
	ctx := config.NewContext()
	ctx.Engines(config.IAA)
	var iaadevices []iaadevice

	for _, d := range ctx.Devices {
		if d.Type != config.IAA {
			continue
		}
		if d.State == config.DeviceStateEnabled {
			newd := iaadevice{ID: d.ID, NumaNode: d.NumaNode, NumWorkQueues: len(d.WorkQueues)}
			iaadevices = append(iaadevices, newd)
		}
	}
	return iaadevices
}
