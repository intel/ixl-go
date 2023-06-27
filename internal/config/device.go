// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/intel/ixl-go/internal/log"
)

// Device represents a device with parameters.
type Device struct {
	MaxGroups         int    `binding:""` // Maximum number of groups.
	MaxWorkQueues     int    `binding:""` // Maximum number of work queues.
	MaxEngines        int    `binding:""` // Maximum number of engines.
	MaxWorkQueuesSize int    `binding:""` // Maximum size of work queues.
	NumaNode          int    `binding:""` // NUMA node.
	ImsSize           int    `binding:""` // IMS size.
	MaxBatchSize      int    `binding:""` // Maximum batch size.
	MaxTransferSize   uint64 `binding:""` // Maximum transfer size.
	GenCap            uint64 `binding:""` // Generation capabilities.
	Configurable      int    `binding:""` // Configurable parameters.
	PasidEnabled      bool   `binding:""` // PASID enabled.
	MaxReadBuffers    int    `binding:""` // Maximum number of read buffers.
	ReadBufferLimit   uint64 `binding:""` // Read buffer limit.
	CdevMajor         uint64 `binding:""` // cdev major number.
	Version           uint64 `binding:""` // Device version number.
	State             string `binding:""` // Device state.

	ID       uint64     // Device ID.
	OpCap    [4]uint64  // Device operation capabilities.
	Path     string     // Device's sysfs path.
	MDevPath string     // Device's mediated device path.
	Type     DeviceType // Device type.
	BusType  string     // Bus type.

	Groups     []*Group     // Groups associated with the device.
	Engines    []*Engine    `json:"-"` // Engines associated with the device.
	WorkQueues []*WorkQueue `json:"-"` // Work queues associated with the device.
}

// Device state constants.
const (
	DeviceStateUnknown  = "unknown"
	DeviceStateEnabled  = "enabled"
	DeviceStateDisabled = "disabled"
)

// Device capabilities constants.
const (
	BlockOnFaultSupport                   = 1 << iota // Block on fault support.
	OverlappingSupport                                // Overlapping support.
	CacheControlSupportMemory                         // Cache control support memory.
	CacheControlSupportCacheFlush                     // Cache control support cache flush.
	CommandCapabilitiesSupport                        // Command capabilities support.
	_                                                 // Reserved.
	_                                                 // Reserved.
	_                                                 // Reserved.
	DestinationReadbackSupport                        // Destination readback support.
	DrainDescriptorReadbackAddressSupport             // Drain descriptor readback address support.
	ConfigurationSupport                  = 1 << 31   // Configuration support.
)

// GenCap represents a device's generation capabilities.
type GenCap struct {
	BlockOnFault                   bool   // Block on fault support.
	Overlapping                    bool   // Overlapping support.
	CacheControlMemory             bool   // Cache control support memory.
	CacheControlCacheFlush         bool   // Cache control support cache flush.
	CommandCapabilities            bool   // Command capabilities support.
	DestinationReadback            bool   // Destination readback support.
	DrainDescriptorReadbackAddress bool   // Drain descriptor readback address support.
	ConfigurationSupport           bool   // Configuration support.
	MaxTransferSize                uint64 // Maximum transfer size.
	MaxBatchSize                   uint64 // Maximum batch size.
	InterruptMessageStorageSize    uint64 // Interrupt message storage size.
}

// ParseGenCap parses a uint64 value representing a device's generation capabilities.
func ParseGenCap(gencap uint64) (g *GenCap) {
	g = &GenCap{}
	if gencap&BlockOnFaultSupport == BlockOnFaultSupport {
		g.BlockOnFault = true
	}
	if gencap&OverlappingSupport == OverlappingSupport {
		g.Overlapping = true
	}
	if gencap&CacheControlSupportMemory == CacheControlSupportMemory {
		g.CacheControlMemory = true
	}
	if gencap&CommandCapabilitiesSupport == CommandCapabilitiesSupport {
		g.CommandCapabilities = true
	}
	if gencap&DestinationReadbackSupport == DestinationReadbackSupport {
		g.DestinationReadback = true
	}
	if gencap&DrainDescriptorReadbackAddressSupport == DrainDescriptorReadbackAddressSupport {
		g.DrainDescriptorReadbackAddress = true
	}
	if gencap&ConfigurationSupport == ConfigurationSupport {
		g.ConfigurationSupport = true
	}
	g.MaxTransferSize = 1 << ((gencap << (63 - 20)) >> (64 - 5))
	g.MaxBatchSize = 1 << ((gencap << (63 - 24)) >> (64 - 4))
	g.InterruptMessageStorageSize = 256 * (gencap << (63 - 30)) >> (64 - 6)
	return g
}

// ReadOPCap reads the op_cap file for a device and updates the device's op cap array.
func (d *Device) ReadOPCap() error {
	data, err := readTrimString(filepath.Join(d.Path, "op_cap"))
	if err != nil {
		return err
	}

	parts := strings.Split(data, " ")
	if len(parts) != 4 {
		return errors.New("unknown op cap format:" + data)
	}
	for i, p := range parts {
		if len(p) < 2 {
			return errors.New("unknown op cap format:" + data)
		}
		p := p[:2]
		num, err := strconv.ParseUint(p, 16, 64)
		if err != nil {
			return errors.New("unknown op cap format:" + data)
		}
		d.OpCap[i] = num
	}
	return nil
}

// Clients returns the number of clients for a device.
func (d *Device) Clients() (int, error) {
	return readInt(filepath.Join(d.Path, "clients"))
}

// AddDevice adds a Device to the context.
func (c *Context) AddDevice(id int, path string, busType string) (err error) {
	var dev Device
	dev.ID = uint64(id)
	err = bindSysfs(&dev, path)
	if err != nil {
		log.Debug("read sysfs %s failed: %v", path, err)
		return err
	}
	dev.Path = path

	mdevPath := filepath.Join("/sys/class/mdev_bus", filepath.Base(filepath.Dir(path)))
	dev.MDevPath = mdevPath

	if busType == IAA.Name() {
		dev.Type = IAA
	}

	if busType == DSA.Name() {
		dev.Type = DSA
	}

	dev.BusType = busType
	c.Devices = append(c.Devices, &dev)
	return
}
