// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"path/filepath"

	"github.com/intel/ixl-go/internal/log"
)

// WorkQueue represents a work queue of a device.
type WorkQueue struct {
	Device          *Device `json:"-"` // Device associated with the work queue.
	Group           *Group  `json:"-"` // Group associated with the work queue.
	ID              int     // Work queue ID.
	NumaNode        int     // NUMA node.
	DeviceName      string  // Device name.
	Path            string  // Path to the work queue.
	GroupID         int     `binding:""` // ID of the group associated with the work queue.
	Size            int     `binding:""` // Work queue size.
	Priority        int     `binding:""` // Queue priority.
	BlockOnFault    int     `binding:""` // Block on fault flag.
	CdevMinor       int     `binding:""` // cdev minor number.
	Type            string  `binding:""` // Type of the work queue.
	Name            string  `binding:""` // Name of the work queue.
	Mode            string  `binding:""` // Mode of the work queue.
	State           string  `binding:""` // State of the work queue.
	DriverName      string  `binding:""` // Name of the driver.
	Threshold       uint    `binding:""` // Threshold limit.
	MaxBatchSize    uint    `binding:""` // Maximum batch size.
	MaxTransferSize uint64  `binding:""` // Maximum transfer size.
	AtsDisable      int     `binding:""` // ATS disable flag.
}

// Clients returns the number of clients for a work queue.
func (w *WorkQueue) Clients() (int, error) {
	return readInt(filepath.Join(w.Path, "clients"))
}

// DevicePath returns the device path of a work queue.
func (w *WorkQueue) DevicePath() string {
	return "/dev/" + w.Device.Type.Name() + "/" + w.DeviceName
}

// WQType represents the type of a work queue.
type WQType uint

const (
	_            = iota
	WQTypeKernel // WQTypeKernel indicates current wq is kernel mode.
	WQTypeUser   // WQTypeUser indicates current wq is User mode.
	WQTypeMDev   // WQTypeMDev indicates current wq is Mdev mode.
)

// addQueue adds a work queue to a device.
func (d *Device) addQueue(id int, path string) {
	wq := &WorkQueue{}
	wq.Path = path
	wq.ID = id
	wq.Device = d

	err := bindSysfs(wq, path)
	if err != nil {
		log.Debug("read sysfs failed: ", path, err)
		return
	}
	wq.DeviceName = filepath.Base(path)

	for _, g := range d.Groups {
		if g.ID == wq.GroupID {
			wq.Group = g
			wq.Group.GroupedWorkQueues = append(wq.Group.GroupedWorkQueues, wq)
		}
	}
	d.WorkQueues = append(d.WorkQueues, wq)
}
