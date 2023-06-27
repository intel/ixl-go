// Copyright (c) 2023, Intel Corporation.
// SPDX-License-Identifier: BSD-3-Clause

// Package datamove provides functions which leverages DSA hardware abilities to copy memory.
package datamove

import (
	"log"
	"runtime"
	"sync"
	"unsafe"

	"github.com/intel/ixl-go/internal/config"
	"github.com/intel/ixl-go/internal/device"
	"github.com/intel/ixl-go/internal/dsa"

	"github.com/intel/ixl-go/util/mem"
)

// globalCtx is the device context used for all memory copy operations.
var globalCtx = device.CreateContext(config.DSA)

// Copy copies the content of the source byte slice to the destination byte slice.
// It returns true if the copy operation is successful, and false otherwise.
func Copy(dst []byte, src []byte) (ok bool) {
	if !Ready() {
		log.Println("[warn]no DSA device detected, fallback to software")
		copy(dst, src)
		return false
	}
	// Get a context from the context pool, and reset it.
	ctx, _ := pool.Get().(*Context)
	return ctx.Copy(dst, src)
}

// pool is the context pool for this package.
var pool *sync.Pool = &sync.Pool{
	New: func() any {
		return NewContext()
	},
}

// NewContext creates a new context.
// The context should be reused if possible.
func NewContext() *Context {
	return mem.Alloc32Align[Context]()
}

// Context represents a context for a memory copy operation.
// It should be reset before each use.
// It must be created by NewContext.
type Context struct {
	record dsa.CompletionRecord
	desc   dsa.Descriptor
}

// reset resets the context to its initial state.
func (c *Context) reset() {
	c.record = dsa.CompletionRecord{}
	c.desc = dsa.Descriptor{}
}

// Copy copies the content of the source byte slice to the destination byte slice.
// It returns true if the copy operation is successful, and false otherwise.
func (c *Context) Copy(dest, src []byte) bool {
	return c.CopyCheckError(dest, src) == nil
}

// CopyCheckError copies the content of the source byte slice to the destination byte slice.
// It returns nil if the copy operation is successful, and error otherwise.
func (c *Context) CopyCheckError(dest, src []byte) error {
	if !Ready() {
		log.Println("[warn]no DSA device detected, fallback to software")
		copy(dest, src)
		return nil
	}
	if len(dest) == 0 || len(src) == 0 {
		return nil
	}
	size := len(dest)
	if len(src) < size {
		size = len(src)
	}
	offset := 0
	// should check max transfer size
	for size > int(globalCtx.MaxTransferSize()) {
		c.reset()
		c.desc.SetFlags(dsa.OpFlagCRAddrValid | dsa.OpFlagReqCR | dsa.OpFlagBlockOnFault)
		c.desc.SetOpcode(dsa.OpcodeMemmove)
		c.desc.SrcAddr = uintptr(unsafe.Pointer(&src[offset]))
		c.desc.DestAddr = uintptr(unsafe.Pointer(&dest[offset]))
		c.desc.Size = (globalCtx.MaxTransferSize())
		c.desc.CompletionAddr = uintptr(unsafe.Pointer(&c.record))
		globalCtx.Submit(uintptr(unsafe.Pointer(&c.desc)), c.record.GetHeader())
		runtime.KeepAlive(dest)
		runtime.KeepAlive(src)
		offset += int(globalCtx.MaxTransferSize())
		err := c.record.CheckError()
		if err != nil {
			return err
		}
		size -= int(globalCtx.MaxTransferSize())
	}

	c.desc.SetFlags(dsa.OpFlagCRAddrValid | dsa.OpFlagReqCR | dsa.OpFlagBlockOnFault)
	c.desc.SetOpcode(dsa.OpcodeMemmove)
	c.desc.SrcAddr = uintptr(unsafe.Pointer(&src[offset]))
	c.desc.DestAddr = uintptr(unsafe.Pointer(&dest[offset]))
	c.desc.Size = uint32(size)
	c.desc.CompletionAddr = uintptr(unsafe.Pointer(&c.record))
	globalCtx.Submit(uintptr(unsafe.Pointer(&c.desc)), c.record.GetHeader())

	runtime.KeepAlive(dest)
	runtime.KeepAlive(src)
	return c.record.CheckError()
}

// Ready returns true if the device is ready for use.
func Ready() bool {
	return globalCtx != nil
}
