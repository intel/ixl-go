// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

import "github.com/intel/ixl-go/internal/log"

// Context represents the devices information which supported on the local machine.
type Context struct {
	Devices []*Device
}

const (
	ModeDedicated = "dedicated" // ModeDedicated represents a dedicated workqueue mode
	ModeShared    = "shared"    // ModeShared represents a shared workqueue mode
)

// Engines returns the engines of the specified device type.
func (c *Context) Engines(typ DeviceType) (results []*Engine) {
	for _, d := range c.Devices {
		if d.Type != typ {
			continue
		}
		if d.State == DeviceStateEnabled {
			results = append(results, d.Engines...)
		}
	}
	return
}

// SharedWQs returns the shared work queues of the specified device type.
func (c *Context) SharedWQs(typ DeviceType) []*WorkQueue {
	results := []*WorkQueue{}
	for _, d := range c.Devices {
		if d.Type != typ {
			continue
		}
		if d.State == DeviceStateEnabled {
			for _, wq := range d.WorkQueues {
				if wq.Mode == ModeShared {
					results = append(results, wq)
				}
			}
		}
	}
	return results
}

// DedicatedWQs returns the dedicated work queues of the specified device type.
func (c *Context) DedicatedWQs(typ DeviceType) []*WorkQueue {
	results := []*WorkQueue{}
	for _, d := range c.Devices {
		if d.Type != typ {
			continue
		}
		if d.State == DeviceStateEnabled {
			for _, wq := range d.WorkQueues {
				if wq.Mode == ModeDedicated {
					results = append(results, wq)
				}
			}
		}
	}
	return results
}

// WorkQueues returns the work queues of the specified device type.
func (c *Context) WorkQueues(typ DeviceType) []*WorkQueue {
	results := []*WorkQueue{}
	for _, d := range c.Devices {
		if d.Type != typ {
			continue
		}
		if d.State == DeviceStateEnabled {
			results = append(results, d.WorkQueues...)
		}
	}
	return results
}

// devicesLocation is the path of the devices directory.
const devicesLocation = "/sys/bus/dsa/devices"

// NewContext create a Context
func NewContext() *Context {
	ctx := &Context{}
	ctx.Init()
	return ctx
}

// Init starts to load devices information from the local sysfs.
func (c *Context) Init() {
	err := devDirs(
		devicesLocation,
		DSA.Name(),
		func(ids []int, path string) {
			err := c.AddDevice(ids[0], path, DSA.Name())
			if err != nil {
				log.Debug("error occurred while scan device:", err)
			}
		})
	if err != nil {
		log.Debug("error occurred while scan device:", err)
	}

	err = devDirs(
		devicesLocation,
		IAA.Name(),
		func(ids []int, path string) {
			err := c.AddDevice(ids[0], path, IAA.Name())
			if err != nil {
				log.Debug("error occurred while scan device:", err)
			}
		})
	if err != nil {
		log.Debug("error occurred while scan device:", err)
	}
	for _, d := range c.Devices {
		err = devDirs(d.Path, "group", func(ids []int, path string) {
			err = d.addGroup(ids[1], path)
			if err != nil {
				log.Debug("error occurred while add group:", err)
			}
		})
		if err != nil {
			log.Debug("error occurred while scan device:", err)
		}
		err = devDirs(d.Path, "engine", func(ids []int, path string) {
			err = d.addEngine(ids[1], path)
			if err != nil {
				log.Debug("error occurred while add group:", err)
			}
		})
		if err != nil {
			log.Debug("error occurred while scan device:", err)
		}
		err = devDirs(d.Path, "wq", func(ids []int, path string) {
			d.addQueue(ids[1], path)
		})
		if err != nil {
			log.Debug("error occurred while scan device:", err)
		}
	}
}

// DeviceType is a type of device.
type DeviceType uint8

const (
	DSA DeviceType = 0 // DSA represents a device type of DSA.
	IAA DeviceType = 1 // IAA represents a device type of IAA.
)

// String returns the string representation of the device type.
func (t DeviceType) String() string {
	switch t {
	case DSA:
		return "DSA"
	case IAA:
		return "IAA"
	}
	return ""
}

// Name returns the name of the device type.
func (t DeviceType) Name() string {
	switch t {
	case DSA:
		return "dsa"
	case IAA:
		return "iax"
	}
	return ""
}

// Names returns the names of the device type.
func (t DeviceType) Names() []string {
	switch t {
	case DSA:
		return []string{"dsa"}
	case IAA:
		return []string{"iaa", "iax"}
	}
	return nil
}
