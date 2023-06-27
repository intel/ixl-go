// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

// Group represents a group of devices managed by a device.
type Group struct {
	// Device is the device that manages the group.
	Device *Device `json:"-"`
	// ID is the ID of the group.
	ID int
	// Path is the path to the group.
	Path string
	// NumaNode is the NUMA node the group belongs to.
	NumaNode int
	// GroupID is the ID of the group.
	GroupID int

	// Engines is the list of engines in the group.
	Engines string `binding:""`
	// WorkQueues is the list of work queues in the group.
	WorkQueues string `binding:""`
	// ReadBuffersReserved is the number of read buffers reserved.
	ReadBuffersReserved uint64 `binding:""`
	// ReadBufferAllowed is the number of read buffers allowed.
	ReadBufferAllowed uint64 `binding:""`
	// UseReadBufferLimit is the number of read buffers to use.
	UseReadBufferLimit uint64 `binding:""`
	// TrafficClassA is the number of traffic class A messages.
	TrafficClassA uint64 `binding:""`
	// TrafficClassB is the number of traffic class B messages.
	TrafficClassB uint64 `binding:""`

	// GroupedEngines is the list of engines in the group.
	GroupedEngines []*Engine
	// GroupedWorkQueues is the list of work queues in the group.
	GroupedWorkQueues []*WorkQueue
}

// addGroup adds a new group to the device.
func (d *Device) addGroup(id int, path string) (err error) {
	var g Group
	g.Device = d
	g.ID = id
	g.Path = path
	err = bindSysfs(&g, path)
	if err != nil {
		return err
	}
	d.Groups = append(d.Groups, &g)
	return nil
}
