// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package config

// Engine represents an engine in a device.
type Engine struct {
	// Device is the device this engine belongs to.
	Device *Device `json:"-"`
	// Group is the group this engine belongs to.
	Group *Group `json:"-"`
	// Type is the type of this engine.
	Type int
	// ID is the ID of this engine.
	ID int
	// GroupID is the ID of the group this engine belongs to.
	GroupID int `binding:""`
}

// addEngine adds an engine to a device, given its ID and path.
func (d *Device) addEngine(id int, path string) (err error) {
	// name format :engine${groupid}.${engineid}

	eng := &Engine{}
	eng.ID = id

	err = bindSysfs(eng, path)
	if err != nil {
		return err
	}

	for _, g := range d.Groups {
		if g.ID == eng.GroupID {
			eng.Group = g
			eng.Group.GroupedEngines = append(eng.Group.GroupedEngines, eng)
		}
	}

	eng.Device = d

	eng.Device.Engines = append(eng.Device.Engines, eng)

	return nil
}
