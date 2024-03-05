// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package deviceinfo

import (
	"os"
	"syscall"

	"github.com/intel/ixl-go/internal/config"
)

// Devices return all devices on current machine
func Devices() []*Device {
	ctx := config.NewContext()
	return ctx.Devices
}

// AvailableWorkQueues return all available workqueues.
func AvailableWorkQueues(typ DeviceType) []*WorkQueue {
	ctx := config.NewContext()
	wqs := ctx.WorkQueues(typ)
	temp := []*WorkQueue{}
	for _, wq := range wqs {
		_, err := os.Stat(wq.DevicePath())
		if err != nil {
			continue
		}
		// try to open workqueue files
		fd, err := syscall.Open(wq.DevicePath(), syscall.O_RDWR, 0)
		if err != nil {
			continue
		}
		syscall.Close(fd)
		temp = append(temp, wq)
	}
	return temp
}

type (
	WorkQueue = config.WorkQueue
	Device    = config.Device
	Group     = config.Group
	Engine    = config.Engine
)

// DeviceType is the type of device.
type DeviceType = config.DeviceType

const (
	DSA DeviceType = 0 // DSA represents a device type of DSA.
	IAA DeviceType = 1 // IAA represents a device type of IAA.
)
