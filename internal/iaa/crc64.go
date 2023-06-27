// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package iaa

import (
	"fmt"
	"unsafe"

	"github.com/intel/ixl-go/internal/device"
	"github.com/intel/ixl-go/internal/errors"
)

// CRC64Descriptor represents an Intel(R) Ethernet CRC64 descriptor.
type CRC64Descriptor struct {
	Header         uint32  `bf:"pasid 20;reversed 11;priv 1"` // The descriptor header.
	FlagsAndOpCode uint32  `bf:"flags 24;opcode 8"`           // The descriptor flags and opcode.
	CompletionAddr uintptr // The completion address for the descriptor.
	Src1Addr       uintptr // The source address for the descriptor.
	_              uintptr // Unused field.
	Size           uint32  // The size of the descriptor.
	IntHandle      uint16  // The interrupt handle for the descriptor.
	CRCFlag        CRCFlag // The CRC flag for the descriptor.
	_              uint64  // Unused field.
	_              uint64  // Unused field.
	CRCPolynomial  uint64  // The CRC polynomial for the descriptor.
}

// CRCFlag is an enumeration type for CRC flags.
type CRCFlag uint16

const (
	// CRCMostSignificant Indicates using the most significant bit for CRC calculation.
	CRCMostSignificant CRCFlag = 1 << 15
	// InvertCRC Indicates inverting the CRC result.
	InvertCRC CRCFlag = 1 << 14
)

// CRC64CompletionRecord represents an IAA CRC64 completion record.
type CRC64CompletionRecord struct {
	Header       device.CompletionRecordHeader // The completion record header.
	FaultAddress uint64                        // The fault address for the completion record.
	InvalidFlags uint32                        // The invalid flags for the completion record.
	_            uint32                        // Unused field.
	_            [2]uint64                     // Unused field.
	CRC64        uint64                        // The CRC result for the completion record.
	_            [2]uint64                     // Unused field.
}

// GetHeader returns the completion record header.
func (c *CRC64CompletionRecord) GetHeader() (crh CompletionRecordHeader) {
	return fromHeader(c.Header)
}

// CheckError checks if error happened.
func (c *CRC64CompletionRecord) CheckError() error {
	h := c.GetHeader()
	if h.StatusCode == Success {
		return nil
	}
	return errors.HardwareError{Status: h.StatusCode.String(), ErrorCode: h.ErrorCode.String()}
}

// SetCRCFlag sets the CRC flag for the descriptor.
func (d *CRC64Descriptor) SetCRCFlag(c CRCFlag) {
	d.CRCFlag = c
}

// AddCRCFlag adds a CRC flag to the descriptor.
func (d *CRC64Descriptor) AddCRCFlag(c CRCFlag) {
	d.CRCFlag = CRCFlag(uint16(d.CRCFlag) | uint16(c))
}

// GetCRCFlag returns the CRC flag for the descriptor.
func (d *CRC64Descriptor) GetCRCFlag() CRCFlag {
	return d.CRCFlag
}

// SetFlags sets the flags value for the descriptor.
func (d *CRC64Descriptor) SetFlags(value DescriptorFlag) {
	rv := uint32(value)

	left := (d.FlagsAndOpCode >> 24) << 24

	d.FlagsAndOpCode = left | rv
}

// SetPoly sets the CRC polynomial value for the descriptor.
func (d *CRC64Descriptor) SetPoly(poly uint64) {
	d.CRCPolynomial = poly
}

// SetCompleteRecord sets the completion address for the descriptor.
func (d *CRC64Descriptor) SetCompleteRecord(cr uintptr) {
	d.CompletionAddr = cr
}

// SetSourceData sets the source data for the descriptor.
func (d *CRC64Descriptor) SetSourceData(data []byte) {
	d.Src1Addr = uintptr(unsafe.Pointer(&data[0]))
	d.Size = uint32(len(data))
}

// SetOpcode sets the opcode value for the descriptor.
func (d *CRC64Descriptor) SetOpcode(value Opcode) {
	rv := uint32(value)

	right := (d.FlagsAndOpCode << (32 - 24)) >> (32 - 24)

	d.FlagsAndOpCode = (rv << 24) | right
}

// String returns a string representation of the CRC64 completion record.
func (c *CRC64CompletionRecord) String() string {
	return fmt.Sprintf("header:[%s]"+
		" fault_address:[%d]"+
		" invalid_flags:[%b]"+
		c.GetHeader().String(),
		c.FaultAddress,
		c.InvalidFlags,
	)
}

// fromHeader converts a device completion record header to a  completion record header.
func fromHeader(h device.CompletionRecordHeader) (crh CompletionRecordHeader) {
	crh.StatusCode = StatusCode((1<<7 - 1) & h.ComplexStatus)
	crh.StatusReadFault = (h.ComplexStatus >> 7) == 1
	crh.ErrorCode = ErrorCode(h.ErrorCode)
	crh.Completed = h.BytesCompleted
	return
}
