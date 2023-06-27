// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package crc

import (
	"math/bits"
	"unsafe"

	"github.com/intel/ixl-go/errors"
	"github.com/intel/ixl-go/internal/device"
	"github.com/intel/ixl-go/internal/iaa"
)

const (
	// ISO polynomial value
	ISO uint64 = 0xD800000000000000
	// ECMA polynomial value
	ECMA uint64 = 0xC96C5795D7870F42
)

const (
	// IEEE polynomial value
	IEEE uint32 = 0xedb88320
	// Castagnoli polynomial value
	Castagnoli uint32 = 0x82f63b78
	// Koopman polynomial value
	Koopman uint32 = 0xeb31d82e
)

const (
	// CCITT polynomial value
	CCITT uint16 = 0x8408
	// T10DIF polynomial value
	T10DIF uint16 = 0x8BB7
)

// Calculator is used for CRC64 calculation
// Notice: the data size should be less than your device's max_transfer_size.
type Calculator struct {
	d   *iaa.CRC64Descriptor
	cr  *iaa.CRC64CompletionRecord
	ctx *device.Context
}

// NewCalculator creates a new Calculator to be used for CRC64 calculation
func NewCalculator() (*Calculator, error) {
	ctx := iaa.LoadContext()
	if ctx == nil {
		// no device found
		return nil, errors.NoHardwareDeviceDetected
	}
	// create crc64 completion record
	cr := &iaa.CRC64CompletionRecord{}

	// create crc64 descriptor
	d := &iaa.CRC64Descriptor{}
	d.SetOpcode(iaa.OpCRC64)
	d.SetFlags(
		iaa.FlagBlockOnFault |
			iaa.FlagCacheControl |
			iaa.FlagRequestCompletionRecord |
			iaa.FlagCompletionRecordValid)
	d.SetCompleteRecord(uintptr(unsafe.Pointer(cr)))
	d.SetCRCFlag(iaa.CRCMostSignificant | iaa.InvertCRC)

	calc := &Calculator{
		ctx: ctx,
		d:   d,
		cr:  cr,
	}
	return calc, nil
}

// makeIAAPoly64 creates an IAA polynomial from a uint64
func makeIAAPoly64(poly uint64) uint64 {
	return bits.Reverse64(poly)
}

// makeIAAPoly32 creates an IAA polynomial from a uint32
func makeIAAPoly32(poly uint32) uint64 {
	return bits.Reverse64(uint64(poly))
}

// makeIAAPoly16 creates an IAA polynomial from a uint16
func makeIAAPoly16(poly uint16) uint64 {
	return bits.Reverse64(uint64(poly))
}

// prepare prepares the Calculator for CRC64 calculation
func (calc *Calculator) prepare(data []byte, iaaPoly uint64) {
	calc.d.SetPoly(iaaPoly)
	calc.d.SetSourceData(data)
}

// CheckSum64 calculates the CRC64 checksum for the given data and polynomial value
func (calc *Calculator) CheckSum64(data []byte, poly uint64) (uint64, error) {
	if len(data) == 0 {
		return 0, nil
	}
	if len(data) > int(calc.ctx.MaxTransferSize()) {
		return 0, errors.DataSizeTooLarge
	}
	calc.prepare(data, makeIAAPoly64(poly))
	status := iaa.StatusCode(calc.ctx.Submit(uintptr(unsafe.Pointer(calc.d)), &calc.cr.Header))
	if status != iaa.Success {
		return 0, calc.cr.CheckError()
	}
	return calc.cr.CRC64, nil
}

// CheckSum32 calculates the CRC32 checksum for the given data and polynomial value
func (calc *Calculator) CheckSum32(data []byte, poly uint32) (uint32, error) {
	if len(data) == 0 {
		return 0, nil
	}
	if len(data) > int(calc.ctx.MaxTransferSize()) {
		return 0, errors.DataSizeTooLarge
	}
	calc.prepare(data, makeIAAPoly32(poly))
	status := iaa.StatusCode(calc.ctx.Submit(uintptr(unsafe.Pointer(calc.d)), &calc.cr.Header))
	if status != iaa.Success {
		return 0, calc.cr.CheckError()
	}
	return uint32(calc.cr.CRC64), nil
}

// CheckSum16 calculates the CRC16 checksum for the given data and polynomial value
func (calc *Calculator) CheckSum16(data []byte, poly uint16) (uint16, error) {
	if len(data) == 0 {
		return 0, nil
	}
	if len(data) > int(calc.ctx.MaxTransferSize()) {
		return 0, errors.DataSizeTooLarge
	}
	calc.prepare(data, makeIAAPoly16(poly))
	status := iaa.StatusCode(calc.ctx.Submit(uintptr(unsafe.Pointer(calc.d)), &calc.cr.Header))
	if status != iaa.Success {
		return 0, calc.cr.CheckError()
	}
	return uint16(calc.cr.CRC64), nil
}
