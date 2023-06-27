// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package dsa provides some useful types and functions for enable Intel DSA Engine.
package dsa

import (
	"unsafe"

	"github.com/intel/ixl-go/internal/device"
	"github.com/intel/ixl-go/internal/errors"
	"github.com/intel/ixl-go/util/mem"
)

// Descriptor is a structure defining a DSA descriptor.
type Descriptor struct {
	Header         uint32   `bf:"pasid 20;reversed 11;priv 1"` // Header of the descriptor.
	FlagsAndOpCode uint32   `bf:"flags 24;opcode 8"`           // Flags and opcode of the descriptor.
	CompletionAddr uintptr  // Completion address
	SrcAddr        uintptr  // Source address
	DestAddr       uintptr  // Destination address
	Size           uint32   // Transfer Size
	IntHandle      uint16   // Interrupt handle
	_              uint16   // Reserved field
	_              [24]byte // Reserved field
}

// GetFlags returns the flags of a descriptor.
func (d Descriptor) GetFlags() Flag { return Flag(d.FlagsAndOpCode<<8) >> 8 }

// SetFlags sets the flags of a descriptor.
func (d *Descriptor) SetFlags(value Flag) {
	left := (d.FlagsAndOpCode >> 24) << 24
	d.FlagsAndOpCode = left | uint32(value)
}

// GetOpcode returns the opcode of a descriptor.
func (d Descriptor) GetOpcode() Opcode { return Opcode(d.FlagsAndOpCode >> 24) }

// SetOpcode sets the opcode of a descriptor.
func (d *Descriptor) SetOpcode(value Opcode) {
	right := (d.FlagsAndOpCode << (32 - 24)) >> (32 - 24)
	d.FlagsAndOpCode = (uint32(value) << 24) | right
}

// String returns a string representation of a descriptor.
func (d *Descriptor) String() string {
	return ""
}

// CompletionRecord is a structure defining a DSA completion record.
type CompletionRecord struct {
	Header    uint64    // Header of the completion record.
	FaultAddr uintptr   // Fault address of the completion record.
	Record    [16]uint8 // Record of the completion record.
}

// NewCompletionRecord returns a new completion record.
func NewCompletionRecord() *CompletionRecord {
	return mem.Alloc32Align[CompletionRecord]()
}

// GetHeader returns the header of a completion record.
func (r *CompletionRecord) GetHeader() *device.CompletionRecordHeader {
	return (*device.CompletionRecordHeader)(unsafe.Pointer(&r.Header))
}

// CheckError checks if error happened, and return wrapped error or nil
func (r *CompletionRecord) CheckError() error {
	h := r.GetHeader()
	status := h.Status()
	if status == uint8(StatusSuccess) {
		return nil
	}
	return errors.HardwareError{Status: StatusCode(h.Status()).String()}
}
