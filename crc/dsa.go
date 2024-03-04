// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

package crc

import (
	"hash"
	"runtime"
	"unsafe"

	"github.com/intel/ixl-go/errors"
	"github.com/intel/ixl-go/internal/dsa"
	"github.com/intel/ixl-go/util/mem"
)

var _ hash.Hash32 = (*CRC32C)(nil)

// NewCRC32C function initializes a new CRC32C hasher.
// This hasher is compliant with the hash.Hash32 interface.
// Returns an error if no hardware device is detected (Intel® DSA).
func NewCRC32C(opts ...CRC32COption) (*CRC32C, error) {
	// Check if the hardware device context is available
	if dsaContext == nil {
		// Return an error if no hardware device is detected
		return nil, errors.NoHardwareDeviceDetected
	}
	// Allocate aligned memory for the CRC32C hasher and return it
	return mem.Alloc32Align[CRC32C](), nil
}

type CRC32COption func(c *CRC32C)

// YieldProcessor yields the processor while  submitting  the CRC job to hardware, instead of busy polling the result.
func YieldProcessor(c *CRC32C) {
	c.yieldProccesor = true
}

// CRC32C represents a CRC32C calculator.
// This calculator uses Castagnoli's polynomial, which is widely used in iSCSI.
// The calculator is compliant with the hash.Hash32 interface.
type CRC32C struct {
	desc           dsa.CRCDescriptor
	record         dsa.CRCCompletionRecord
	crc            uint64
	yieldProccesor bool
}

func (c *CRC32C) reset() {
	c.record = dsa.CRCCompletionRecord{}
}

var dsaContext = dsa.LoadContext()

func (c *CRC32C) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	n = len(data)
	for len(data) > int(dsaContext.MaxTransferSize()) {
		slice := data[:dsaContext.MaxTransferSize()]
		_, err = c.write(slice)
		if err != nil {
			c.crc = 0
			return 0, err
		}
		data = data[dsaContext.MaxTransferSize():]
	}
	_, err = c.write(data)
	if err != nil {
		c.crc = 0
		return 0, err
	}
	return n, nil
}

func (c *CRC32C) write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	c.reset()
	c.desc.CRCSeed = c.crc
	c.desc.SetFlags(dsa.OpFlagCRAddrValid | dsa.OpFlagReqCR | dsa.OpFlagBlockOnFault)
	c.desc.SetOpcode(dsa.OpcodeCRCGen)
	c.desc.SrcAddr = uintptr(unsafe.Pointer(&data[0]))
	c.desc.Size = uint32(len(data))
	c.desc.CompletionAddr = uintptr(unsafe.Pointer(&c.record))
	if c.yieldProccesor {
		dsaContext.Submit(uintptr(unsafe.Pointer(&c.desc)), c.record.GetHeader())
	} else {
		dsaContext.SubmitBusyPoll(uintptr(unsafe.Pointer(&c.desc)), c.record.GetHeader())
	}
	runtime.KeepAlive(data)
	err = c.record.CheckError()
	if err != nil {
		return 0, err
	}
	c.crc = c.record.CRCValue
	return len(data), nil
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (c *CRC32C) Sum(b []byte) []byte {
	s := c.Sum32()
	return append(b, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}

// Reset resets the Hash to its initial state.
func (c *CRC32C) Reset() {
	c.reset()
	c.crc = 0
}

// Size returns the number of bytes Sum will return.
func (c *CRC32C) Size() int {
	return 4
}

// BlockSize returns the hash's underlying block size.
// The Write method must be able to accept any amount
// of data, but it may operate more efficiently if all writes
// are a multiple of the block size.
func (c *CRC32C) BlockSize() int {
	return 1
}

func (c *CRC32C) Sum32() uint32 {
	return uint32(c.crc)
}
